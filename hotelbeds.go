// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/0x9ef/clientx"
	"golang.org/x/exp/constraints"
)

type (
	API struct {
		*clientx.API
		options   *Options
		apiKey    string
		apiSecret string
	}

	Client interface {
		ContentClient
		BookingClient
	}

	Option  func(*Options)
	Options struct {
		DefaultHeaders http.Header
		Limit          *clientx.OptionRateLimit
		Retry          *clientx.OptionRetry
	}
)

var _ Client = (*API)(nil)

// New returns new API with provided apiKey, apiSecret, applies all options.
func New(apiKey, apiSecret string, opts ...Option) *API {
	api := &API{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	api.options = &options
	api.API = clientx.NewAPI(api.options.toClientxOptions()...)
	return api
}

func (opts *Options) toClientxOptions(options ...Option) []clientx.Option {
	clientxOptions := make([]clientx.Option, 0, len(options)+1)
	clientxOptions = append(clientxOptions, clientx.WithBaseURL("https://api.test.hotelbeds.com"))
	if opts.Limit != nil {
		clientxOptions = append(clientxOptions,
			clientx.WithRateLimit(opts.Limit.Limit, opts.Limit.Burst, opts.Limit.Per))
	}
	if opts.Retry != nil {
		clientxOptions = append(clientxOptions,
			clientx.WithRetry(opts.Retry.MaxAttempts, opts.Retry.MinWaitTime, opts.Retry.MaxWaitTime, opts.Retry.Fn, opts.Retry.Conditions...))
	}
	return clientxOptions
}

func (api *API) buildHeaders() http.Header {
	return http.Header{
		"Accept":          []string{"application/json"},
		"Accept-Encoding": []string{"application/json"},
		"Content-Type":    []string{"application/json"},
		"Api-key":         []string{api.apiKey},
		"X-Signature":     []string{api.hashSignature()},
	}
}

func (api *API) hashSignature() string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s%s%d", api.apiKey, api.apiSecret, time.Now().Unix())))
	return hex.EncodeToString(hasher.Sum(nil))
}

func joinInts[T constraints.Integer](values []int) string {
	var sb strings.Builder
	for i := range values {
		sb.WriteString(strconv.Itoa(values[i]) + ",")
	}
	return sb.String()[:sb.Len()-1]
}
