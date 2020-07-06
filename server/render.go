package server

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) RenderUserPage(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	u, err := s.Store.GetUserByID(ctx, eCtx.Param("id"))
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	names, err := s.Store.GetFriendsNames(ctx, u.Friends)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, err)
	}

	u.Friends = names

	t, err := template.New("user").Parse(userTemplate)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}

	var w bytes.Buffer
	err = t.Execute(&w, u)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}

	return eCtx.HTML(http.StatusOK, w.String())
}

func (s Server) RenderLoginPage(eCtx echo.Context) error {
	t, err := template.New("login").Parse(loginTemplate)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}

	var w bytes.Buffer
	err = t.Execute(&w, nil)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}

	return eCtx.HTML(http.StatusOK, w.String())
}

func (s Server) RenderRegistrationPage(eCtx echo.Context) error {
	t, err := template.New("login").Parse(registrationTemplate)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}

	var w bytes.Buffer
	err = t.Execute(&w, nil)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}

	return eCtx.HTML(http.StatusOK, w.String())
}
