package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fajhrinazgul/mysite-api/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "runserver":
			handleServer()
		case "createsuperuser":
			handleCreateSuperuser()
		case "migrate":
			models.Migrate(models.GetDB())
			fmt.Println("Selesai migrate")
		default:
			fmt.Printf("Command %s tidak dikenal.\n", command)
		}
	}
}

func handleCreateSuperuser() {
	fs := flag.NewFlagSet("createsuperuser", flag.ExitOnError)
	firstName := fs.String("firstname", "", "type your first name")
	lastName := fs.String("lastname", "", "type your last name")
	username := fs.String("username", "", "type your username")
	email := fs.String("email", "", "type your email")
	password := fs.String("password", "", "type your password")

	if len(os.Args) > 2 {
		fs.Parse(os.Args[2:])
	} else {
		fmt.Println("User: createsuperuser -firstname=... -lastname=... -username=... -email=... --password=...")
		return
	}

	if *firstName != "" && *lastName != "" && *username != "" && *email != "" && *password != "" {
		db := models.GetDB()
		err := models.NewUserModel(db).CreateUser(&models.User{
			FirstName: *firstName,
			LastName:  *lastName,
			Username:  *username,
			Email:     *email,
			Password:  encryptionPassword(*password),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("Successfully create new user.")
	}
}

func handleServer() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Static("/media", "./media")
	router.Static("/uploads", "./uploads")

	c := cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	})

	router.Use(c)

	router.POST("/get-token", getTokenHandler())
	router.POST("/uploads", uploadImageHandler())

	postGroupWithAuth := router.Group("")
	postGroupWithAuth.Use(authMiddleware())
	postControllerWithAuth(postGroupWithAuth)

	postGroupNoAuth := router.Group("")
	postControllerNoAuth(postGroupNoAuth)

	router.Run(":8000")
}
