package main

import (
	"fmt"
	"net/http"

	"github.com/fajhrinazgul/mysite-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func createPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "you not logged as a admin.",
			})
			return
		}

		var payload PostPayload
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		validate = validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(payload)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err.Error())
				return
			}

			var validations []Validation

			for _, err := range err.(validator.ValidationErrors) {
				validations = append(validations, Validation{
					Field: err.Field(),
					Tag:   err.Tag(),
				})
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "validation error",
				"errors":  validations,
			})
			return
		}

		post := models.Post{
			Title:   payload.Title,
			Content: payload.Content,
			Status:  payload.Status,
			// Tags:    payload.Tags,
		}
		err = models.NewPostModel(models.GetDB()).CreatePost(&post)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "success create new post",
			"post":    post,
		})
	}
}
