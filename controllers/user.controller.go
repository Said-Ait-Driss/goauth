package controllers

import (
	// "context"
	// "fmt"
	// "http/net"

	"github.com/Said-Ait-Driss/go-auth/database"
	// "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.Collection(database.Client,"user")

var validate = validator.New()

func HashPassword() {

}

func VerifyPassword() {

}

func Signup() {

}

func Login() {

}

func GetUser() {

}

func GetUsers() {

}
