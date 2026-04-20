package types

// NOTE: DRAFT

// TODO: what if request partially structured?
//   - interface for each field or these types should return NilT on Struct*() when non-structured fields?
type StructuredRequest[UriT any, QueryT any, HeaderT any, BodyT any] interface {
	Request
	StructURI() UriT
	StructQuery() QueryT
	HeaderT() HeaderT
	BodyT() BodyT
}
