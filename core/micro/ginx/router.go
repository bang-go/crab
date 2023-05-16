package ginx

import "github.com/gin-gonic/gin"

type RouterGroup interface {
	RouterHandle(httpMethod, relativePath string, handlers ...gin.HandlerFunc)
}

type routerGroupEntity struct {
	group *gin.RouterGroup
}

// RouterHandle Handle 路由
func (r *routerGroupEntity) RouterHandle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	r.group.Handle(httpMethod, relativePath, handlers...)
}
