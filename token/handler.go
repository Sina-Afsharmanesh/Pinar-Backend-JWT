package token

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"template/log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	tokenMap      sync.Map
	cacheTokenMap sync.Map
	secret        = []byte(`P0pogHm:{"Rp%&%>~vSfY]-;7Uzlxq`)
)

func authorize(c echo.Context) error {
	//todo: fetch user role from database, see if it matches the target table and the permission if it does, pass the request
	return c.NoContent(200)
}

func generate(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}
	log.ServLogger.Info("generate - user: " + user.Id.String())

	claims := jwt.MapClaims{
		"sub":  user.Id.String(),
		"name": user.Name,
		"role": user.Role,
		"iss":  "PinatJwtService",
		"aud":  "PinarFrontend",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return c.NoContent(echo.ErrInternalServerError.Code)
	}
	log.ServLogger.Info("generate - tokenString: " + tokenString)

	tokenMap.Store(user.Id.String(), tokenString)

	return c.String(200, tokenString)
}

func verify(c echo.Context) error {
	id := c.Param("id")
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		log.ServLogger.Error("verify - authHeader is empty")
		return c.String(echo.ErrUnauthorized.Code, "Authorization is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.String(echo.ErrBadRequest.Code, "Authorization is not Bearer")
	}

	tokenString := parts[1]

	if load, ok := cacheTokenMap.Load(id); ok {
		if load.(string) == tokenString {
			return c.NoContent(200)
		}
	}

	if _, ok := tokenMap.Load(id); !ok {
		log.ServLogger.Error("verify - token not found")
		return c.String(echo.ErrUnauthorized.Code, "token not found")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			log.ServLogger.Error("verify - unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		log.ServLogger.Error("verify - " + err.Error())
		return c.String(http.StatusUnauthorized, err.Error())
	}

	if !token.Valid {
		log.ServLogger.Error("verify - Invalid token")
		return c.String(http.StatusUnauthorized, "Invalid token")
	}

	cacheTokenMap.Store(id, tokenString)

	return c.NoContent(http.StatusOK)

}

func revoke(c echo.Context) error {
	id := c.Param("id")
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		log.ServLogger.Error("revoke - authHeader is empty")
		return c.String(echo.ErrUnauthorized.Code, "Authorization header is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.String(echo.ErrBadRequest.Code, "Authorization is not Bearer")
	}

	if load, ok := tokenMap.LoadAndDelete(id); ok {
		cacheTokenMap.Delete(id)
		log.ServLogger.Info("revoke - removed token: " + load.(string))
	}

	return c.NoContent(200)
}
