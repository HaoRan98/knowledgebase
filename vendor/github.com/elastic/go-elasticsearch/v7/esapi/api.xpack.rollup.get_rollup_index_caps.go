// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 7.8.0: DO NOT EDIT

package esapi

import (
	"context"
	"net/http"
	"strings"
)

func newRollupGetRollupIndexCapsFunc(t Transport) RollupGetRollupIndexCaps {
	return func(index string, o ...func(*RollupGetRollupIndexCapsRequest)) (*Response, error) {
		var r = RollupGetRollupIndexCapsRequest{Index: index}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// RollupGetRollupIndexCaps - Returns the rollup capabilities of all jobs inside of a rollup index (e.g. the index where rollup data is stored).
//
// This API is experimental.
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/rollup-get-rollup-index-caps.html.
//
type RollupGetRollupIndexCaps func(index string, o ...func(*RollupGetRollupIndexCapsRequest)) (*Response, error)

// RollupGetRollupIndexCapsRequest configures the Rollup Get Rollup Index Caps API request.
//
type RollupGetRollupIndexCapsRequest struct {
	Index string

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r RollupGetRollupIndexCapsRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "GET"

	path.Grow(1 + len(r.Index) + 1 + len("_rollup") + 1 + len("data"))
	path.WriteString("/")
	path.WriteString(r.Index)
	path.WriteString("/")
	path.WriteString("_rollup")
	path.WriteString("/")
	path.WriteString("data")

	params = make(map[string]string)

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f RollupGetRollupIndexCaps) WithContext(v context.Context) func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f RollupGetRollupIndexCaps) WithPretty() func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f RollupGetRollupIndexCaps) WithHuman() func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f RollupGetRollupIndexCaps) WithErrorTrace() func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f RollupGetRollupIndexCaps) WithFilterPath(v ...string) func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f RollupGetRollupIndexCaps) WithHeader(h map[string]string) func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
//
func (f RollupGetRollupIndexCaps) WithOpaqueID(s string) func(*RollupGetRollupIndexCapsRequest) {
	return func(r *RollupGetRollupIndexCapsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
