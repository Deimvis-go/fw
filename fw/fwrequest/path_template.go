package fwrequest

type PathTemplate string

type PathTemplateGenerationOption func(*canonicalPathCfg)

type PathTemplateFormat string

type canonicalPathCfg struct {
	format PathTemplateFormat
}

var (
	// BraceEscaping follows RFC 6570: https://datatracker.ietf.org/doc/html/rfc6570#section-2
	BraceEscaping PathTemplateFormat = "brace-escaping"
	ColonEscaping PathTemplateFormat = "colon-escaping"
)
