package httpdeliver

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"sudoku/internal/domain"
)

type UserHandler struct {
	AUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		AUsecase: us,
	}
	e.POST("signup", handler.SignUp)
	e.POST("signin", handler.SignIn)
	e.GET("account/:username", handler.Account)
	e.POST("board/new/:username", handler.CreateBoard)
	e.POST("board/save/:username", handler.Save)
	e.GET("board/load/:username/:id", handler.Load)
	e.GET("board/clear/:username/:id", handler.Clear)
	e.POST("board/submit/:username", handler.Submit)

	e.Logger.Fatal(e.Start(":4000"))
}

func (m *UserHandler) SignUp(e echo.Context) error {
	var user domain.User
	if err := e.Bind(user); err != nil {
		return err
	}
	if err := m.AUsecase.SignUp(&user); err != nil {
		return err
	}
	u := domain.UserResponse{UserName: user.UserName, Email: user.Email}
	msg := &domain.AuthMessage{
		Text:     "you logged in as, ",
		UserInfo: &u,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) SignIn(e echo.Context) error {
	loginForm := new(domain.LoginForm)
	if err := e.Bind(loginForm); err != nil {
		return err
	}
	u, err := m.AUsecase.SignIn(loginForm.PassWord, loginForm.Email)
	if err != nil {
		return err
	}
	if u.UserName == "" {
		msg := domain.AuthMessage{
			Text: "incorrect password",
		}
		return e.JSON(200, msg)
	}
	msg := domain.AuthMessage{
		Text:     "you logged in successfully",
		UserInfo: &u,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Account(e echo.Context) error {
	username := e.Param("username")
	user, err := m.AUsecase.Account(username)
	if err != nil {
		msg := &domain.AuthMessage{
			Text: "User not found",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.AuthMessage{
		Text:     "User info: ",
		UserInfo: &user,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Save(e echo.Context) error {
	board := new(domain.Board)
	username := e.Param("username")
	if err := e.Bind(board); err != nil {
		return err
	}
	b, err := m.AUsecase.Save(board, username)
	if err != nil {
		msg := &domain.GameMessage{
			Text: "board not found",
		}
		return e.JSON(200, msg)
	}
	if b.Username == "" {
		msg := &domain.GameMessage{
			Text: "this board doesn't belong to you",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.GameMessage{
		Text:    "board saved successfully, remember the board id for reload",
		Board: b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Load(e echo.Context) error {
	id, _ := strconv.Atoi(e.Param("id"))
	username := e.Param("username")
	b, err := m.AUsecase.Load(id, username)
	if err != nil {
		msg := &domain.GameMessage{
			Text: "board not found",
		}
		return e.JSON(200, msg)
	}
	if b.Username == "" {
		msg := &domain.GameMessage{
			Text: "this board doesn't belong to you",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.GameMessage{
		Text:  "board loaded successfully",
		Board: &b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Clear(e echo.Context) error {
	id, _ := strconv.Atoi(e.Param("id"))
	username := e.Param("username")
	b, err := m.AUsecase.Clear(id, username)
	if err != nil {
		return err
	}
	if b.Username == "" {
		msg := &domain.GameMessage{
			Text: "this board doesn't belong to you",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.GameMessage{
		Text:  "board cleared successfully",
		Board: b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Submit(e echo.Context) error {
	board := new(domain.Board)
	if err := e.Bind(board); err != nil {
		return err
	}
	b, err := m.AUsecase.Submit(board)
	if err != nil {
		msg := &domain.GameMessage{
			Text: "submit failed",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.GameMessage{
		Text:  "wrong numbers are deleted, remember the id for reload",
		Board: b,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) CreateBoard(e echo.Context) error {
	board := new(domain.Board)
	username := e.Param("username")
	if err := e.Bind(board); err != nil {
		return err
	}
	b, err := m.AUsecase.CreateBoard(board, username)
	if err != nil {
		msg := &domain.GameMessage{
			Text: "Creation failed",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.GameMessage{
		Text:  "Board created successfully",
		Board: b,
	}
	return e.JSON(200, msg)
}
