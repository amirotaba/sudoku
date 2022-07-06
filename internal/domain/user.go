package domain

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName 	string `json:"username"`
	PassWord 	string `json:"password"`
	Email 		string `json:"email"`
	Boards 		string `json:"boards"`
}

type Number struct {
	Ebox int
	Boxes int
}

type LoginForm struct {
	PassWord string `json:"password"`
	Email    string `json:"email"`
}

type Board struct {
	Ebox 		int
	Boxes		int
	UserID 		int `json:"userid"`
	BoardID 	int `json:"boardId"`
	BoardData 	string `json:"boarddata"`
}

type Message struct {
	Text		string
	UserInfo	*User
	Device		string
	Token		string
	Err			string
	Board		*Board

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
	Save(board *Board) error
	Load(boardID int) (Board, error)
	Clear(board *Board) (*Board, error)
	Submit(board *Board) (*Board, error)
	CreateBoard(board *Board) (Board, error)
}