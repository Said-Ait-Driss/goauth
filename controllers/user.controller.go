package controllers

import (
	// "context"
	// "fmt"
	// "http/net"

	"context"
	"log"
	"net/http"
	"time"

	"github.com/Said-Ait-Driss/go-auth/database"
	"github.com/Said-Ait-Driss/go-auth/helpers"
	"github.com/Said-Ait-Driss/go-auth/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.Collection(database.Client, "user")

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, foundUserPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(foundUserPassword), []byte(userPassword))
	valide := true
	msg := ""

	if err != nil {
		valide = false
		msg = "email or password incorect !"
	}

	return valide, msg

}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationError := validate.Struct(&user)
		defer cancel()
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError.Error()})
			return
		}

		count, erro := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if erro != nil {
			log.Panic(erro)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": erro.Error(),
			})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		count, erro = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if erro != nil {
			log.Panic(erro)
			c.JSON(http.StatusInternalServerError, gin.H{"error": erro.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone or email already exists !"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = primitive.NewObjectID().Hex()

		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)

		user.Token = &token
		user.Refresh_token = &refreshToken

		nbreInserted, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while inserting user"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, nbreInserted)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isValidPassword, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !isValidPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, user_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func GetUsers() {

}
