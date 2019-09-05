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

// Plugin for gin
func (proxy *Proxy) Plugin(httpProxy *gin.Engine) {
	httpProxy.Use(middleware(proxy))
}
