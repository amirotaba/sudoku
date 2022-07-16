package domain

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
	Email    string `json:"email"`
}

type UserResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Token    string
}

type LoginForm struct {
	Email    string `json:"email"`
	PassWord string `json:"pass_word"`
}

type Board struct {
	Defaultid int
	ID        int    `json:"id" gorm:"primary_key"`
	Ebox      int    `json:"ebox"`
	Boxes     int    `json:"boxes"`
	Username  string `json:"username"`
	BoardData string `json:"boarddata"`
}

type AuthMessage struct {
	Text     string
	UserInfo *UserResponse
}

type SignUpMessage struct {
	Text     string
	UserName string
	Email    string
}

type GameMessage struct {
	Text  string
	Board *Board
}

type UserRepository interface {
	SignUp(newuser *User) error
	SignIn(password, email string) (User, error)
	Account(username string) (User, error)
	Save(board *Board) error
	Load(boardID int) (Board, error)
	Remove(id int) error
	CreateBoard(board *Board) (*Board, error)
}

type UserUsecase interface {
	SignUp(newuser *User) error
	SignIn(password, email string) (UserResponse, error)
	Account(username string) (UserResponse, error)
	Save(board *Board, username string) (*Board, error)
	Load(boardID int, username string) (Board, error)
	Clear(id int, username string) (*Board, error)
	Submit(board *Board) (*Board, error)
	CreateBoard(board *Board, username string) (*Board, error)
}
