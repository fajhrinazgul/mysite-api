package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/fajhrinazgul/mysite-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// createPostHandler is handler to create new post
func createPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

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
		fmt.Println(payload.IsFeatured)

		post := models.Post{
			Title:      payload.Title,
			Content:    payload.Content,
			Status:     payload.Status,
			IsFeatured: payload.IsFeatured,
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

		duration := time.Since(start)

		ctx.JSON(http.StatusCreated, gin.H{
			"status":         "success",
			"message":        "success create new post",
			"post":           post,
			"execution_time": float64(duration.Microseconds()) / 1000,
		})
	}
}

func getAllPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		status := ctx.DefaultQuery("status", "published")

		offset := (page - 1) * limit

		postModel := models.NewPostModel(models.GetDB())
		posts, total, err := postModel.GetPagedPosts(status, limit, offset)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch data",
			})
			return
		}

		totalPage := int(math.Ceil(float64(total) / float64(limit)))

		duration := time.Since(start)

		ctx.JSON(http.StatusOK, gin.H{
			"status":         "success",
			"message":        "Posts retrieve successfully",
			"data":           posts,
			"execution_time": float64(duration.Microseconds()) / 1000,
			"pagination": gin.H{
				"current_page": page,
				"limit":        limit,
				"total_items":  total,
				"total_page":   totalPage,
				"has_next":     page < totalPage,
				"has_prev":     page > 1,
			},
			"queries": gin.H{
				"limit":  limit,
				"offset": offset,
				"status": status,
			},
		})
	}
}
