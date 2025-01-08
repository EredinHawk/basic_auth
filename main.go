package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Login    string
	Password string
}

func main() {
	router := gin.Default()
	router.GET("/", handler)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}

func handler(c *gin.Context) {
	if err := headerCheck(c.GetHeader("Authorization")); err != nil {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm='user_pages'")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	credential, err := decodeBase64(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"login": credential.Login, "password": credential.Password})
}

// headerCheck проверяет наличие заголовка 'Authorization' в http запросе
func headerCheck(s string) error {
	if s == "" {
		return fmt.Errorf("headerCheck error: the Authorization header is empty")
	}
	return nil
}

// decodeBase64 декодирует заголовок 'Authorization' в котором содержится логин и пароль
func decodeBase64(s string) (Credentials, error) {
	result, err := base64.StdEncoding.DecodeString(strings.Split(s, " ")[1])
	if err != nil {
		return Credentials{}, fmt.Errorf("decodeBase64 error: %v", err.Error())
	}

	cred := Credentials{
		Login:    strings.Split(string(result), ":")[0],
		Password: strings.Split(string(result), ":")[1],
	}
	return cred, nil
}
