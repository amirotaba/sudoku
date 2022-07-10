package domain

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
	Email    string `json:"email"`
}

type Number struct {
	Ebox  int
	Boxes int
}

type LoginForm struct {
	PassWord string `json:"pass_word"`
	Email    string `json:"email"`
}

type Board struct {
	Defaultid int
	ID        int    `json:"id" gorm:"primary_key"`
	Ebox      int    `json:"ebox"`
	Boxes     int    `json:"boxes"`
	Username  string `json:"username"`
	BoardData string `json:"boarddata"`
}

type Message struct {
	Text     string
	UserInfo *User
	Device   string
	Token    string
	Err      string
	Board    *Board
	BoardID  int
}
type UserRepository interface {
	SignUp(newuser *User) error
	SignIn(password, email string) (User, error)
	Account(username string) (User, error)
	Save(board *Board) error
	Load(boardID int) (Board, error)
	CreateBoard(board *Board) error
}

type UserUsecase interface {
	SignUp(newuser *User) error
	SignIn(password, email string) (User, error)
	Account(username string) (User, error)
	Save(board *Board) (int, error)
	Load(boardID int) (Board, error)
	Clear(board *Board) (*Board, error)
	Submit(board *Board) (*Board, error)
	CreateBoard(board *Board, username string) (*Board, error)
}
