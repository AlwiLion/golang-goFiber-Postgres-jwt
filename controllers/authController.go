package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func authentication(c *fiber.Ctx) (*jwt.Token, error) {

	//authorizationHeader := token
	authorizationHeader := c.Get("Authorization")

	token, err := jwt.ParseWithClaims(authorizationHeader, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	return token, err
	// if err != nil {
	// 	return token, err
	// }else{
	// 	return token, err
	// }
}
