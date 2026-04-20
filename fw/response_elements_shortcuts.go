package fw

import "github.com/Deimvis-go/fw/fw/fwresponse"

type Response200WithHeader struct {
	Response200
	ResponseHeader[NoHeaderPreset]
	ResponseNoBody
}

type Response200WithJSONHeader struct {
	Response200
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type Response200NoHeader struct {
	Response200
	ResponseNoHeader
	ResponseNoBody
}

type Response202WithJSONHeader struct {
	Response202
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type Response202NoHeader struct {
	Response202
	ResponseNoHeader
	ResponseNoBody
}

type Response302WithHeader struct {
	Response302
	ResponseHeader[NoHeaderPreset]
	ResponseNoBody
}

type Response302WithJSONHeader struct {
	Response302
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type Response400NoHeader struct {
	Response400
	ResponseNoHeader
	ResponseNoBody
}

type Response400WithJSONHeader struct {
	Response400
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type Response403NoHeader struct {
	Response403
	ResponseNoHeader
	ResponseNoBody
}

type Response403WithJSONHeader struct {
	Response403
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type Response404NoHeader struct {
	Response404
	ResponseNoHeader
	ResponseNoBody
}

type Response404WithJSONHeader struct {
	Response404
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type Response500NoHeader struct {
	Response500
	ResponseNoHeader
	ResponseNoBody
}

type Response500WithJSONHeader struct {
	Response500
	ResponseHeader[JSONHeaderPreset]
	ResponseNoBody
}

type ErrorResponse400NoHeader struct {
	Response400
	ResponseNoHeader
	fwresponse.ErrorResponse
}

type ErrorResponse401NoHeader struct {
	Response401
	ResponseNoHeader
	fwresponse.ErrorResponse
}

type ErrorResponse403NoHeader struct {
	Response403
	ResponseNoHeader
	fwresponse.ErrorResponse
}

type ErrorResponse404NoHeader struct {
	Response404
	ResponseNoHeader
	fwresponse.ErrorResponse
}

type ErrorResponse409NoHeader struct {
	Response409
	ResponseNoHeader
	fwresponse.ErrorResponse
}

type ErrorResponse500NoHeader struct {
	Response500
	ResponseNoHeader
	fwresponse.ErrorResponse
}

type ErrorResponse503NoHeader struct {
	Response503
	ResponseNoHeader
	fwresponse.ErrorResponse
}
