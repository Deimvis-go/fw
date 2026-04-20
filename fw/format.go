package fw

import (
	"github.com/Deimvis-go/fw/fw/fwfmt"
)

// FormatRequest formats request into string.
// Implementation is not considered to be deterministic
// (do not check equality using format results).
func FormatRequest(r Request) string {
	return fwfmt.Request(r)
}

// FormatResponse formats response into string.
// Implementation is not considered to be deterministic
// (do not check equality using format results).
func FormatResponse(r Response) string {
	return fwfmt.Response(r)
}
