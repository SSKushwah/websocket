package middlewares

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
)

const AuthContext = "userContext"

//const AuthexpContext = "userexpContext"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mySigningKey := []byte("AllYourBase")
		jwtKey := request.Header.Get("token") //taking input from the postman header
		fmt.Printf("got token: %s\n", jwtKey)

		decryptedToken, err := jwt.Parse(jwtKey, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})

		if claims, ok := decryptedToken.Claims.(jwt.MapClaims); ok && decryptedToken.Valid {
			userID := claims["iss"].(string)

			ctx := context.WithValue(request.Context(), AuthContext, userID)

			next.ServeHTTP(writer, request.WithContext(ctx))
		} else {
			fmt.Println(err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

func GetUserFromContext(request *http.Request) int {

	userID := request.Context().Value(AuthContext)
	userIDInt, convErr := strconv.Atoi(userID.(string))
	if convErr != nil {
		panic(convErr)
	}
	return userIDInt

}
