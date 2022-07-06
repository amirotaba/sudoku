package usecase

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"sudoku/internal/domain"
)

type userUsecase struct {
 	UserRepo domain.UserRepository
}

func NewUserUsecase(a domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		UserRepo: a,
	}
}

func (a *userUsecase) SignUp(user *domain.User) error {
	err := a.UserRepo.SignUp(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *userUsecase) SignIn(password, email string) (domain.User, error) {
	u, err := a.UserRepo.SignIn(password, email)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (a *userUsecase) Account(username string) (domain.User, error) {
	user, err := a.UserRepo.Account(username)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (a *userUsecase) Save(board *domain.Board) error {
	if err := a.UserRepo.Save(board); err != nil {
		return err
	}
	return nil
}

func (a *userUsecase) Load(boardID int) (domain.Board, error) {
	b, err := a.UserRepo.Load(boardID)
	if err != nil {
		return domain.Board{}, err
	}
	return b, nil
}

func (a *userUsecase) Clear(board *domain.Board) (*domain.Board, error) {
	board = Delete(board)
	if err := a.UserRepo.Save(board); err != nil {
		return &domain.Board{}, err
	}
	return board, nil
}

func (a *userUsecase) Submit(board *domain.Board) (*domain.Board, error) {
	board = Check(board)
	if err := a.UserRepo.Save(board); err != nil {
		return &domain.Board{}, err
	}
	return board, nil
}

func (a *userUsecase) CreateBoard(number *domain.Board) (domain.Board, error) {
	board := makeboard(number)
	err := a.UserRepo.CreateBoard(&board)
	if err != nil {
		return domain.Board{}, err
	}
	return board, nil
}

func Delete(board *domain.Board) *domain.Board {
	b := board.BoardData
	re := regexp.MustCompile(".")
	b = re.ReplaceAllString(b, "0")
	return &domain.Board{BoardData: b}
}

func Check(board *domain.Board) *domain.Board {
	ebox, boxes := float64(board.Ebox), float64(board.Boxes)
	n := int(math.Sqrt(boxes) * ebox)
	str := board.BoardData
	span := int(math.Pow(float64(n), 2))
	//vertical
	for j := 0; j < n; j++ {
		m := make(map[int]int)
		list := make([]int, 0)
		for i := 0; i < span; i++ {
			x, _ := strconv.Atoi(string(str[i+j]))
			list = append(list, x)
			i = i + n
		}
		for _, i := range list {
			m[list[i]] += 1
		}
		for i := range m {
			if m[i] > 1 {
				str = strings.ReplaceAll(str, string(i), "0")
			}
		}
	}
	//horizontal
}

func makeboard(board *domain.Board) domain.Board {
	var str string
	ebox, boxes := float64(board.Ebox), float64(board.Boxes)
	n := int(math.Sqrt(boxes) * ebox)
	for i := 0; i < n; i++{
		for j := 0; j < n; j++{
			str = str + " 0 "
		}
		str = str + "\n"
	}
	 return domain.Board{BoardData: str}
}
