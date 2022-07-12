package mysqlHandler

import (
	"gorm.io/gorm"
	"sudoku/internal/domain"
)

type mysqlUserRepository struct {
	Conn *gorm.DB
}

func NewMysqlUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &mysqlUserRepository{Conn: Conn}
}

func (m *mysqlUserRepository) SignUp(user *domain.User) error {
	if err := m.Conn.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (m *mysqlUserRepository) SignIn(password, email string) (domain.User, error) {
	var user domain.User
	if err := m.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (m *mysqlUserRepository) Account(username string) (domain.User, error) {
	var user domain.User
	if err := m.Conn.Where("user_name = ?", username).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (m *mysqlUserRepository) Save(board *domain.Board) error {
	if err := m.Conn.Where("id = ?", board.ID).Update("boarddata", board.BoardData).Error; err != nil {
		return err
	}
	return nil
}

func (m *mysqlUserRepository) Load(id int) (domain.Board, error) {
	var board domain.Board
	if err := m.Conn.Where("id = ?", id).First(&board).Error; err != nil {
		return domain.Board{}, err
	}
	return board, nil
}

func (m *mysqlUserRepository) Remove(id int) error {
	if err := m.Conn.Delete(&domain.Board{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (m *mysqlUserRepository) CreateBoard(board *domain.Board) (*domain.Board, error) {
	if err := m.Conn.Create(&board).Error; err != nil {
		return &domain.Board{}, err
	}
	return board, nil
}
