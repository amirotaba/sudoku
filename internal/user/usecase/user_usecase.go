package usecase

import (
	"math"
	"math/rand"
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

func (a *userUsecase) Save(board *domain.Board) (int, error) {
	b, err := a.UserRepo.Load(board.Defaultid)
	if err != nil {
		return 0, err
	}
	if b.Defaultid != 0 {
		if err := a.UserRepo.Save(board); err != nil {
			return b.ID, err
		}
	}
	board.Defaultid = b.ID
	board.ID = 0
	if err := a.UserRepo.CreateBoard(board); err != nil {
		return 0, err
	}
	return board.ID, nil
}

func (a *userUsecase) Load(ID int) (domain.Board, error) {
	b, err := a.UserRepo.Load(ID)
	if err != nil {
		return domain.Board{}, err
	}
	return b, nil
}

func (a *userUsecase) Clear(board *domain.Board) (*domain.Board, error) {
	b, err := a.UserRepo.Load(board.Defaultid)
	if err != nil {
		return &domain.Board{}, err
	}
	return &b, nil
}

func (a *userUsecase) Submit(board *domain.Board) (*domain.Board, error) {
	board = Check(board)
	if err := a.UserRepo.Save(board); err != nil {
		return &domain.Board{}, err
	}
	return board, nil
}

func (a *userUsecase) CreateBoard(number *domain.Board, username string) (*domain.Board, error) {
	board := makeboard(number)
	board.Username = username
	err := a.UserRepo.CreateBoard(board)
	if err != nil {
		return &domain.Board{}, err
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
	n := int(math.Sqrt(boxes) * math.Sqrt(ebox))
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
		for k := range list {
			m[k] += 1
		}
		for k := range m {
			if m[k] > 1 {
				str = strings.ReplaceAll(str, string(str[k]), "0")
			}
		}
	}
	//horizontal
	for i := 0; i < span; i++ {
		m := make(map[int]int)
		list := make([]int, 0)
		for j := 0; j < n; j++ {
			x, _ := strconv.Atoi(string(str[i+j]))
			list = append(list, x)
		}
		i = i + n
		for k := range list {
			m[k] += 1
		}
		for k := range m {
			if m[k] > 1 {
				str = strings.ReplaceAll(str, string(str[k]), "0")
			}
		}
	}
	//box check
	//even boxes
	for i := 0; i < int(boxes); i++ {
		m := make(map[int]int)
		list := make([]int, 0)
		span := int(math.Sqrt(ebox))
		for j := 0; j < span; j++ {
			for l := 0; l < int(math.Pow(float64(n), 2)); l++ {
				for k := 0; k < span; k++ {
					x, _ := strconv.Atoi(string(str[k+l+j+i]))
					list = append(list, x)
				}
				l = l + n
			}
			for k := range list {
				m[k] += 1
			}
			for k := range m {
				if m[k] > 1 {
					str = strings.ReplaceAll(str, string(str[k]), "0")
				}
			}
		}
		i = i + (n * span)
	}
	board.BoardData = str
	return board
}

func makeboard(board *domain.Board) *domain.Board {
	var str string
	ebox, boxes := float64(board.Ebox), float64(board.Boxes)
	n := int(math.Sqrt(boxes) * math.Sqrt(ebox))
	var e int
	var l int
	random2 := make([]int, 0)
	randInt := make([]int, 0)
	for {
		if len(randInt) != n || len(random2) != n {
			randInt = append(randInt, rand.Intn(n+1))
			random2 = append(random2, rand.Intn(n))
			randInt = zero(randInt)
			randInt = duplicate(randInt)
			random2 = duplicate(random2)
		} else {
			break
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j == random2[l] && l < 8 {
				str = str + " " + strconv.Itoa(randInt[e]) + " "
				l += 1
				e += 1
			}
			str = str + " 0 "
		}
		str = str + "\n"
	}
	str = Check(&domain.Board{BoardData: str}).BoardData
	board.BoardData = str
	board.Defaultid = board.ID

	return board
}

func duplicate(l []int) []int {
	m := make(map[int]int)
	for i := 0; i < len(l); i++ {
		m[l[i]] += 1
	}
	for _, k := range m {
		if k > 1 {
			l = l[:len(l)-1]
		}
	}
	return l
}

func zero(l []int) []int {
	for i := 0; i < len(l); i++ {
		if l[i] == 0 {
			l = l[:i]
		}
	}
	return l
}
