package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jetuuuu/hl_homework/mysql"

	"github.com/jetuuuu/hl_homework/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jetuuuu/hl_homework/server"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := database.Open("mysql", "root:@tcp(0.0.0.0:3306)/otus", "main")
	if err != nil {
		panic(err)
	}

	s := server.Server{Store: mysql.New(db)}

	e := echo.New()

	e.POST("/api/login", s.Login)
	e.POST("/api/new", s.CreateUser)

	e.GET("/api/:id", s.GetUser, checkToken)
	e.POST("/api/make_friends", s.MakeFriends, checkToken)

	e.GET("/", s.RenderLoginPage)
	e.GET("/registration", s.RenderRegistrationPage)
	e.GET("/users/:id", s.RenderUserPage, checkToken)

	fmt.Println(e.Start(":1323"))
}

func checkToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("x-auth-token")
		err := parseToken(c, tokenString)
		if err != nil {
			cookie, err := c.Request().Cookie("token")
			if err != nil {
				return err
			}

			err = parseToken(c, cookie.Value)
			if err != nil {
				return c.JSON(http.StatusForbidden, err)
			}
		}

		return next(c)
	}
}

func parseToken(c echo.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("super_secret_key"), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_id", claims["user_id"])
	} else {
		return err
	}

	return nil
}
