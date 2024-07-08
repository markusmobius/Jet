package http

import (
	"io"
	"jet/API"
	"jet/API/Keygen"
	"jet/StorageEngines"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Serve(hostAndPort string, storage StorageEngines.Storage, keyGenerator Keygen.KeyGenerator) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(textPlain)
	r.Use(cors.Default())
	r.GET("/", func(c *gin.Context) {
		content, err := API.AllKeys(storage)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, content)
	})
	r.PUT("/:key", func(c *gin.Context) {
		key := c.Param("key")
		content, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		_, err = API.Set(storage, key, string(content))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Writer.WriteString(key)
	})
	r.POST("/", func(c *gin.Context) {
		content, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		key, err := API.Store(storage, keyGenerator, string(content))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Writer.WriteString(key)
	})
	r.GET("/:key", func(c *gin.Context) {
		key := c.Param("key")
		content, err := API.Get(storage, key)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		c.Writer.WriteString(content)
	})
	r.DELETE("/:key", func(c *gin.Context) {
		key := c.Param("key")

		err := API.Delete(storage, key)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		c.Status(204)
	})

	return r.Run(hostAndPort)
}

func textPlain(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Next()

}
