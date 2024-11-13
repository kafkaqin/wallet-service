package controllers

const (
	// 通用错误码
	CODE_SUCCESS             = 000000
	CODE_NOT_FOUND           = 100001
	CODE_TIMEOUT             = 100002
	CODE_INTERNALSERVER      = 100003
	CODE_INVALID_PARAMS      = 100004 // 请求参数错误
	CODE_DATA_LEN_ERROR      = 100005 //数据格式错误
	CODE_REQUEST_TOO_QUICKLY = 100006
	// 用户
	CODE_USER_ROLE_NOT_EXISTS = 201001 //用户不存在

)

const (
	// 通用错误码
	ERRMSG_SUCCESS             string = "OK"
	ERRMSG_REQUEST_TOO_QUICKLY string = "request_too_quickly"
	ERRMSG_NOT_FOUND           string = "not_found"
	ERRMSG_TIMEOUT             string = "timeout"
	ERRMSG_INTERNAL_SERVER     string = "internal_server_error"
	ERRMSG_INVALID_PARAMS      string = "invalid params" // 请求参数错误
	ERRMSG_DATA_LEN_ERROR      string = "data_len_error" //数据格式错误

	// 用户
	ERRMSG_USER_ROLE_NOT_EXISTS string = "user role not exists" //角色不存在
)

var codeToERRMsgMap = map[int]string{
	CODE_SUCCESS:             ERRMSG_SUCCESS,
	CODE_NOT_FOUND:           ERRMSG_NOT_FOUND,
	CODE_REQUEST_TOO_QUICKLY: ERRMSG_REQUEST_TOO_QUICKLY,
	CODE_TIMEOUT:             ERRMSG_TIMEOUT,
	CODE_INTERNALSERVER:      ERRMSG_INTERNAL_SERVER,
	CODE_INVALID_PARAMS:      ERRMSG_INVALID_PARAMS, // 请求参数错误
	CODE_DATA_LEN_ERROR:      ERRMSG_DATA_LEN_ERROR, //数据格式错误

	// 用户
	CODE_USER_ROLE_NOT_EXISTS: ERRMSG_USER_ROLE_NOT_EXISTS, //角色不存在

}
