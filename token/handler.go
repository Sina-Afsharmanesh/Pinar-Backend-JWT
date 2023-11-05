package token

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"jwt/log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	tokenMap      sync.Map
	cacheTokenMap sync.Map
	_, secret, _  = ed25519.GenerateKey(rand.Reader)
)

func Authorize(c echo.Context) error {
	//todo: fetch user role from database, see if it matches the target table and the permission if it does, pass the request
	return c.NoContent(200)
}

func Generate(c echo.Context) error {
	var user User
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		log.ServLogger.Error("generate - user: " + err.Error())
		return c.NoContent(echo.ErrBadRequest.Code)
	}
	log.ServLogger.Info("generate - user: " + user.Id)

	claims := jwt.MapClaims{
		"sub":  user.Id,
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

	tokenMap.Store(user.Id, tokenString)

	return c.String(200, tokenString)
}

func Verify(c echo.Context) error {
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

func Revoke(c echo.Context) error {
	id := c.Param("id")
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		log.ErrLogger.Error("revoke - authHeader is empty")
		return c.String(echo.ErrUnauthorized.Code, "Authorization header is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.String(echo.ErrBadRequest.Code, "Authorization header is not Bearer")
	}

	if load, ok := tokenMap.LoadAndDelete(id); ok {
		cacheTokenMap.Delete(id)
		log.ErrLogger.Info("revoke - removed token: " + load.(string))
	}

	return c.NoContent(200)
}

func GetClaims(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		log.ServLogger.Error("GetClaims - authHeader is empty")
		return c.String(echo.ErrUnauthorized.Code, "Authorization is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.String(echo.ErrBadRequest.Code, "Authorization is not Bearer")
	}

	tokenString := parts[1]

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.ErrLogger.Error("GetClaims - error parsing claims " + err.Error())
		return c.String(echo.ErrBadRequest.Code, err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.ErrLogger.Error("GetClaims - error mapping claims " + err.Error())
		return c.String(echo.ErrInternalServerError.Code, err.Error())
	}

	id := fmt.Sprint(claims["sub"])
	name := fmt.Sprint(claims["name"])
	role := fmt.Sprint(claims["role"])

	result := User{
		Id:   id,
		Name: name,
		Role: role,
	}

	return c.JSON(200, result)
}
