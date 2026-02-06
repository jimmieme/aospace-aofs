// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"aofs/internal/log4bp"
	"time"

	"github.com/gin-gonic/gin"
)

var httpLogger = log4bp.New("", gin.Mode())

func LoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			start += "?" + c.Request.URL.RawQuery
		}

		reqId := c.GetHeader("Request-Id")
		if reqId == "" {
			reqId = c.GetHeader("X-Request-Id")
		}

		t0 := time.Now()
		c.Next()

		latency := time.Since(t0).Milliseconds()
		status := c.Writer.Status()
		evt := httpLogger.LogI()
		if status >= 500 {
			evt = httpLogger.LogE()
		} else if status >= 400 {
			evt = httpLogger.LogW()
		}

		evt.Str("method", c.Request.Method).
			Str("path", start).
			Int("status", status).
			Int64("latency_ms", latency).
			Str("client_ip", c.ClientIP()).
			Str("request_id", reqId).
			Str("proto", c.Request.Proto).
			Str("user_agent", c.Request.UserAgent()).
			Int("size", c.Writer.Size())

		if errMsg := c.Errors.String(); errMsg != "" {
			evt.Str("errors", errMsg)
		}

		evt.Msg("http request")
	}
}
