package headerenc

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// source: json package
// patch log:
// - no field.encoder
// - no structFields.byExactName
// - modified structFields.byFoldedName (same as json's, but we do not truncate name to 32 bytes so we are slower)
// - no field dominance ranking, return error on tag name collision
// - no tag opts
// - no field.quoted
// - no field.omitEmpty
// - no field.nameNonEsc
// - no field.nameEscHTML
// - no go:linkname for typeFields

// A field represents a single field found in a struct.
type field struct {
	name      string
	nameBytes []byte // []byte(name)

	// nameNonEsc  string // `"` + name + `":`
	// nameEscHTML string // `"` + HTMLEscape(name) + `":`

	tag   bool
	index []int
	typ   reflect.Type
	// omitEmpty bool
	// quoted    bool

	// encoder encoderFunc
}

type structFields struct {
	list []field
	// byExactName map[string]*field
	byFoldedName map[string]*field
}

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
//
// typeFields should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/bytedance/sonic
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
func typeFields(t reflect.Type) (structFields, error) {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	// // Buffer to run appendHTMLEscape on field names.
	// var nameEscBuf []byte

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("header")
				if tag == "-" {
					continue
				}
				name := strings.ToLower(tag)
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := field{
						name:  name,
						tag:   tagged,
						index: index,
						typ:   ft,
						// omitEmpty: opts.Contains("omitempty"),
						// quoted:    quoted,
					}
					field.nameBytes = []byte(field.name)

					// // Build nameEscHTML and nameNonEsc ahead of time.
					// nameEscBuf = appendHTMLEscape(nameEscBuf[:0], field.nameBytes)
					// field.nameEscHTML = `"` + string(nameEscBuf) + `":`
					// field.nameNonEsc = `"` + field.name + `":`

					fields = append(fields, field)
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 and 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	// NOTE: we ignore field dominance ranking,
	// we return error on header collision

	// slices.SortFunc(fields, func(a, b field) int {
	// 	// sort field by name, breaking ties with depth, then
	// 	// breaking ties with "name came from json tag", then
	// 	// breaking ties with index sequence.
	// 	if c := strings.Compare(a.name, b.name); c != 0 {
	// 		return c
	// 	}
	// 	if c := cmp.Compare(len(a.index), len(b.index)); c != 0 {
	// 		return c
	// 	}
	// 	if a.tag != b.tag {
	// 		if a.tag {
	// 			return -1
	// 		}
	// 		return +1
	// 	}
	// 	return slices.Compare(a.index, b.index)
	// })

	// // Delete all fields that are hidden by the Go rules for embedded fields,
	// // except that fields with JSON tags are promoted.

	// // The fields are sorted in primary order of name, secondary order
	// // of field index length. Loop over names; for each name, delete
	// // hidden fields by choosing the one dominant field that survives.
	// out := fields[:0]
	// for advance, i := 0, 0; i < len(fields); i += advance {
	// 	// One iteration per name.
	// 	// Find the sequence of fields with the name of this first field.
	// 	fi := fields[i]
	// 	name := fi.name
	// 	for advance = 1; i+advance < len(fields); advance++ {
	// 		fj := fields[i+advance]
	// 		if fj.name != name {
	// 			break
	// 		}
	// 	}
	// 	if advance == 1 { // Only one field with this name
	// 		out = append(out, fi)
	// 		continue
	// 	}
	// 	return structFields{}, fmt.Errorf("multiple fields in %s with same name %s", t.String(), name)

	// 	// dominant, ok := dominantField(fields[i : i+advance])
	// 	// if ok {
	// 	// 	out = append(out, dominant)
	// 	// }
	// }

	// fields = out
	// slices.SortFunc(fields, func(i, j field) int {
	// 	return slices.Compare(i.index, j.index)
	// })

	foldedNameIndex := make(map[string]*field, len(fields))
	for i, field := range fields {
		fname := strings.ToLower(field.name)
		if _, ok := foldedNameIndex[fname]; ok {
			return structFields{}, fmt.Errorf("multiple fields in %s with same name %s", t.String(), fname)
		}
		foldedNameIndex[fname] = &fields[i]
	}
	return structFields{fields, foldedNameIndex}, nil
}

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []field) (field, bool) {
	// The fields are sorted in increasing index-length order, then by presence of tag.
	// That means that the first field is the dominant one. We need only check
	// for error cases: two fields at top level, either both tagged or neither tagged.
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tag == fields[1].tag {
		return field{}, false
	}
	return fields[0], true
}

var fieldCache sync.Map // map[reflect.Type]structFields

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) (structFields, error) {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields), nil
	}
	fields, err := typeFields(t)
	if err != nil {
		return structFields{}, err
	}
	f, _ := fieldCache.LoadOrStore(t, fields)
	return f.(structFields), nil
}
