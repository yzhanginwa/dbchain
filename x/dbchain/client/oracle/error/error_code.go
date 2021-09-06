package error

const (
	Success = "success"
	PasswordErr = "password err"
	UnLoginErr = "not login"
	QueryErr = "query err"
	ServerErr = "server err"
	UndefinedErr = "undefined Err"
	TelNoVerify = "tel is not verify"
	FormatErr = "format err"
	ParamsErr = "params err"
	UnregisterErr = "unregister err"
	SoldOutErr = "sold out"
	RegisteredErr = "registered"
	AllBookedErr = "all booked"
	Unauthorized = "Unauthorized"
	)

const (
	SuccessCode = "0"
	PasswordErrCode = "1"
	UnLoginErrCode = "2"
	QueryErrCode = "3"
	ServerErrCode = "4"
	UndefinedErrCode = "5"
	TelNoVerifyCode = "6"
	FormatErrCode = "7"
	ParamsErrCode = "8"
	UnregisterErrCode = "9"
	SoldOutErrCode = "10"
	RegisteredErrCode = "11"
	AllBookedErrCode = "12"
	UnauthorizedErrCode = "13"
)

var ErrDescription = map[string]string{
	SuccessCode : Success,
	PasswordErrCode : PasswordErr,
	UnLoginErrCode : UnLoginErr,
	QueryErrCode : QueryErr,
	ServerErrCode : ServerErr,
	UndefinedErrCode : UndefinedErr,
	TelNoVerifyCode :  TelNoVerify,
	FormatErrCode : FormatErr,
	ParamsErrCode : ParamsErr,
	UnregisterErrCode : UnregisterErr,
	SoldOutErrCode : SoldOutErr,
	RegisteredErrCode : RegisteredErr,
	AllBookedErrCode : AllBookedErr,
	UnauthorizedErrCode : Unauthorized,
}
