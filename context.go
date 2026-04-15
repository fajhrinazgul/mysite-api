package main

import (
	"net/http"

	"github.com/fajhrinazgul/mysite-api/models"
)

func getUserContext(r *http.Request) (models.User, error) {
	context := r.Context().Value(&userAuth{}).(claims)
	user, err := models.NewUserModel(models.GetDB()).GetUserByID(context.credential.UserID)
	return user, err
}
