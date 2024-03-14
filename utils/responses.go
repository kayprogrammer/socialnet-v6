package utils

type ErrorResponse struct {
	Status				string				`json:"status"`
	Code				string				`json:"code"`
	Message				string				`json:"message"`
	Data				*map[string]string	`json:"data,omitempty"`
}

func (obj ErrorResponse) Init() ErrorResponse {
	if obj.Status == "" {
		obj.Status = "failure"
	}
	return obj
}

// Error codes
var ERR_UNAUTHORIZED_USER =	"unauthorized_user"
var ERR_NETWORK_FAILURE =	"network_failure"
var ERR_SERVER_ERROR =	"server_error"
var ERR_INVALID_REQUEST =	"invalid_request"
var ERR_INVALID_ENTRY =	"invalid_entry"
var ERR_INCORRECT_EMAIL =	"incorrect_email"
var ERR_INCORRECT_OTP =	"incorrect_otp"
var ERR_EXPIRED_OTP =	"expired_otp"
var ERR_INVALID_AUTH =	"invalid_auth"
var ERR_INVALID_TOKEN =	"invalid_token"
var ERR_INVALID_CREDENTIALS =	"invalid_credentials"
var ERR_UNVERIFIED_USER =	"unverified_user"
var ERR_NON_EXISTENT =	"non_existent"
var ERR_INVALID_OWNER =	"invalid_owner"
var ERR_INVALID_PAGE =	"invalid_page"
var ERR_INVALID_VALUE =	"invalid_value"
var ERR_NOT_ALLOWED =	"not_allowed"
var ERR_INVALID_DATA_TYPE =	"invalid_data_type"

func RequestErr(code string, message string, opts ...map[string]string) ErrorResponse {
	var data *map[string]string
	// Check if data is provided as an argument
	if len(opts) > 0 {
		data = &opts[0]
	}
	resp := ErrorResponse{Code: code, Message: message, Data: data}.Init()
	return resp
}