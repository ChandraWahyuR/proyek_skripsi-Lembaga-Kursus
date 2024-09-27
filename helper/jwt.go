package helper

import (
	"fmt"
	"skripsi/constant"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserJWT struct {
	ID    string
	Email string
	Role  string
}

type AdminJWT struct {
	ID       string
	Username string
	Role     string
	Email    string
}

type ForgotPassJWT struct {
	ID    string
	Email string
}

type JWT struct {
	signKey string
}

type JWTInterface interface {
	// User Token JWT
	GenerateUserToken(user UserJWT) string
	GenerateUserJWT(user UserJWT) (string, error)
	ExtractUserToken(token *jwt.Token) map[string]interface{}
	// Admin Token JWT
	GenerateAdminToken(admin AdminJWT) string
	GenerateAdminJWT(user AdminJWT) (string, error)
	ExtractAdminToken(token *jwt.Token) map[string]interface{}

	//
	GenerateForgotPassToken(user ForgotPassJWT) string
	GenerateForgotPassJWT(user ForgotPassJWT) (string, error)

	ValidateToken(token string) (*jwt.Token, error)
}

func NewJWT(signKey string) JWTInterface {
	return &JWT{
		signKey: signKey,
	}
}

func (j *JWT) GenerateUserToken(user UserJWT) string {
	var claims = jwt.MapClaims{}
	claims[constant.JWT_ID] = user.ID
	claims[constant.JWT_EMAIL] = user.Email
	claims[constant.JWT_ROLE] = constant.RoleUser
	claims[constant.JWT_IAT] = time.Now().Unix()
	// Sengaja token masa berlaku 1 bulan
	claims[constant.JWT_EXP] = time.Now().Add(time.Hour * 24 * 31).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}
func (j *JWT) GenerateUserJWT(user UserJWT) (string, error) {
	var accessToken = j.GenerateUserToken(user)
	if accessToken == "" {
		return "", constant.ErrGenerateJWT
	}

	return accessToken, nil
}

func (j *JWT) ExtractUserToken(token *jwt.Token) map[string]interface{} {
	if token.Valid {
		var claims = token.Claims
		expTime, _ := claims.GetExpirationTime()
		if expTime.Time.Compare(time.Now()) > 0 {
			var mapClaim = claims.(jwt.MapClaims)
			var result = map[string]interface{}{}
			result[constant.JWT_ID] = mapClaim[constant.JWT_ID]
			result[constant.JWT_EMAIL] = mapClaim[constant.JWT_EMAIL]
			result[constant.JWT_NAME] = mapClaim[constant.JWT_NAME]
			result[constant.JWT_ROLE] = mapClaim[constant.JWT_ROLE]
			return result
		}
		return nil
	}
	return nil
}

func (j *JWT) GenerateAdminToken(admin AdminJWT) string {
	var claims = jwt.MapClaims{}
	claims[constant.JWT_ID] = admin.ID
	claims[constant.JWT_NAME] = admin.Username
	claims[constant.JWT_EMAIL] = admin.Email
	claims[constant.JWT_ROLE] = constant.RoleAdmin
	claims[constant.JWT_IAT] = time.Now().Unix()
	// Sengaja token masa berlaku 1 bulan
	claims[constant.JWT_EXP] = time.Now().Add(time.Hour * 24 * 31).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}

func (j *JWT) GenerateAdminJWT(user AdminJWT) (string, error) {
	var accessToken = j.GenerateAdminToken(user)
	if accessToken == "" {
		return "", constant.ErrGenerateJWT
	}

	return accessToken, nil
}

func (j *JWT) ExtractAdminToken(token *jwt.Token) map[string]interface{} {
	if token.Valid {
		var claims = token.Claims
		expTime, _ := claims.GetExpirationTime()
		if expTime.Time.Compare(time.Now()) > 0 {
			var mapClaim = claims.(jwt.MapClaims)
			var result = map[string]interface{}{}
			result[constant.JWT_ID] = mapClaim[constant.JWT_ID]
			result[constant.JWT_EMAIL] = mapClaim[constant.JWT_EMAIL]
			result[constant.JWT_ROLE] = mapClaim[constant.JWT_ROLE]
			return result
		}
		return nil
	}
	return nil
}

func (j *JWT) GenerateForgotPassToken(user ForgotPassJWT) string {
	var claims = jwt.MapClaims{}
	claims[constant.JWT_ID] = user.ID
	claims[constant.JWT_EMAIL] = user.Email
	claims[constant.JWT_IAT] = time.Now().Unix()
	// Sengaja token masa berlaku 1 bulan
	claims[constant.JWT_EXP] = time.Now().Add(time.Hour * 24 * 31).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}

func (j *JWT) GenerateForgotPassJWT(user ForgotPassJWT) (string, error) {
	var accessToken = j.GenerateForgotPassToken(user)
	if accessToken == "" {
		return "", constant.ErrGenerateJWT
	}

	return accessToken, nil
}

func (j *JWT) ValidateToken(token string) (*jwt.Token, error) {
	if token == "" {
		return nil, constant.ErrValidateJWT
	}
	if len(token) < 7 {
		return nil, constant.ErrValidateJWT
	}

	var authHeader = token[7:]
	parsedToken, err := jwt.Parse(authHeader, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.signKey), nil
	})
	if err != nil {
		return nil, constant.ErrValidateJWT
	}
	return parsedToken, nil
}
