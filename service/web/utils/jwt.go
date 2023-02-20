package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"toktik/service/web/config"
)

var jwtConf config.JWT

func GenerateToken(username string) (string, error) {
	now := time.Now()
	claims := jwt.StandardClaims{
		Audience:  username,
		ExpiresAt: now.Add(jwtConf.TokenExpireDuration * time.Minute).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    jwtConf.Issuer,
		NotBefore: now.Unix(),
		Subject:   jwtConf.IdentityKey,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtConf.Secrete))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(tokenStr string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConf.Secrete), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}

func InitConfig(conf config.JWT) {
	jwtConf = conf
}
