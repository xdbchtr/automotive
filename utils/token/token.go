package token

import (
	"automotive/utils"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type DataClaims struct {
	Authorized bool   `json:"authorized"`
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	Exp        int64  `json:"exp"`
}

func (d DataClaims) Valid() error {
	// Check if the expiration time has passed
	if time.Unix(d.Exp, 0).Before(time.Now()) {
		return errors.New("token has expired")
	}
	return nil
}

var API_SECRET = utils.GetEnv("API_SECRET", "rahasiasekali")

func GenerateToken(username string, userID uint) (string, error) {
	tokenLifespan, err := strconv.Atoi(utils.GetEnv("TOKEN_HOUR_LIFESPAN", "1"))

	if err != nil {
		return "", err
	}

	claims := DataClaims{
		Authorized: true,
		UserID:     userID,
		Username:   username,
		Exp:        time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(API_SECRET))
}
