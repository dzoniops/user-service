package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/dzoniops/user-service/models"
)

// TODO: change secret key and ExpiresAt time
// TODO: maybe add handler like in old project
var jwtKey = []byte("my_secret_key")

type JwtClaims struct {
	*jwt.RegisteredClaims
	Id       int64
	Username string
	Role     string
}

func HashPassword(pass string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user models.User) (signedToken string, err error) {
	claims := &JwtClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Local().Add(time.Hour * time.Duration(100))}},
		Id:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = t.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(signedToken string) (claims *JwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*JwtClaims)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	return claims, nil
}
