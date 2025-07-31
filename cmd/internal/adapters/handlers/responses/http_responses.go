package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HTTPStatus(context *gin.Context, status int, message string, err error) {
	context.JSON(status, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

func BadRequest(context *gin.Context, message string, err error) {
	HTTPStatus(context, http.StatusBadRequest, message, err)
}

func InternalServerError(context *gin.Context, message string, err error) {
	HTTPStatus(context, http.StatusInternalServerError, message, err)
}

func Unauthorized(context *gin.Context, message string, err error) {
	HTTPStatus(context, http.StatusUnauthorized, message, err)
}

func NotFound(context *gin.Context, message string, err error) {
	HTTPStatus(context, http.StatusNotFound, message, err)
}
