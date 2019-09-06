// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	opath "path"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func resolve(path string, proxy *Proxy) string {
	url := path
	if proxy.URL != "" {
		url = proxy.URL
		r, _ := regexp.Compile("^http")
		if !r.MatchString(url) {
			if proxy.Host != "" {
				u, _ := nurl.Parse(proxy.Host)
				u.Path = opath.Join(u.Path, url)
				url = u.String()
			}
			url = ""
		}
		return url
	}
	if proxy.Map != nil {
		url = proxy.Map(url)
	}
	url = strings.Split(url, "?")[0]
	if proxy.Host != "" {
		u, _ := nurl.Parse(proxy.Host)
		u.Path = opath.Join(u.Path, url)
		return u.String()
	}
	return ""
}

func request(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func middleware(proxy *Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// MatchString url
		if proxy.Match != "" {
			r, _ := regexp.Compile(proxy.Match)
			if !r.MatchString(c.Request.RequestURI) {
				c.Next()
				return
			}
		}

		// Resolve url from cfg
		url := resolve(c.Request.URL.Path, proxy)
		if url == "" {
			RushLogger.Debug("resolve empty url %v", url)
			c.Next()
			return
		}

		// Parse body
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		}

		// Reassgin to body
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		url = fmt.Sprintf("%s?%s", url, c.Request.URL.RawQuery)

		proxyReq, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
		proxyReq.Header = make(http.Header)

		// Reassgin headers
		for key, val := range c.Request.Header {
			proxyReq.Header[key] = val
		}

		// replace host
		if proxy.Host != "" {
			proxyReq.Header["Host"] = []string{string(proxy.Host[strings.Index(proxy.Host, "://")+3:])}
		}

		// custom proxyReq if you need
		if proxy.req != nil {
			proxyReq = proxy.req(proxyReq, c.Request)
		}

		resp, err := request(proxyReq)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		c.Status(resp.StatusCode)

		for key, val := range resp.Header {
			c.Writer.Header()[key] = val
		}

		// custom proxyRes if you need
		if proxy.res != nil {
			proxy.res(c.Writer, resp)
		}
		c.Writer.Write(bodyContent)
		c.Abort()
	}
}
