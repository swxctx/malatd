package td

import "encoding/json"

type (
	// Rerror 发生错误时返回的错误信息 response error
	Rerror struct {
		// Code 错误状态码
		Code int `json:"code"`
		// Message 错误信息
		Message string `json:"message"`
		// Reason 错误原因
		Reason string `json:"reason"`
	}
)

// NewRerror
func NewRerror(code int, message, reason string) *Rerror {
	return &Rerror{
		Code:    code,
		Message: message,
		Reason:  reason,
	}
}

// SetCode
func (r *Rerror) SetCode(message string) *Rerror {
	r.Message = message
	return r
}

// SetMessage
func (r *Rerror) SetMessage(message string) *Rerror {
	r.Message = message
	return r
}

// SetReason
func (r *Rerror) SetReason(reason string) *Rerror {
	r.Reason = reason
	return r
}

// String prints error info.
func (r *Rerror) String() string {
	if r == nil {
		return "<nil>"
	}
	b, _ := r.MarshalRerror()
	return BytesToString(b)
}

// MarshalRerror 错误信息编码
func (r *Rerror) MarshalRerror() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalRerror 解析错误信息
func (r *Rerror) UnmarshalRerror(data []byte) error {
	return json.Unmarshal(data, &r)
}
