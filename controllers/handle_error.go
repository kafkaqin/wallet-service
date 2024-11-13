package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleError 通用错误处理函数
func handleError(c *gin.Context, errorCode int, err error) {
	if err != nil {
		var statusCode int
		switch errorCode {
		case CODE_NOT_FOUND, CODE_USER_ROLE_NOT_EXISTS:
			statusCode = http.StatusOK
		case CODE_INVALID_PARAMS, CODE_DATA_LEN_ERROR:
			statusCode = http.StatusBadRequest
		case CODE_TIMEOUT:
			statusCode = http.StatusGatewayTimeout
		default:
			statusCode = http.StatusInternalServerError
		}
		var responseData ResponseData
		if msg, ok := codeToERRMsgMap[errorCode]; ok {
			responseData.ErrorMsg = msg
			responseData.ErrorCode = errorCode
		}

		responseData.Detail = err.Error()
		c.JSON(statusCode, responseData)
	}
}

// handleSuccess
func handleSuccess(c *gin.Context, data interface{}) {
	var responseData ResponseData
	responseData.ErrorMsg = codeToERRMsgMap[CODE_SUCCESS]
	responseData.ErrorCode = CODE_SUCCESS
	responseData.Data = data
	c.JSON(http.StatusOK, responseData)
}

type ResponseData struct {
	ErrorMsg  string      `json:"error_msg"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data,omitempty"`
	Detail    interface{} `json:"detail,omitempty"`
}
