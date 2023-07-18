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

type jwtClaims struct {
	*jwt.StandardClaims
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
	claims := &jwtClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(100)).Unix(),
		},
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

func ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}

	return claims, nil
}
