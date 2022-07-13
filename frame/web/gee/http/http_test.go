package http

import (
	"awesome/frame/web/gee/base"
	"awesome/frame/web/gee/middleware"
	"log"
	"net/http"
	"testing"
	"time"
)

/*
	curl "http://localhost:6789/login?username=zjf&password=hello" -H "host:static-zjf.snssdk.com" -X POST
*/
func TestHttpServer(t *testing.T) {
	GetHelloServer().init()

}

// base server
func TestGeeServer(t *testing.T) {
	r := base.New()

	r.GET("/", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1> Hello Gee </h1>")
	})

	r.GET("/hello", func(c *base.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *base.Context) {
		c.JSON(http.StatusOK, base.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":6789")
}

// wild route server
func TestGeeServerWild(t *testing.T) {
	r := base.New()
	r.GET("/", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *base.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *base.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *base.Context) {
		c.JSON(http.StatusOK, base.H{"filepath": c.Param("filepath")})
	})

	r.Run(":6789")
}

// group route server
func TestGeeServerGroup(t *testing.T) {
	r := base.New()
	r.GET("/index", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	{
		v1 := r.Group("/v1")
		v1.GET("/hello", func(c *base.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	{
		v2 := r.Group("/v2")
		v2.GET("/hello/:name", func(c *base.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *base.Context) {
			c.JSON(http.StatusOK, base.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":6789")
}

func onlyForV2() base.HandlerFunc {
	return func(c *base.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func TestGeeServerWithMid(t *testing.T) {
	r := base.New()
	r.Use(middleware.Logger())

	r.GET("/", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	{
		v2 := r.Group("/v2")
		v2.Use(onlyForV2()) // v2 group middleware
		v2.GET("/hello/:name", func(c *base.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":6789")

}
