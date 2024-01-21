package routes

import (
	controllers "github.com/Said-Ait-Driss/go-auth/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	// r.GET("users", controllers.GetUsers())
	r.GET("users/:user_id", controllers.GetUser())
}
