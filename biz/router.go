package biz

import (
	v1 "webhook2syslog/biz/handler/v1"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	apiv1Group := r.Group("/api/v1")
	{
		apiv1Group.POST("/webhook", v1.Webhook)
	}
}
