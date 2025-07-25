package middleware

import (
	"RestGoTest/src/config"
	"RestGoTest/src/constant"
	"RestGoTest/src/helper"
	"RestGoTest/src/helper/service_errors"
	"RestGoTest/src/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	var tokenUsecase = services.NewTokenUsecase(cfg)

	return func(c *gin.Context) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.GetHeader(constant.AuthorizationHeaderKey)
		token := strings.Split(auth, " ")
		if auth == "" {
			err = &service_errors.ServiceError{EndUserMessage: constant.TokenRequired}
		} else if len(token) == 1 {
			token = append([]string{"Bearer"}, token[0])
		} else {
			claimMap, err = tokenUsecase.GetClaims(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &service_errors.ServiceError{EndUserMessage: constant.TokenExpired}
				default:
					err = &service_errors.ServiceError{EndUserMessage: constant.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.GenerateBaseResponseWithError(
				nil, false, helper.AuthError, err,
			))
			return
		}

		c.Set(constant.UserIdKey, claimMap[constant.UserIdKey])
		c.Set(constant.FirstNameKey, claimMap[constant.FirstNameKey])
		c.Set(constant.LastNameKey, claimMap[constant.LastNameKey])
		c.Set(constant.UsernameKey, claimMap[constant.UsernameKey])
		c.Set(constant.EmailKey, claimMap[constant.EmailKey])
		c.Set(constant.MobileNumberKey, claimMap[constant.MobileNumberKey])
		c.Set(constant.RolesKey, claimMap[constant.RolesKey])
		c.Set(constant.ExpireTimeKey, claimMap[constant.ExpireTimeKey])

		c.Next()
	}
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Keys) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
			return
		}
		rolesVal := c.Keys[constant.RolesKey]
		fmt.Println(rolesVal)
		if rolesVal == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
			return
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}

		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
	}
}
