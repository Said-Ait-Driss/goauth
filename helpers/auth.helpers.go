package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func checkUserType(c *gin.Context, role string) error {
	user_type := c.GetString("user_type")

	if user_type != role {
		return errors.New("Unauthorized")
	}
	return nil
}

func MatchUserTypeToUid(c *gin.Context, user_id string) error {
	user_type := c.GetString("user_type")
	uid := c.GetString("uid")

	if user_type == "USER" && uid != user_id {
		return errors.New("Unauthorized")
	}
	return checkUserType(c, user_type)
}
