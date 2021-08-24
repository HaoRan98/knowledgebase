package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

type CustomClaims struct {
	Username string `json:"username"`
	DeptID   string `json:"dept_id"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username, dept_id string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(330 * 24 * time.Hour)

	claims := CustomClaims{
		username,
		dept_id,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "dingtalk",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (interface{}, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method:%v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	//log.Println("jwt:",tokenClaims)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok {
			if tokenClaims.Valid {
				return claims, nil
			} else { //if token valid false,if validation error is expired,return claims
				if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
					return claims, err
				}
			}
		}
	}
	return nil, err
}

// according to token, return username
func GetLoginID(token string, c *gin.Context) string {
	if token == "" {
		token = c.GetHeader("X-Access-Token")
		auth := c.GetHeader("Authorization")
		//token := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
	}
	claims, _ := ParseToken(token)
	return claims.(jwt.MapClaims)["username"].(string)
}

func GetDeptID(token string, c *gin.Context) string {
	if token == "" {
		token = c.GetHeader("X-Access-Token")
		auth := c.GetHeader("Authorization")
		//token := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
	}
	claims, _ := ParseToken(token)
	return claims.(jwt.MapClaims)["dept_id"].(string)
}
