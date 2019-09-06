// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Proxy http proxy
type Proxy struct {
	URL   string
	Host  string
	Match string
	Map   func(string) string
	req   func(*http.Request, *http.Request) *http.Request
	res   func(gin.ResponseWriter, *http.Response)
}

// AddOptions defined add option
func (proxy *Proxy) AddOptions(opts ...Option) *Proxy {
	for _, v := range opts {
		v.apply(proxy)
	}
	return proxy
}

// Plugin for gin
func (proxy *Proxy) Plugin(httpProxy *gin.Engine) {
	httpProxy.Use(middleware(proxy))
}
