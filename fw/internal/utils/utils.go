package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/xstrconv"
)

func MakeHttpHeader(headers any) http.Header {
	v := reflect.ValueOf(headers)
	vt := v.Type()
	h := make(http.Header, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Pointer && f.IsNil() {
			// ignore field as not-filled optional
			continue
		}
		ft := vt.Field(i)
		if ft.Anonymous {
			// ignore embedded field
			continue
		}
		key := vt.Field(i).Tag.Get("header")
		if key == "" {
			panic(fmt.Errorf(`no "header" tag in field %s of type %s`, vt.Field(i).Name, vt.Name()))
		}
		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				h.Add(key, toString(f.Index(j)))
			}
		} else {
			h.Add(key, toString(f))
		}
	}
	return h
}

func MakeQueryString(query any) string {
	v := reflect.ValueOf(query)
	vt := v.Type()
	kvs := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Pointer && f.IsNil() {
			continue
		}
		key := vt.Field(i).Tag.Get("query")
		if key == "" {
			panic(fmt.Errorf(`no "query" tag in field %s of type %s`, vt.Field(i).Name, vt.Name()))
		}
		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				elem := f.Index(j)
				queryValue := url.QueryEscape(toString(elem))
				kvs = append(kvs, fmt.Sprintf("%s=%s", key, queryValue))
			}
		} else {
			queryValue := url.QueryEscape(toString(f))
			kvs = append(kvs, fmt.Sprintf("%s=%s", key, queryValue))
		}
	}
	return strings.Join(kvs, "&")
}

func MakeJsonBodyRaw(body any) []byte {
	res, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return res
}

func toString(v reflect.Value) string {
	// explicitly handle nil pointer case.
	// otherwise, indirect will use underlying type default value.
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return ""
	}
	v = reflect.Indirect(v)
	return formatter.MustFormat(v.Interface())
}

var (
	formatter xstrconv.Formatter
)

func init() {
	formatter = xstrconv.NewFormatter(
		xstrconv.FormatFns{
			PerKind: xstrconv.KindFormatFns{
				Int64: func(i int64) (string, error) {
					return strconv.FormatInt(i, 10), nil
				},
				Uint64: func(ui uint64) (string, error) {
					return strconv.FormatUint(ui, 10), nil
				},
				Float32: func(f float32) (string, error) {
					return strconv.FormatFloat(float64(f), 'f', -1, 32), nil
				},
				Float64: func(f float64) (string, error) {
					return strconv.FormatFloat(float64(f), 'f', -1, 64), nil
				},
				String: func(s string) (string, error) {
					return s, nil
				},
				Bool: func(b bool) (string, error) {
					return strconv.FormatBool(b), nil
				},
			},
		},
		xstrconv.WithIntFormatPropagation(),
		xstrconv.WithUintFormatPropagation(),
	)
}
