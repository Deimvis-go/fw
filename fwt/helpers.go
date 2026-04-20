package fwt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis-go/fw/fw"
	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xfmt"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

func Test(t *testing.T, req fw.Request, resp fw.Response) {
	if t != nil {
		t.Helper()
	}
	w := Request(t, req)
	RequireResponse(t, resp, w)
}

func RequestAndDecode(t *testing.T, req fw.Request, resp fw.Response) {
	if t != nil {
		t.Helper()
	}
	w := Request(t, req)
	DecodeResponse(t, w, resp)
}

func Request(t *testing.T, r fw.Request) *httptest.ResponseRecorder {
	if t != nil {
		t.Helper()
	}
	rb := fw.MustBufferizeRequest(r)
	return RequestRaw(t, r.Method(), fmt.Sprintf("%s?%s", r.Path(), r.QueryString()), bytes.NewReader(rb.BodyRaw()), r.Header())
}

func RequestRaw(t *testing.T, method string, url string, body io.Reader, headers ...xhttp.ConstHeader) *httptest.ResponseRecorder {
	if t != nil {
		t.Helper()
	}
	refreshFn(cfg)

	req := xmust.Do(http.NewRequest(method, url, body))
	req.Header = mergeHeaders(headers...)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *cfg.ACCESS_TOKEN))

	cfg.LOGGER.Debugw("Will send request", "req", formatRequest(req))

	w := httptest.NewRecorder()
	cfg.SERVER.Handler.ServeHTTP(w, req)
	return w
}

func DecodeResponse(t *testing.T, w *httptest.ResponseRecorder, resp fw.Response) {
	if t != nil {
		t.Helper()
	}
	require.Equal(t, resp.Code(), w.Code)
	headersPtr := xmust.Do(ext.GetFieldAddr(resp, "Headers"))
	if headersPtr != nil {
		h := headersPtr.(*http.Header)
		(*h) = make(http.Header)
		for k, vs := range w.Header() {
			for _, v := range vs {
				(*h).Add(k, v)
			}
		}
	}
	bodyPtr := xmust.Do(ext.GetFieldAddr(resp, "Body"))
	if bodyPtr != nil {
		err := json.NewDecoder(w.Body).Decode(bodyPtr)
		require.NoError(t, err)
	}
}

func RequireResponse(t *testing.T, exp fw.Response, act *httptest.ResponseRecorder) {
	if t != nil {
		t.Helper()
	}
	RequireResponseRaw(t, exp.Code(), fw.MustBufferizeResponse(exp).BodyRaw(), act)
}

func RequireResponseRaw(t *testing.T, expCode int, expBody []byte, act *httptest.ResponseRecorder) {
	if t != nil {
		t.Helper()
	}
	require.Equal(t, expCode, act.Code)
	if len(expBody) > 0 {
		require.JSONEq(t, string(expBody), act.Body.String())
	} else {
		require.Len(t, act.Body.String(), 0)
	}
}

func mergeHeaders(headers ...xhttp.ConstHeader) http.Header {
	if len(headers) == 0 {
		return make(http.Header)
	}
	res := headers[0].Clone()
	for i := 0; i < len(headers); i++ {
		headers[i].Range(func(k string, vs []string) bool {
			for _, v := range vs {
				res.Add(k, v)
			}
			return true
		})
	}
	return res
}

func formatRequest(r *http.Request) string {
	body, bodyRaw := dupReader(r.Body)
	r.Body = io.NopCloser(body)
	kvs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"body", string(bodyRaw),
	}
	return fmt.Sprintf("request(%s)", xfmt.Sprintfkv("%v=%v", ", ", kvs...))
}

func dupReader(r io.Reader) (io.Reader, []byte) {
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)
	return r, buf.Bytes()
}
