package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
)

var currentUserKey struct{}

type CurrentUser struct {
	Username string
}

func GenerateToken(secret, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func JWTAuth(secret string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				tokenString := tr.RequestHeader().Get("Authorization")
				auths := strings.SplitN(tokenString, " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], "Token") {
					return nil, errors.New("jwt token missing")
				}
				token, err := jwt.Parse(auths[1], func(token *jwt.Token) (interface{}, error) {
					fmt.Println("token,,token", token)
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(secret), nil
				})

				if err != nil {
					return nil, err
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if u, ok := claims["username"]; ok {
						ctx = WithContext(ctx, &CurrentUser{Username: u.(string)})
					}
					// put CurrentUser into ctx

				} else {
					return nil, errors.New("Token Invalid")
				}
			}
			return handler(ctx, req)
		}
	}
}

func FromContext(ctx context.Context) *CurrentUser {
	if user := ctx.Value(currentUserKey); user != nil {
		return user.(*CurrentUser)
	}
	return nil
}

func WithContext(ctx context.Context, user *CurrentUser) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}
