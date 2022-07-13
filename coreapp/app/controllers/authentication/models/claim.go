package models

import jwt "github.com/golang-jwt/jwt"

type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}
