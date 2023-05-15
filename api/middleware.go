package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chensheep/simple-bank-backend/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is not provided")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	authorizationType := strings.ToLower(jwtToken[0])
	if authorizationType != authorizationTypeBearer {
		return "", fmt.Errorf("unsupported authorization type %s", authorizationType)
	}

	return jwtToken[1], nil
}
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := extractBearerToken(c.GetHeader(authorizationHeaderKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse(err))
			c.Abort()
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}
