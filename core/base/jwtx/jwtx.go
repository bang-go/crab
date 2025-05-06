package jwtx

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type UuidClaims struct {
	Uuid string `json:"uuid"`
	jwt.RegisteredClaims
}

type UidClaims struct {
	Uid uint64 `json:"uid"`
	jwt.RegisteredClaims
}

type DefaultClaims struct {
	jwt.RegisteredClaims
}

type EntityFunc func(token *jwt.Token) error

type TokenConfig struct {
	Secret        string
	SigningMethod jwt.SigningMethod
	Claims        jwt.Claims
}

// NewToken 新建token
func NewToken(config *TokenConfig) (string, error) {
	var signingMethod jwt.SigningMethod
	if config.SigningMethod != nil {
		signingMethod = config.SigningMethod
	}
	if signingMethod == nil {
		signingMethod = jwt.SigningMethodHS256
	}
	token := jwt.NewWithClaims(signingMethod, config.Claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(config.Secret))
}

// ParseToken 解析并校验Token
func ParseToken(tokenStr string, secret string, claims jwt.Claims, entityFunc EntityFunc) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return
	}
	if err = entityFunc(token); err != nil {
		return
	}
	if !token.Valid {
		err = fmt.Errorf("key is invalid")
		return
	}
	return
}
