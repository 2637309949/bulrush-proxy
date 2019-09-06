// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Option defined implement of option
type (
	Option func(*Proxy) interface{}
)

// option defined implement of option
func (o Option) apply(p *Proxy) *Proxy { return o(p).(*Proxy) }

// option defined implement of option
func (o Option) check(p *Proxy) interface{} { return o(p) }

// ReqOption defined req
func ReqOption(req func(*http.Request, *http.Request) *http.Request) Option {
	return Option(func(p *Proxy) interface{} {
		p.req = req
		return p
	})
}

// ResOption defined res
func ResOption(res func(gin.ResponseWriter, *http.Response)) Option {
	return Option(func(p *Proxy) interface{} {
		p.res = res
		return p
	})
}
