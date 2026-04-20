package headerenc

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xinvar"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// TODO: support custom encoding.
//       header:"key,default" or header:"key" - "default" encoding (works for int, float, string; int is 10-based)
//       header:"key,json" - "json" encoding: decode json, cast to type
//       header:"key,bool,booltrue=1 true" - "bool" encoding: decode bool by matching header value string to true values
//       header:"key,bool,boolfalse=0 false" - "bool" encoding: decode bool by matching header value string to false values
//         _other impl_ header:"key,bool=1 true ! 0 false" - "bool" encoding: decode bool by matching header value string. if not matched - error
//         ^^^ require explicit of true/false values
//         _enhance impl_ header:"key,bool=...,boolopts=fb(false) regex_match()" - fallback value is false, when not matched (not error), match as regex
// TODO: support match by case sensetive header name  as tag option?

func Unmarshal(h xhttp.ConstHeader, out any, notMatchedOut http.Header /* optional */) error {
	// - TODO: implement xgenc.GenericMarshaler + xenccore.Marshaler (xencoding/xgenc, xencoding/xenccore)
	//         add settings:
	// 		   - max depth,
	// 		   - dominant field resolver,
	//         - no field matched input/output action (ignore/error/store somewhere/...) (e.g. non-matched header, non-matched out struct field),
	//         - whether optimistic - whether to enable validation before unmarshalling
	// 		     validation includes reader's data validation and output value validation (e.g. that out type belongs to fixed set, works on each recursion level)
	// 		     (developer should be able to provide interfaces used for validation heck and their methods:
	//			  e.g. interface { ValidHeaderStruct() error }, or specify validatoin by type)
	// 	 - TODO: how to implement Marshaling from map[string]any to struct without encoding map in bytes?

	rout := reflect.ValueOf(out)
	if rout.Kind() != reflect.Pointer || rout.IsNil() {
		return fmt.Errorf("invalid output (non pointer or nil pointer %s)", rout.Type().String())
	}

	rout = indirect(rout, false)
	switch rout.Kind() {
	case reflect.Map:
		return errors.New("map output is not supported yet")
	case reflect.Struct:
		v := rout
		vt := v.Type()
		fields, err := cachedTypeFields(vt)
		if err != nil {
			return err
		}
		err = nil
		h.Range(func(key string, values []string) bool {
			if len(values) == 0 {
				return true
			}
			// source: json package
			f, ok := fields.byFoldedName[strings.ToLower(key)]
			if !ok {
				if notMatchedOut != nil {
					notMatchedOut.Del(key)
					for _, val := range values {
						notMatchedOut.Add(key, val)
					}
				}
				return true
			}

			subv := v
			for _, i := range f.index {
				if subv.Kind() == reflect.Pointer {
					if subv.IsNil() {
						// If a struct embeds a pointer to an unexported type,
						// it is not possible to set a newly allocated value
						// since the field is unexported.
						//
						// See https://golang.org/issue/21357
						if !subv.CanSet() {
							err = fmt.Errorf("fwheader: cannot set embedded pointer to unexported struct: %v", subv.Type().Elem())
							// Invalidate subv to ensure d.value(subv) skips over
							// the JSON value without assigning it to subv.
							subv = reflect.Value{}
							break
						}
						subv.Set(reflect.New(subv.Type().Elem()))
					}
					subv = subv.Elem()
				}
				subv = subv.Field(i)
			}
			if err != nil {
				return false
			}

			subv = indirect(subv, false)
			err = unmarshalValues(key, values, subv)
			if err != nil {
				return false
			}

			return true

		})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported input kind: %s", rout.Kind())
	}
	return nil
}

func unmarshalValues(key string, vs []string, out reflect.Value) error {
	var err error
	switch out.Kind() {
	case reflect.Array, reflect.Slice:
		if len(vs) > out.Cap() {
			if out.Kind() == reflect.Array {
				return fmt.Errorf("output array is smaller than number of header values: %d < %d", out.Cap(), len(vs))
			}
			out.Grow(len(vs) - out.Cap())
		}
		if len(vs) > out.Len() {
			if out.Kind() == reflect.Array {
				return fmt.Errorf("output array is smaller than number of header values: %d < %d", out.Len(), len(vs))
			}
			out.SetLen(len(vs))
		}
		for i := 0; i < len(vs); i++ {
			err = unmarshalScalar(vs[i], out.Index(i))
			if err != nil {
				break
			}
		}
	default:
		if len(vs) > 1 {
			return fmt.Errorf("header key has multiple values, but expected single: %s", key)
		}
		xinvar.Eq(len(vs), 1)
		err = unmarshalScalar(vs[0], out)
	}
	return err
}

func unmarshalScalar(v string, out reflect.Value) error {
	switch out.Kind() {
	case reflect.Int,
		reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8,
		reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		vv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		out.SetInt(vv)
	case reflect.Float64, reflect.Float32:
		vv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		out.SetFloat(vv)
	case reflect.String:
		out.SetString(v)
	default:
		return fmt.Errorf("unsupported header value kind: %T", out.Kind())
	}
	return nil
}
