package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/olzh2102/golang-hotel-reservation/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {

		fmt.Println("-- JWT authentication")

		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		fmt.Println("claims: ", claims)

		// * check token expiration
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		// * set the current authenticated user to the context.
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("never ever print it in real: ", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse token", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
