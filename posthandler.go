package main

import (
	"fmt"
	"math"
	"net/http"
	"path/filepath"
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
			"data":           post,
			"execution_time": float64(duration.Microseconds()) / 1000,
		})
	}
}

// getAllPostHandler
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

// postDetailHandler
func postDetailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		slug := ctx.Param("slug")
		post, err := models.NewPostModel(models.GetDB()).GetPostBySlug(slug)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Post not found.",
			})
			return
		}

		duration := time.Since(start)

		ctx.JSON(http.StatusOK, gin.H{
			"status":         "success",
			"message":        "",
			"execution_time": float64(duration.Microseconds()) / 1000,
			"data":           post,
		})
	}
}

// uploadImageHandler
func uploadImageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("image")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan>"})
			return
		}

		// validation format
		extension := filepath.Ext(file.Filename)
		if extension != ".jpg" && extension != ".png" && extension != ".jpeg" && extension != ".webp" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Format file not support",
			})
			return
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		savePath := filepath.Join("uploads/posts", fileName)

		// save file
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to save file.",
			})
			return
		}

		fileURL := fmt.Sprintf("/%s", savePath)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Image upload successfully",
			"url":     fileURL,
		})
	}
}

// function for edit
func postEditHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		slugQuery := ctx.Param("slug")

		postInstance, err := models.NewPostModel(models.GetDB()).GetPostBySlug(slugQuery)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		_, err = getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "you not logged as a admin.",
			})
			return
		}

		var payload PostEditPayload
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
		// fmt.Println(payload.IsFeatured)

		if payload.Title != nil {
			postInstance.Title = *payload.Title
		}
		if payload.Content != nil {
			postInstance.Content = *payload.Content
		}
		if payload.Status != nil {
			postInstance.Status = *payload.Status
		}
		if payload.IsFeatured != nil {
			postInstance.IsFeatured = *payload.IsFeatured
		}
		if err := models.GetDB().Save(&postInstance).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		duration := time.Since(start)

		ctx.JSON(http.StatusOK, gin.H{
			"status":         "success",
			"message":        "success update post",
			"data":           postInstance,
			"execution_time": float64(duration.Microseconds()) / 1000,
		})
	}
}

// deletePostHandler
func deletePostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Permission danied",
			})
			return
		}

		slug := ctx.Param("slug")

		post, err := models.NewPostModel(models.GetDB()).GetPostBySlug(slug)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Post not found",
			})
			return
		}

		if err = models.GetDB().Delete(&post).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}
