package httpd

import (
	"github.com/labstack/echo/v4"
	"runtime"
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
	user := new(domain.User)
	if err := e.Bind(user); err != nil {
		return err
	}
	m.AUsecase.SignUp(user)
	msg := &domain.Message{
		Text:     "you logged in as, ",
		UserInfo: user,
		Device:   runtime.GOOS,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) SignIn(e echo.Context) error {
	loginform := new(domain.LoginForm)
	if err := e.Bind(loginform); err != nil {
		return err
	}
	user, err := m.AUsecase.SignIn(loginform.PassWord, loginform.Email)
	if err == nil {
		if user.PassWord == loginform.PassWord {
			//jwt
			//if errs != nil {
			//	fmt.Println(errs)
			//}
			msg := &domain.Message{
				Text:     "you logged in successfuly",
				UserInfo: &user,
				//Token: jwttoken,
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
		Err:  err.Error(),
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Account(e echo.Context) error {
	username := e.Param("username")
	user, err := m.AUsecase.Account(username)
	if err != nil {
		msg := &domain.Message{
			Text: "User not found",
			Err:  err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text:     "User info: ",
		UserInfo: &user,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Save(e echo.Context) error {
	board := new(domain.Board)
	b2 := new(domain.Board)
	username := e.Param("username")
	if err := e.Bind(board); err != nil {
		return err
	}
	b, err := m.AUsecase.Load(board.ID)
	if err != nil {
		msg := &domain.Message{
			Text: "board not found",
		}
		return e.JSON(200, msg)
	}
	if b.Username != username {
		msg := &domain.Message{
			Text: "this board doesn't belong to you",
		}
		return e.JSON(200, msg)
	}
	board.ID = b2.ID
	id, err := m.AUsecase.Save(board)
	if err != nil {
		msg := &domain.Message{
			Text: "Saving failed",
			Err:  err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text:    "board saved successfully, remember the board id for loading",
		BoardID: id,
	}
	return e.JSON(200, msg)
}

func (m *UserHandler) Load(e echo.Context) error {
	id, _ := strconv.Atoi(e.Param("id"))
	username := e.Param("username")
	b, err := m.AUsecase.Load(id)
	if err != nil {
		msg := &domain.Message{
			Text: "board not found",
		}
		return e.JSON(200, msg)
	}
	if b.Username != username {
		msg := &domain.Message{
			Text: "this board doesn't belong to you",
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text:  "board loaded successfully",
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
		Text:  "board cleared successfully",
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
			Err:  err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text:  "wrong numbers colored as red",
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
		msg := &domain.Message{
			Text: "Creation failed",
			Err:  err.Error(),
		}
		return e.JSON(200, msg)
	}
	msg := &domain.Message{
		Text:  "Board created successfully",
		Board: b,
	}
	return e.JSON(200, msg)
}
