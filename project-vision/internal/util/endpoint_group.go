package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EndpointGroup []gin.IRoutes

func (e EndpointGroup) Use(handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.Use(handlerFunc...)
	}
	return e
}

func (e EndpointGroup) Handle(s string, s2 string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.Handle(s, s2, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) Any(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.Any(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) GET(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.GET(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) POST(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.POST(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) DELETE(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.DELETE(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) PATCH(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.PATCH(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) PUT(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.PUT(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) OPTIONS(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.OPTIONS(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) HEAD(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	for _, router := range e {
		router.HEAD(s, handlerFunc...)
	}
	return e
}

func (e EndpointGroup) StaticFile(s string, s2 string) gin.IRoutes {
	for _, router := range e {
		router.StaticFile(s, s2)
	}
	return e
}

func (e EndpointGroup) Static(s string, s2 string) gin.IRoutes {
	for _, router := range e {
		router.Static(s, s2)
	}
	return e
}

func (e EndpointGroup) StaticFS(s string, system http.FileSystem) gin.IRoutes {
	for _, router := range e {
		router.StaticFS(s, system)
	}
	return e
}
