package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ThrottleHandlerFunc func() bool

func ThrottleMiddleware(h ThrottleHandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !h() {
			context.AbortWithStatusJSON(http.StatusServiceUnavailable, map[string]string{"msg": "you request is rejected by gin throttle control, please retry later"})
		}
		return
	}
}
