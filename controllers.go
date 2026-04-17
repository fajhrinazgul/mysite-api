package main

import "github.com/gin-gonic/gin"

func postControllerWithAuth(group *gin.RouterGroup) {
	group.POST("/posts", createPostHandler())
}

func postControllerNoAuth(group *gin.RouterGroup) {
	group.GET("/posts", getAllPostHandler())
}
