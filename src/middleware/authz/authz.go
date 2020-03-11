package authz

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/khoa5773/go-server/src/constant"
	"github.com/khoa5773/go-server/src/helpers"
	"gopkg.in/mgo.v2/bson"
)

func Scopes(requiredScopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hasPersonalScopes := false
		hasProjectScopes := false

		projectID := c.Param("projectID")
		projectIDs := c.MustGet("userProjectIDs")
		userScopes := c.MustGet("userScopes")
		userScope := []string{}
		basicScopes := constant.Scopes["BASIC"]

		countBasicScopes := 0
		for _, v := range basicScopes {
			_, result := helpers.CheckItemExists(requiredScopes, v)
			if result {
				countBasicScopes++
			}
		}

		if projectID == "" && countBasicScopes > 0 {
			hasPersonalScopes = true
		} else {
			for index, v := range projectIDs.([]bson.ObjectId) {
				if v.Hex() == projectID {
					userScope = userScopes.([][]string)[index]
				}
			}

			roleAccess, _ := helpers.GetRoleScopes(userScope)

			_, hasItemExisted := helpers.CheckObjectIdExists(projectIDs.([]bson.ObjectId), projectID)

			if !hasItemExisted {
				_ = c.Error(errors.New("no permission"))
				c.Abort()
				return
			}

			countPersonalScopes := 0
			for _, v := range roleAccess.PersonalScopes {
				_, result := helpers.CheckItemExists(requiredScopes, v)
				if result {
					countPersonalScopes++
					break
				}
			}

			if len(roleAccess.PersonalScopes) > 0 && countPersonalScopes > 0 {
				hasPersonalScopes = true
			}

			countProjectScopes := 0
			for _, v := range roleAccess.ProjectScopes {
				_, result := helpers.CheckItemExists(requiredScopes, v)
				if result {
					countProjectScopes++
				}
			}

			if len(roleAccess.ProjectScopes) > 0 && countProjectScopes > 0 {
				hasProjectScopes = true
			}
		}

		if hasPersonalScopes || hasProjectScopes {
			c.Set("hasProjectScopes", hasProjectScopes)
			c.Set("hasPersonalScopes", hasPersonalScopes)
			c.Next()
			return
		}

		_ = c.Error(errors.New("no permission"))
		c.Abort()
	}
}
