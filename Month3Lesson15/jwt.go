package jwt

import (
	"github.com/golang-jwt/jwt"
	"test/config"
	"time"
)

func GenerateJWT(m map[interface{}]interface{}) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	aClaims := accessToken.Claims.(jwt.MapClaims)
	rClaims := refreshToken.Claims.(jwt.MapClaims)

	for key, value := range m {
		aClaims[key.(string)] = value
		rClaims[key.(string)] = value
	}

	aClaims["exp"] = time.Now().Add(config.AccessExpireTime).Unix()
	aClaims["iat"] = time.Now().Unix()

	rClaims["exp"] = time.Now().Add(config.RefreshExpireTime).Unix()
	rClaims["iat"] = time.Now().Unix()

	accessToken.Claims = aClaims
	refreshToken.Claims = rClaims

	accessTokenStr, err := accessToken.SignedString(config.SignKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenStr, err := refreshToken.SignedString(config.SignKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func ExtractClaims(tokenString string) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.SignKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	userID, ok := token.Claims.(jwt.MapClaims)["user_id"]
	if ok {
		m["user_id"] = userID
	}

	userRole, ok := token.Claims.(jwt.MapClaims)["user_role"]
	if ok {
		m["user_role"] = userRole
	}

	return m, nil
}