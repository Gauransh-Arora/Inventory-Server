package utils

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func InitJWT() error {
	rawPriv := strings.Trim(os.Getenv("JWT_PRIVATE_KEY"), "'\"")
	rawPub := strings.Trim(os.Getenv("JWT_PUBLIC_KEY"), "'\"")

	privPEM := strings.ReplaceAll(rawPriv, "\\n", "\n")
	pubPEM := strings.ReplaceAll(rawPub, "\\n", "\n")

	var err error
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	return nil
}

func GenerateAccessToken(userID uuid.UUID, username, role string) (string, string, int64, error) {
	jti := uuid.New().String()
	exp := time.Now().Add(30 * 24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"userId":   userID.String(),
		"username": username,
		"role":     role,
		"jti":      jti,
		"iss":      os.Getenv("JWT_ISSUER"),
		"aud":      os.Getenv("JWT_AUDIENCE"),
		"exp":      exp,
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	str, err := token.SignedString(signKey)
	return str, jti, exp, err
}

func GetVerifyKey() *rsa.PublicKey {
	return verifyKey
}