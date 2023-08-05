package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	token "github.com/dgrijalva/jwt-go/v4"
)

type Token struct {
	Value     string `json:"token"`
	ExpiresIn string `json:"expiresIn"`
}

type jwtClaims struct {
	Value string `json:"value"`
	token.StandardClaims
}

var (
	ExpiryTime         = 31
	PublicKeyLocation  = ""
	PrivateKeyLocation = ""
)

var (
	publicKey     *rsa.PublicKey
	privateKey    *rsa.PrivateKey
	signingMethod token.SigningMethodRSA
	parser        *token.Parser
	keyFunc       token.Keyfunc
)

func Configure() error {
	publicFile, err := ioutil.ReadFile(PublicKeyLocation)
	if err != nil {
		return fmt.Errorf("could not read file %v: %v", PublicKeyLocation, err)
	}
	publicKey, err = token.ParseRSAPublicKeyFromPEM(publicFile)
	if err != nil {
		return fmt.Errorf("could not parse public key: %v", err)
	}

	privateFile, err := ioutil.ReadFile(PrivateKeyLocation)
	if err != nil {
		return fmt.Errorf("could not read file %v: %v", PrivateKeyLocation, err)
	}
	privateKey, err = token.ParseRSAPrivateKeyFromPEM(privateFile)
	if err != nil {
		return fmt.Errorf("could not parse private key: %v", err)
	}

	signingMethod = *token.SigningMethodRS256

	parser = token.NewParser()
	keyFunc = token.KnownKeyfunc(&signingMethod, publicKey)

	return nil
}

func EncodeID(id string) (Token, error) {
	duration, _ := time.ParseDuration(fmt.Sprintf("%vh", ExpiryTime*24))
	expiresAt := time.Now().Add(duration)

	claims := jwtClaims{
		id,
		token.StandardClaims{
			ExpiresAt: token.At(expiresAt),
			Issuer:    "auth.gigfinder.io",
		},
	}

	tkn := token.NewWithClaims(&signingMethod, claims)
	ss, err := tkn.SignedString(privateKey)
	if err != nil {
		return Token{}, fmt.Errorf("cannot sign to token: %v", err)
	}

	return Token{Value: ss, ExpiresIn: fmt.Sprintf("%vd", ExpiryTime)}, nil
}

func Parse(input string) (string, error) {
	tkn, err := parser.ParseWithClaims(input, &jwtClaims{}, keyFunc)
	if err != nil {
		return "", fmt.Errorf("could not parse input token: %v", err)
	}

	if claims, ok := tkn.Claims.(*jwtClaims); ok && tkn.Valid {
		return claims.Value, nil
	} else {
		return "", fmt.Errorf("could not cast claims to required type")
	}
}
