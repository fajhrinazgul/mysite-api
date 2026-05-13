package main

import (
	"net/http"

	"github.com/fajhrinazgul/mysite-api/models"
	"github.com/gin-gonic/gin"
)

func getTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		user, err := models.NewUserModel(models.GetDB()).GetUserByUsername(payload.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Username or password is incorrect!",
			})
			return
		}

		if !decryptionPassword(user.Password, payload.Password) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Username or password is incorrect!",
			})
			return
		}

		token, err := getToken(credential{UserID: user.ID})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully generate new token.",
			"token":   token,
		})
	}
}
