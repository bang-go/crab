package jwtx

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"testing"
	"time"
)

var secret string = "password"

func TestNewToken(t *testing.T) {
	claims := &UuidClaims{
		Uuid: "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}
	tokenStr, err := NewToken(&TokenConfig{Secret: secret, Claims: claims})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(tokenStr)
}

func TestParseToken(t *testing.T) {
	claims := &UuidClaims{}
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMTIzIiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiYXVkIjpbInNvbWVib2R5X2Vsc2UiXSwiZXhwIjoxNjgwMjQ2ODYxLCJuYmYiOjE2ODAxNjA0NjEsImlhdCI6MTY4MDE2MDQ2MSwianRpIjoiMSJ9.gP4lKcnuX7bpGahqJ5C49PlSetpGAkPRnvtyQCi1Gmg"
	_, err := ParseToken(tokenStr, secret, claims, func(token *jwt.Token) error {
		var ok bool
		if claims, ok = token.Claims.(*UuidClaims); !ok {
			return errors.New("err entity")
		}
		log.Println(claims)
		return nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Fatal("expired")
		}
		log.Fatal(err)
	}
	if claims.Uuid != "" { //有效
		log.Println("token valid")
	}
}
