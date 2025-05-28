package codes

type Code struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e Code) New(msg string) Code {
	return Code{
		Code: e.Code,
		Msg:  msg,
	}
}

const (
	SuccessCode              = 0
	UnsupportedCurrencyType  = 100100
	UnsupportedMessageType   = 200100
	UnsupportedMessageFormat = 200200
	InvalidArgumentCode      = 400000
	InternalCode             = 500000
	UnimplementedCode        = 500001
)

var (
	Success                     = Code{Code: SuccessCode, Msg: "success"}
	ErrUnsupportedCurrencyType  = Code{Code: UnsupportedCurrencyType, Msg: "unsupported currency type"}
	ErrUnsupportedMessageType   = Code{Code: UnsupportedMessageType, Msg: "unsupported message type"}
	ErrUnsupportedMessageFormat = Code{Code: UnsupportedMessageFormat, Msg: "unsupported message format"}
	ErrInvalidArgument          = Code{Code: InvalidArgumentCode, Msg: "invalid params"}
	ErrUnimplemented            = Code{Code: UnimplementedCode, Msg: "unimplemented"}
	ErrInternal                 = Code{Code: InternalCode, Msg: "Unknown error"}
)
