package service

import (
	"fmt"
	"net/http"
)

type Error struct {
	HttpCode int    `json:"-"`
	Code     int    `json:"errno"`
	Reason   string `json:"errmsg"`
}

func NewError(httpCode int, code int, reason string) *Error {
	return &Error{
		HttpCode: httpCode,
		Code:     code,
		Reason:   reason,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("<Error(http:%d code:%d reason:%s)>", e.HttpCode, e.Code, e.Reason)
}

//error object define
var (
	ErrUnauthorized                = &Error{HttpCode: http.StatusUnauthorized, Code: 101, Reason: "unauthorized"}
	ErrNotAllowed                  = &Error{HttpCode: http.StatusForbidden, Code: 106, Reason: "not allowed"}
	ErrInvalidParam                = &Error{HttpCode: http.StatusBadRequest, Code: 107, Reason: "invalid param"}
	ErrResourceNotFound            = &Error{HttpCode: http.StatusNotFound, Code: 108, Reason: "resouece not found"}
	ErrLessDownloadFlowQuota       = &Error{HttpCode: http.StatusForbidden, Code: 110, Reason: "less download flow quota"}
	ErrLessStorageQuota            = &Error{HttpCode: http.StatusForbidden, Code: 112, Reason: "less storage quota"}
	ErrUploadBusy                  = &Error{HttpCode: http.StatusForbidden, Code: 115, Reason: "upload busy"}
	ErrAppIsNotActived             = &Error{HttpCode: http.StatusForbidden, Code: 119, Reason: "App is not actived"}
	ErrInternalServiceAccessFailed = &Error{HttpCode: http.StatusInternalServerError, Code: 120, Reason: "internal service access failed"}
	ErrNodeAccessFailed            = &Error{HttpCode: http.StatusInternalServerError, Code: 147, Reason: "node access failed"}
	ErrNodeNotFound                = &Error{HttpCode: http.StatusInternalServerError, Code: 148, Reason: "node not found"}
	ErrZoneAccessFailed            = &Error{HttpCode: http.StatusInternalServerError, Code: 150, Reason: "zone access failed"}
	ErrDBFailed                    = &Error{HttpCode: http.StatusInternalServerError, Code: 151, Reason: "db failed"}
	ErrStorageNotFound             = &Error{HttpCode: http.StatusInternalServerError, Code: 152, Reason: "storage not found"}
	ErrNoStorage                   = &Error{HttpCode: http.StatusInternalServerError, Code: 153, Reason: "no storage space"}
	ErrPortAllocationFailed        = &Error{HttpCode: http.StatusInternalServerError, Code: 154, Reason: "port allocation failed"}
)

func InternalError(err error) *Error {
	return &Error{
		HttpCode: http.StatusInternalServerError,
		Code:     100,
		Reason:   err.Error(),
	}
}

func InvalidParamError(reason string) *Error {
	return &Error{http.StatusBadRequest, 106, reason}
}
