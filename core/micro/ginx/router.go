package ginx

import "github.com/gin-gonic/gin"

type RouterGroup interface {
	RouterHandle(httpMethod, relativePath string, handlers ...gin.HandlerFunc)
}

type RouterGroupWrapper struct {
	Group *gin.RouterGroup
}

// RouterHandle Handle 路由
func (r *RouterGroupWrapper) RouterHandle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	r.Group.Handle(httpMethod, relativePath, handlers...)
}
