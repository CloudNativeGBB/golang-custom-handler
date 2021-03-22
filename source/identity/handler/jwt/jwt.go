package jwt

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"identity/settings"
	"identity/structs"
	"math/big"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var validIssuer = "https://robece.b2clogin.com/d80e5eff-2a87-4324-9b19-cdcc76772d3a/v2.0/"

func init() {
	fmt.Println("package: jwt - initialized")
}

// ValidateToken request
func ValidateToken(request *http.Request, trustFrameworkPolicy string) error {

	strToken := getTokenFromHeader(request)

	if len(strToken) == 0 {
		returnedErr := fmt.Errorf("Token not exists")
		return returnedErr
	}

	keys := settings.GlobalPublicKeys

	key := keys.Keys[0]
	performTokenVerificationErr := performTokenVerification(strToken, key)

	if performTokenVerificationErr != nil {
		return performTokenVerificationErr
	}

	return nil
}

func performTokenVerification(strToken string, key structs.Key) error {

	strPublicKey, generatePublicKeyErr := generatePublicKey(key)

	if generatePublicKeyErr != nil {
		return generatePublicKeyErr
	}

	token, err := jwt.ParseWithClaims(strToken, &structs.UserClaims{}, func(token *jwt.Token) (interface{}, error) {

		parsedKey, parsedKeyErr := jwt.ParseRSAPublicKeyFromPEM([]byte(strPublicKey))

		if parsedKeyErr != nil {
			returnedErr := fmt.Errorf("Parse RSA public key from PEM error: %v", parsedKeyErr)
			return nil, returnedErr
		}

		return parsedKey, nil
	})

	if err != nil {
		return err
	}

	if token.Valid {

		if claims, ok := token.Claims.(*structs.UserClaims); ok {
			if claims.Issuer != validIssuer {
				returnedErr := fmt.Errorf("Invalid issuer")
				return returnedErr
			}
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {

		if ve.Errors&jwt.ValidationErrorMalformed != 0 {

			returnedErr := fmt.Errorf("Malformed token")
			return returnedErr
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {

			returnedErr := fmt.Errorf("Expired token")
			return returnedErr
		} else {

			returnedErr := fmt.Errorf("Couldn't handle this token: %v", err)
			return returnedErr
		}
	} else {

		returnedErr := fmt.Errorf("Couldn't handle this token: %v", err)
		return returnedErr
	}

	return nil
}

func getTokenFromHeader(req *http.Request) string {

	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		return ""
	}

	authHeaderParts := strings.Split(authHeader, " ")

	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return ""
	}

	return authHeaderParts[1]
}

func generatePublicKey(key structs.Key) (string, error) {

	if key.Kty != "RSA" {
		returnedErr := fmt.Errorf("Invalid key type: %v", key.Kty)
		return "", returnedErr
	}

	nb, decodeErr := base64.RawURLEncoding.DecodeString(key.N)
	if decodeErr != nil {
		returnedErr := fmt.Errorf("Decode error: %v", decodeErr)
		return "", returnedErr
	}

	e := 0
	// the default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if key.E == "AQAB" || key.E == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		returnedErr := fmt.Errorf("Need to decode e: %v", key.E)
		return "", returnedErr
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	der, marshalErr := x509.MarshalPKIXPublicKey(pk)
	if marshalErr != nil {
		returnedErr := fmt.Errorf("Marshal public key error: %v", marshalErr)
		return "", returnedErr
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	var out bytes.Buffer
	encodeErr := pem.Encode(&out, block)

	if encodeErr != nil {
		returnedErr := fmt.Errorf("Encode error: %v", encodeErr)
		return "", returnedErr
	}

	return out.String(), nil
}
