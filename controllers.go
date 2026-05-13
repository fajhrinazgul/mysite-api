package main

import "github.com/gin-gonic/gin"

func postControllerWithAuth(group *gin.RouterGroup) {
	group.POST("/posts", createPostHandler())
	group.PUT("/posts/:slug", postEditHandler())
	group.DELETE("/posts/:slug", deletePostHandler())
}

func postControllerNoAuth(group *gin.RouterGroup) {
	group.GET("/posts", getAllPostHandler())
	group.GET("/posts/:slug", postDetailHandler())
}
