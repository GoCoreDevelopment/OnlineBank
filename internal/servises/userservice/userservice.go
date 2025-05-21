package userservice

import (
	"api/internal/models/user"
	"api/internal/repository"
	"api/internal/utils/hashpassword"
	"errors"
)

type UserService interface {
	Registred(user user.UserRegistred) (int, error)
	Login(user user.UserLogin) (int, error)
	GetUserIDByEmail(email string) (int, error)
	CheckBalance(id int) (int, error)
	Transfer(senderID, receiverID, amount int) error
}

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Registred(user user.UserRegistred) (int, error) {
	if user.FirstName == "" || user.LastName == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		return 0, errors.New("all fields must be filled in")
	}

	hash, err := hashpassword.CreateHash(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = hash

	return s.repo.CreateUser(user)
}

func (s *userService) Login(user user.UserLogin) (int, error) {
	passwordHash, err := s.repo.GetUserPassword(user.Email)
	if err != nil {
		return 0, err
	}

	err = hashpassword.CheckValidHash(passwordHash, user.Password)
	if err != nil {
		return 0, errors.New("error login or password")
	}

	return s.repo.GetUserID(user.Email)
}	

func (s *userService) CheckBalance(id int) (int, error) {
	return s.repo.GetUserBalance(id)
}

func (s *userService) GetUserIDByEmail(email string) (int, error) {
	return s.repo.GetUserID(email)
}

func (s *userService) Transfer(senderID, receiverID, amount int) error {
	if amount <= 0 {
		return errors.New("the number cannot be zero")
	}

	if senderID == receiverID {
		return errors.New("you can't transfer money to yourself")
	}

	return s.repo.TransferMoney(senderID, receiverID, amount)
}

