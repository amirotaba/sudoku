package httpd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"runtime"
	"sudoku/internal/domain"
)

type UserHandler struct {
	AUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo,us domain.UserUsecase) {
	handler := &UserHandler{
		AUsecase: us,
	}
	e.POST("signup", handler.SignUp)
	e.POST("signin", handler.SignIn)
	e.GET("account/:username", handler.Account)
	e.POST("boards/save", handler.Save)
	e.GET("board/load/:id", handler.Load)
	e.GET("board/clear/:id", handler.Clear)
	e.POST("board/submit", handler.Submit)
	e.POST("board/new", handler.CreateBoard)
}

func (m *UserHandler) SignUp(e echo.Context) error {
	user := new(domain.User)
	m.AUsecase.SignUp(user)
	msg := &domain.Message{
		Text:   "you logged in as, ",
		UserInfo: user,
		Device: runtime.GOOS,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) SignIn(e echo.Context) error {
	loginform := new(domain.LoginForm)
	user, err := m.AUsecase.SignIn(loginform.PassWord, loginform.Email)
	if err == nil {
		if user.PassWord == loginform.PassWord {
			//jwt
			if errs != nil {
				fmt.Println(errs)
			}
			msg := &domain.Message{
				Text: "you logged in successfuly",
				UserInfo: &user,
				Token: //jwttoken,
			}
			return e.JSON(200, msg)
		} else {
			msg := &domain.Message{
				Text: "incorrect password",
			}
			return e.JSON(200, msg)
		}
	}
	msg := &domain.Message{
		Text: "User not found",
		Err: err.Error(),
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Account(e echo.Context) error {
	username := e.Param("username")
	user, err := m.AUsecase.Account(username)
	if err != nil {
		msg := &domain.Message{
			Text: "User not found",
			Err: err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text: "User info: ",
		UserInfo: &user,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Save(e echo.Context) error {
	board := new(domain.Board)
	if err := m.AUsecase.Save(board); err != nil {
		msg := &domain.Message{
			Text: "Saving failed",
			Err: err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text: "board saved successfully",
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Load(e echo.Context) error {
	board := new(domain.Board)
	b, err := m.AUsecase.Load(board.BoardID)
	if err != nil {
		msg := &domain.Message{
			Text: "board not found",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text: "board loaded successfully",
		Board: &b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Clear(e echo.Context) error {
	board := new(domain.Board)
	b, err := m.AUsecase.Clear(board)
	if err != nil {
		msg := &domain.Message{
			Text: "clearing failed",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text: "board cleared successfully",
		Board: b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Submit(e echo.Context) error {
	board := new(domain.Board)
	b, err := m.AUsecase.Submit(board)
	if err != nil {
		msg := &domain.Message{
			Text: "submition failed",
			Err: err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text: "wrong numbers colored as red",
		Board: b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) CreateBoard(e echo.Context) error {
	board := new(domain.Board)
	b, err := m.AUsecase.CreateBoard(board)
	if err != nil {
		msg := &domain.Message{
			Text: "Creation failed",
			Err: err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text: "Board created successfully",
		Board: &b,
	}
	return e.JSON(200, msg)
}