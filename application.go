/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush role plugin]
 */

package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
)

// Proxy http proxy
type Proxy struct {
	bulrush.PNBase
	Host  string
	Match string
	Map   func(string) string
}

// Plugin for gin
func (proxy *Proxy) Plugin() bulrush.PNRet {
	return func(httpProxy *gin.Engine) {
		httpProxy.Use(func(c *gin.Context) {
			// MatchString url
			if proxy.Match != "" {
				r, _ := regexp.Compile(proxy.Match)
				if !r.MatchString(c.Request.RequestURI) {
					c.Next()
					return
				}
			}
			// Parse body
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			}

			// Reassgin to body
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
			// Map path
			var url string
			if proxy.Map != nil {
				url = proxy.Map(c.Request.RequestURI)
			} else {
				url = c.Request.RequestURI
			}
			url = fmt.Sprintf("%s%s", proxy.Host, url)
			proxyReq, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
			proxyReq.Header = make(http.Header)
			// Reassgin headers
			for key, val := range c.Request.Header {
				proxyReq.Header[key] = val
			}

			client := &http.Client{}
			resp, err := client.Do(proxyReq)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			bodyContent, _ := ioutil.ReadAll(resp.Body)
			c.Writer.Write(bodyContent)
			for key, val := range resp.Header {
				c.Writer.Header()[key] = val
			}
			c.Abort()
		})
	}
}
