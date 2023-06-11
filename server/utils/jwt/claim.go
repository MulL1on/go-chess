package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	BaseClaims
}

type BaseClaims struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
}
