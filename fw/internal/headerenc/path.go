package headerenc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/ext"
)

type path struct {
	v []string
}

func (p *path) String() string {
	return strings.Join(p.v, "")
}

func (p *path) Clone() *path {
	return &path{v: ext.CopyElements(p.v)}
}

func (p *path) WithRootObject(obj any, fn func()) {
	p.PushRootObject(obj)
	defer p.Pop()
	fn()
}

func (p *path) WithStructField(sf reflect.StructField, fn func()) {
	p.PushStructField(sf)
	defer p.Pop()
	fn()
}

func (p *path) WithSliceElement(ind int, fn func()) {
	p.PushSliceElement(ind)
	defer p.Pop()
	fn()
}

func (p *path) WithMapKey(key reflect.Value, fn func()) {
	p.PushMapKey(key)
	defer p.Pop()
	fn()
}

func (p *path) PushRootObject(obj any) {
	if len(p.v) != 0 {
		panic("bug: object is not root")
	}
	var name string
	v := extractInternalValue(reflect.ValueOf(obj))
	if v.Kind() == reflect.Invalid {
		name = "nil"
	} else {
		vt := v.Type()
		if vt.Name() != "" {
			name = vt.Name()
		} else {
			// NOTE: reflect.Type.String() may include dot in type definitions
			name = fmt.Sprintf("<%s>", vt.String())
			// k := type_.Kind()
			// if k == reflect.Array || k == reflect.Slice || k == reflect.Map || k == reflect.Chan {
			// 	name += "{}"
			// }
		}
	}
	p.v = append(p.v, name)
}

func (p *path) PushStructField(sf reflect.StructField) {
	p.v = append(p.v, fmt.Sprintf(".%s", sf.Name))
}

func (p *path) PushSliceElement(ind int) {
	p.v = append(p.v, fmt.Sprintf("[%d]", ind))
}

func (p *path) PushMapKey(key reflect.Value) {
	keyStr := p.formatMapKey(key)
	p.v = append(p.v, fmt.Sprintf("[%s]", keyStr))
}

func (p *path) Pop() {
	p.v = p.v[:len(p.v)-1]
}

func (p *path) formatMapKey(key reflect.Value) string {
	if key.CanInterface() {
		if s, ok := key.Interface().(fmt.Stringer); ok {
			return s.String()
		}
	}
	if key.Kind() == reflect.String {
		return fmt.Sprintf("\"%s\"", key.String())
	}
	return key.String()
}

// extractInternalValue resolves pointers and interfaces until value becomes of different kind.
// https://github.com/go-playground/validator/blob/a947377040f8ebaee09f20d09a745ec369396793/util.go#L15
func extractInternalValue(current reflect.Value) reflect.Value {

BEGIN:
	switch current.Kind() {
	case reflect.Pointer:

		if current.IsNil() {
			return current
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:

		if current.IsNil() {
			return current
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current

	default:
		return current
	}
}
