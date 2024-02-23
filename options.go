// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"net/http"
	"time"

	"github.com/0x9ef/clientx"
)

func WithRetry(maxAttempts int, minWaitTime, maxWaitTime time.Duration, f clientx.RetryFunc, conditions ...clientx.RetryCond) Option {
	return func(o *Options) {
		o.Retry = &clientx.OptionRetry{
			MaxAttempts: maxAttempts,
			MinWaitTime: minWaitTime,
			MaxWaitTime: maxWaitTime,
			Conditions:  ([]clientx.RetryCond)(conditions),
			Fn:          (clientx.RetryFunc)(f),
		}
	}
}

func WithRateLimit(limit int, burst int, per time.Duration) Option {
	return func(o *Options) {
		o.Limit = &clientx.OptionRateLimit{
			Limit: limit,
			Burst: burst,
			Per:   per,
		}
	}
}

func WithHeaders(set http.Header) Option {
	return func(o *Options) {
		o.DefaultHeaders = set
	}
}
