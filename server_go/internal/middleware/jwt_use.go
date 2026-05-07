package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const TokenExpireDuration = time.Hour * 2

// const TokenExpireDuration = time.Second * 60

var Secret = []byte("subjectcourse")

type MyClaims struct {
	UserName string
	UserType int //0 是超级管理员  1是学校管理员   2是老师用户
	SchoolId int
	UserId   int
	jwt.StandardClaims
}

// get token
func GetToken(username string, userType int, schoolId int, userId int) (string, error) {
	cla := MyClaims{
		username,
		userType,
		schoolId,
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "lx-jwt",                                   // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	fmt.Println("Token = ", token)
	return token.SignedString(Secret) // 进行签名生成对应的token
}

// parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 中间件,认证token合法性
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHandler := c.Request.Header.Get("authorization")
		if authHandler == "" {
			c.JSON(200, gin.H{"code": 2003, "message": "请求头部auth为空", "data": "{}"})
			c.Abort()
			return
		}

		// 前两部门可以直接解析出来
		jwt := strings.Split(authHandler, ".")
		cnt := 0
		for _, val := range jwt {
			cnt++
			if cnt == 3 {
				break
			}
			msg, _ := base64.StdEncoding.DecodeString(val)
			fmt.Println("val ->", string(msg))
		}

		// 我们使用之前定义好的解析JWT的函数来解析它,并且在内部解析时判断了token是否过期
		mc, err := ParseToken(authHandler)
		if err != nil {
			fmt.Println("err = ", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code":    2005,
				"message": "无效的Token",
				"data":    "{}",
			})
			c.Abort()
			return
		}

		if (strings.Contains(c.FullPath(), "/user/") && mc.UserType != 2) ||
			(strings.Contains(c.FullPath(), "/schooladmin/") && mc.UserType != 1) ||
			(strings.Contains(c.FullPath(), "/administrator/") && mc.UserType != 0) ||
			(strings.Contains(c.FullPath(), "/currency/") && !(mc.UserType == 0 || mc.UserType == 1 || mc.UserType == 2)) ||
			(strings.Contains(c.FullPath(), "/backstage/") && !(mc.UserType == 0 || mc.UserType == 1)) ||
			(strings.Contains(c.FullPath(), "/commonuser/") && !(mc.UserType == 2)) {
			c.JSON(http.StatusOK, gin.H{
				"code":    2004,
				"message": "Token不匹配请求的接口",
				"data":    "{}",
			})
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.UserName)
		c.Set("userType", mc.UserType)
		c.Set("schoolId", mc.SchoolId)
		c.Set("userId", mc.UserId)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
