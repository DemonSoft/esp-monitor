package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func success[T any](data T) Response[T] {
	return Response[T]{Data: &data, Error: nil}
}

// Error создает новый ответ с ошибкой.
func restError[T any](message string, code string) Response[T] {

	mistake := RequestError{
		Message: message,
		Code:    code,
		Caution: "Something Happened",
	}

	return Response[T]{Error: &mistake}
}

func Success[T any](c *gin.Context, data T) {
	c.IndentedJSON(http.StatusOK, success(data))
}

func RestError[T any](c *gin.Context, message string, code string) {
	c.IndentedJSON(http.StatusOK, restError[Empty](message, code))
}

func ErrorHtml[T any](c *gin.Context) {
	c.HTML(http.StatusForbidden, "error.html", nil)
}

func SomeError[T any](c *gin.Context, status int) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, RequestError{
		Message: "",
		Code:    strconv.Itoa(status),
		Caution: http.StatusText(status),
	})
}

func PanicError[T any](c *gin.Context, caution string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, RequestError{
		Message: "",
		Code:    strconv.Itoa(http.StatusInternalServerError),
		Caution: caution,
	})
}
