/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.2.0
 * Contact: support@athonet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"
	"time"
	"io/ioutil"
	"bytes"

	"gitlab.lan.athonet.com/core/netconfd/logger"
)

//Logger wrapper for logging before inner HTTP serving
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		body := ""
		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			if err != nil {
				logger.Log.Warning("Failed to process request body")
				return
			}
		}
		if r.Method != "GET" {
			logger.Log.Notice("T:" + time.Since(start).String() + " - API:" + r.Method + ":" + r.URL.Path + ":" + string(body))
		} else {
			logger.Log.Info("T:" + time.Since(start).String() + " - API:" + r.Method + ":" + r.URL.Path + ":" + string(body))
		}
	})
}
