package cache

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/cache/gredis"
	"net/http"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		if gredis.GRedis == nil || !gredis.GRedis.IsConnected() {
			logger.Warn("Redis cache is not available")
			c.Next()
			return
		}

		key := c.Request.URL.RequestURI()
		if c.Request.Method != "GET" {
			c.Next()

			if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
				gredis.GRedis.Remove(key)
			}

			return
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		value := gredis.GRedis.Get(key)
		var data map[string]interface{}
		json.Unmarshal(value, &data)

		if value != nil {
			c.JSON(http.StatusOK, data)

			c.Abort()
			return
		}

		c.Next()

		statusCode := w.Status()
		if statusCode == http.StatusOK {
			gredis.GRedis.Set(key, w.body.Bytes())
		}
	}
}
