package usecase

import (
	"github.com/majidzarephysics/go-jwt/pkg/jwt"
	"math"
	"math/rand"
	"strconv"
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

func (a *userUsecase) SignIn(password, email string) (domain.UserResponse, error) {
	user, err := a.UserRepo.SignIn(password, email)
	if err != nil {
		return domain.UserResponse{}, err
	}
	if user.PassWord != password {
		return domain.UserResponse{}, nil
	}
	jwtsig, errs := jwt.GenerateJWTSigned(user)
	if errs != nil {
		return domain.UserResponse{}, errs
	}
	u := domain.UserResponse{
		UserName: user.UserName,
		Email:    user.Email,
		Token:    jwtsig,
	}
	return u, nil
}

func (a *userUsecase) Account(username string) (domain.UserResponse, error) {
	user, err := a.UserRepo.Account(username)
	if err != nil {
		return domain.UserResponse{}, err
	}
	u := domain.UserResponse{
		UserName: user.UserName,
		Email: user.Email,
	}
	return u, nil
}

func (a *userUsecase) Save(board *domain.Board, username string) (*domain.Board, error) {
	b, err := a.UserRepo.Load(board.ID)
	if err != nil {
		return &domain.Board{}, err
	}
	if b.Username != username {
		return &domain.Board{}, nil
	}
	if b.Defaultid != 0 {
		if err := a.UserRepo.Save(board); err != nil {
			return &domain.Board{}, err
		}
	}
	board.ID = 0
	b2, err := a.UserRepo.CreateBoard(board)
	if err != nil {
		return &domain.Board{}, err
	}
	return b2, nil
}

func (a *userUsecase) Load(id int, username string) (domain.Board, error) {
	b, err := a.UserRepo.Load(id)
	if err != nil {
		return domain.Board{}, err
	}
	if b.Username != username {
		return domain.Board{}, nil
	}

	return b, nil
}

func (a *userUsecase) Clear(id int, username string) (*domain.Board, error) {
	b, err := a.UserRepo.Load(id)
	if err != nil {
		return &domain.Board{}, err
	}
	if b.Username != username {
		return &domain.Board{}, nil
	}
	b, err = a.UserRepo.Load(b.Defaultid)
	if err != nil {
		return &domain.Board{}, err
	}
	if err := a.UserRepo.Remove(id); err != nil {
		return &domain.Board{}, err
	}
	return &b, nil
}

func (a *userUsecase) Submit(board *domain.Board) (*domain.Board, error) {
	board = Check(board)
	b, err := a.Save(board, board.Username)
	if err != nil {
		return &domain.Board{}, err
	}
	board.ID = b.ID
	return board, nil
}

func (a *userUsecase) CreateBoard(number *domain.Board, username string) (*domain.Board, error) {
	board := makeboard(number)
	board.Username = username
	b, err := a.UserRepo.CreateBoard(board)
	if err != nil {
		return &domain.Board{}, err
	}
	b.Defaultid = b.ID
	return b, nil
}

func Check(board *domain.Board) *domain.Board {
	ebox, boxes := float64(board.Ebox), float64(board.Boxes)
	n := int(math.Sqrt(boxes) * math.Sqrt(ebox))
	str := board.BoardData
	span := n * n + n
	//vertical
	for j := 0; j < n; j++ {
		m := make(map[int]int)
		span2 := j + n * n
		for i := 0; i < span2; i++{
			x, _ := strconv.Atoi(string(str[i+j]))
			m[x] += 1
			i = i + n
		}
		for i := 0; i < span2; i++{
			x, _ := strconv.Atoi(string(str[i+j]))
			if x != 0 {
				if m[x] > 1 {
					str = replace(str, i+j)
				}
			}
			i = i + n
		}
	}
	//horizontal
	for i := 0; i < span; i++ {
		m := make(map[int]int)
		for j := 0; j < n; j++ {
			x, _ := strconv.Atoi(string(str[i+j]))
			m[x] += 1
		}
		for j := 0; j < n; j++{
			x, _ := strconv.Atoi(string(str[i+j]))
			if x != 0 {
				if m[x] > 1 {
					str = replace(str, i+j)
				}
			}
		}
		i = i + n
	}
	//box check
	for j := 0; j < n; j++ {
		m := make(map[int]int)
		span := int(math.Sqrt(ebox))
		scale := (n + 1) * span
		for i := 0; i < n * n + 1; i++{
			for l := 0; l < scale; l++ {
				for k := 0; k < span; k++ {
					x, _ := strconv.Atoi(string(str[j+i+l+k]))
					m[x] += 1
				}
				l = l + n
			}
			for l := 0; l < scale; l++ {
				for k := 0; k < span; k++ {
					x, _ := strconv.Atoi(string(str[j+i+l+k]))
					if x != 0 {
						if m[x] > 1 {
							str = replace(str, j+k+l)
						}
					}
				}
				l = l + n
			}
			i = i + (n + 1) * 3 - 1
		}
		j = j + span - 1
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
			if j == random2[l] && l < n-1 {
				str = str + strconv.Itoa(randInt[e])
				l += 1
				e += 1
			}else {
				str = str + "0"
			}
		}
		str = str + "\n"
	}
	str = Check(&domain.Board{BoardData: str}).BoardData
	board.BoardData = str

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

func replace(str string, n int) string {
	str = str[:n] + "0" + str[n+1:]
	return str
}