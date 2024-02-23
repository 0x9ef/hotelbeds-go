// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrorCode represents code of HotelBeds error.
type ErrorCode string

const (
	ErrorCodeConfiguration  ErrorCode = "CONFIGURATION_ERROR"
	ErrorCodeSystem         ErrorCode = "SYSTEM_ERROR"
	ErrorCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrorCodeInvalidData    ErrorCode = "INVALID_DATA"
	ErrorCodeProduct        ErrorCode = "PRODUCT_ERROR"
)

var (
	// General error.
	ErrNoHealthyUpstream = errors.New("no healthy upstead")
	ErrExternal          = errors.New("external error")
	ErrRateLimitExceeded = errors.New("rate limits exceeded")
	ErrQuotaExceeded     = errors.New("quota exceeded")
	ErrConfiguration     = errors.New("configuration error")
	ErrSystem            = errors.New("system error")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidData       = errors.New("invalid data")

	// Product errors.
	ErrAllotmentExceeded                                = errors.New("allotment exceeded")
	ErrInsufficientAllotment                            = errors.New("insufficient allotment")
	ErrPriceHasIncreased                                = errors.New("price has increased")
	ErrPriceHasChanged                                  = errors.New("price has changed")
	ErrStopSales                                        = errors.New("stop sales")
	ErrBookingDoesNotExist                              = errors.New("booking does not exist")
	ErrBookingCannotBeCanceledAfterCheckIn              = errors.New("cannot cancel a booking after the check-in")
	ErrBookingCannotBeCanceledOrModifiedWhenCheckInPast = errors.New("cannot cancel/modify a booking which has a check-in date in the past")
	ErrBookingCannotBeAmended                           = errors.New("this booking cannot be amended")
	ErrBookingConfirmationError                         = errors.New("booking confirmation error")
	ErrMinimumStayViolated                              = errors.New("minimum stay violated")
	ErrCodeIsInvalid                                    = errors.New("code is invalid")
	ErrPleaseDoNotTryAgain                              = errors.New("please do not retry again")
	ErrHotelDoesNotAllowCancellation                    = errors.New("hotel does not allow cancellations")
	ErrReservationUnreachable                           = errors.New("reservation does not exist or the agency does not access")
	ErrUndefined                                        = errors.New("undefined error")
)

type Error struct {
	Audit   *AuditData `json:"auditData"`
	Code    ErrorCode  `json:"code"`
	Message string     `json:"message"`
	// Our internal variables.
	StatusCode  int  `json:"-"`
	IsRetryable bool `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%s,statusCode=%d,message=%s", e.Code, e.StatusCode, e.Message)
}

// IsErrorCode checks if error contains specified code.
func IsErrorCode(err error, code ErrorCode) bool {
	if err, ok := err.(*Error); ok {
		return err.Code == code
	}
	return false
}

// IsErrorRetryable checks if error is retryable.
func IsErrorRetryable(err error) bool {
	if err, ok := err.(*Error); ok {
		return err.IsRetryable
	}
	return false
}

type shortError struct {
	Error string `json:"error"`
}

func decodeError(resp *http.Response) error {
	var shortErr shortError
	if err := json.NewDecoder(resp.Body).Decode(&shortErr); err == nil {
		isRetryable, _ := isRetryableError[decodeErrorMessage(shortErr.Error)]
		return &Error{
			Message:     shortErr.Error,
			StatusCode:  resp.StatusCode,
			IsRetryable: isRetryable,
		}
	}

	var longErr Error
	if err := json.NewDecoder(resp.Body).Decode(&longErr); err == nil {
		isRetryable, _ := isRetryableError[decodeErrorMessage(shortErr.Error)]
		return &Error{
			Audit:       longErr.Audit,
			Code:        longErr.Code,
			Message:     longErr.Message,
			StatusCode:  resp.StatusCode,
			IsRetryable: isRetryable,
		}
	}

	return ErrUndefined
}

var (
	isRetryableError = map[error]bool{
		ErrRateLimitExceeded: true,
		ErrQuotaExceeded:     true,
	}
)

func decodeErrorMessage(msg string) error {
	switch {
	case errorContains(msg, ErrExternal):
		return ErrExternal
	case errorContains(msg, ErrRateLimitExceeded):
		return ErrRateLimitExceeded
	case errorContains(msg, ErrQuotaExceeded):
		return ErrQuotaExceeded
	case errorContains(msg, ErrConfiguration):
		return ErrConfiguration
	case errorContains(msg, ErrSystem):
		return ErrSystem
	case errorContains(msg, ErrInvalidRequest):
		return ErrInvalidRequest
	case errorContains(msg, ErrInvalidData):
		return ErrInvalidData
	case errorContains(msg, ErrAllotmentExceeded):
		return ErrAllotmentExceeded
	case errorContains(msg, ErrInsufficientAllotment):
		return ErrInsufficientAllotment
	case errorContains(msg, ErrPriceHasIncreased):
		return ErrPriceHasIncreased
	case errorContains(msg, ErrPriceHasChanged):
		return ErrPriceHasChanged
	case errorContains(msg, ErrStopSales):
		return ErrStopSales
	case errorContains(msg, ErrBookingDoesNotExist):
		return ErrBookingDoesNotExist
	case errorContains(msg, ErrBookingCannotBeCanceledAfterCheckIn):
		return ErrBookingCannotBeCanceledAfterCheckIn
	case errorContains(msg, ErrBookingCannotBeCanceledOrModifiedWhenCheckInPast):
		return ErrBookingCannotBeCanceledOrModifiedWhenCheckInPast
	case errorContains(msg, ErrBookingCannotBeAmended):
		return ErrBookingCannotBeAmended
	case errorContains(msg, ErrBookingConfirmationError):
		return ErrBookingConfirmationError
	case errorContains(msg, ErrMinimumStayViolated):
		return ErrMinimumStayViolated
	case errorContains(msg, ErrCodeIsInvalid):
		return ErrCodeIsInvalid
	case errorContains(msg, ErrPleaseDoNotTryAgain):
		return ErrPleaseDoNotTryAgain
	case errorContains(msg, ErrHotelDoesNotAllowCancellation):
		return ErrHotelDoesNotAllowCancellation
	case errorContains(msg, ErrReservationUnreachable):
		return ErrReservationUnreachable
	default:
		return errors.New("undefined error")
	}
}

func errorContains(s string, err error) bool {
	return strings.Contains(strings.ToLower(s), err.Error())
}

type ValidationError struct {
	Required  bool
	FieldName string
	Min       int
	Max       int
	Allow     []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field=%s,required=%t,min=%d,max=%d,allow=[%s]", e.FieldName, e.Required, e.Min, e.Max, strings.Join(e.Allow, ","))
}
