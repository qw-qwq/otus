package server

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/jetuuuu/hl_homework/log"

	"github.com/jetuuuu/hl_homework/mysql"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Store *mysql.DB
}

func (s Server) Login(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	params := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := eCtx.Bind(&params)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	log.Debug(ctx, "params", map[string]interface{}{"params": params})

	params.Password = obfuscate(params.Password)

	exists := s.Store.IsUserExist(ctx, params.Login, params.Password)
	if !exists {
		return eCtx.JSON(http.StatusBadRequest, "is not exist")
	}

	token, err := generateToken(eCtx, params.Login)
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (s Server) GetUser(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	u, err := s.Store.GetUserByID(ctx, eCtx.Get("user_id").(string))
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	names, err := s.Store.GetFriendsNames(ctx, u.Friends)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, err)
	}

	u.Friends = names

	return eCtx.JSON(http.StatusOK, u)
}

func (s Server) CreateUser(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	u := mysql.User{}

	err := eCtx.Bind(&u)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	u.Password = obfuscate(u.Password)

	err = s.Store.CreateUser(ctx, u)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, err)
	}

	token, err := generateToken(eCtx, u.Login)
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, map[string]interface{}{"token": token})
}

func (s Server) MakeFriends(eCtx echo.Context) error {
	params := struct {
		Friend string `json:"friend"`
	}{}
	err := eCtx.Bind(&params)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	user := eCtx.Get("user_id").(string)
	if user == params.Friend {
		return eCtx.JSON(http.StatusBadRequest, nil)
	}

	err = s.Store.MakeFriends(eCtx.Request().Context(), user, params.Friend)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

func generateToken(eCtx echo.Context, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": login,
		"login":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte("super_secret_key"))
	if err != nil {
		return "", eCtx.JSON(http.StatusInternalServerError, err)
	}

	return tokenString, nil
}

func obfuscate(s string) string {
	const salt = "fghdsjkirewo84329fnap"
	log.Debug(context.Background(), "obfuscate", map[string]interface{}{"s": s})

	ret := sha256.Sum256([]byte(s + salt))

	return fmt.Sprintf("%x", ret)
}
