package main

import "github.com/gin-gonic/gin"

func postControllerWithAuth(group *gin.RouterGroup) {
	group.POST("", createPostHandler())
}
