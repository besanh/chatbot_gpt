package response

import (
	"net/http"
)

const (
	SUCCESS = "success"
)

type (
	ErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	ValidatorFieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func OKResponse() (int, any) {
	return http.StatusOK, map[string]any{
		"message": "SUCCESS",
		"code":    http.StatusText(http.StatusOK),
	}
}

func BadRequestMsg(msg any) (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusText(http.StatusBadRequest),
		"message": msg,
	}
}

func NotFoundMsg(msg any) (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusText(http.StatusNotFound),
		"message": msg,
	}
}

func Forbidden() (int, any) {
	return http.StatusForbidden, map[string]any{
		"error":   "Do not have permission for the request.",
		"code":    http.StatusText(http.StatusForbidden),
		"message": http.StatusText(http.StatusForbidden),
	}
}

func Unauthorized() (int, any) {
	return http.StatusUnauthorized, map[string]any{
		"error":   http.StatusText(http.StatusUnauthorized),
		"code":    http.StatusText(http.StatusUnauthorized),
		"message": http.StatusText(http.StatusUnauthorized),
	}
}

func ServiceUnavailableMsg(msg any) (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusText(http.StatusServiceUnavailable),
		"message": msg,
	}
}

func Created(data any) (int, any) {
	result := map[string]any{
		"code":    http.StatusCreated,
		"message": "SUCCESS",
		"data":    data,
	}

	return http.StatusCreated, result
}

func Pagination(data, total, limit, offset any) (int, any) {
	return http.StatusOK, map[string]any{
		"data":   data,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
}

func OK(data any) (int, any) {
	return http.StatusOK, data
}
