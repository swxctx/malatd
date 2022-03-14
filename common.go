package td

import "math"

const (
	// 公用错误码
	CodeUnknownError = -1
	CodeWriteFailed  = 104

	CodeInvalidParameter    = 299
	CodeBadPacket           = 400
	CodeUnauthorized        = 401
	CodeNotFound            = 404
	CodeInternalServerError = 500
	CodeBadGateway          = 502
)

// CodeMessage
func CodeMessage(rerrCode int) string {
	switch rerrCode {
	case CodeBadPacket:
		return "Bad Packet"
	case CodeUnauthorized:
		return "Unauthorized"
	case CodeWriteFailed:
		return "Write Failed"
	case CodeNotFound:
		return "Not Found"
	case CodeInternalServerError:
		return "Internal Server Error"
	case CodeBadGateway:
		return "Bad Gateway"
	case CodeUnknownError:
		fallthrough
	case CodeInvalidParameter:
		return "Invalid Parameter"
	default:
		return "Unknown Error"
	}
}

var (
	// 公用错误
	RerrUnknown          = NewRerror(CodeUnknownError, CodeMessage(CodeUnknownError), "")
	RerrWriteFailed      = NewRerror(CodeWriteFailed, CodeMessage(CodeWriteFailed), "")
	RerrInvalidParameter = NewRerror(CodeInvalidParameter, CodeMessage(CodeInvalidParameter), "")
	RerrBadPacket        = NewRerror(CodeBadPacket, CodeMessage(CodeBadPacket), "")
	RerrNotFound         = NewRerror(CodeNotFound, CodeMessage(CodeNotFound), "")
	RerrInternalServer   = NewRerror(CodeInternalServerError, CodeMessage(CodeInternalServerError), "")
	RerrBadGateway       = NewRerror(CodeBadGateway, CodeMessage(CodeBadGateway), "")
)

const (
	abortIndex int = math.MaxInt8 / 2
)
