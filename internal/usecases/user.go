package usecases

import (
	"context"
	"errors"

	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/hunderaweke/metsasft/internal/repository"
	"github.com/hunderaweke/metsasft/pkg"
	"github.com/sv-tools/mongoifc"
)

type userUsecase struct {
	repo      domain.UserRepository
	tokenRepo domain.TokenRepository
}

func NewUserUsecase(db mongoifc.Database, ctx context.Context) (domain.UserUsecase, bool, error) {
	repo, created, err := repository.NewUserRepository(db, ctx)
	if err != nil {
		return nil, false, err
	}
	tokenRepo := repository.NewTokenRepository(db, ctx)
	return &userUsecase{repo: repo, tokenRepo: tokenRepo}, created, nil
}

func (u *userUsecase) CreateUser(user domain.User) (domain.User, error) {
	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword
	return u.repo.CreateUser(user)
}
func (u *userUsecase) GetUsers() ([]domain.User, error) {
	return u.repo.GetUsers(domain.UserFilter{})
}
func (u *userUsecase) UpdateUser(user domain.User) (domain.User, error) {
	return u.repo.UpdateUser(user)
}
func (u *userUsecase) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}
func (u *userUsecase) GetUserByID(id string) (domain.User, error) {
	return u.repo.GetUserByID(id)
}
func (u *userUsecase) ActivateUser(id string) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	user.IsActive = true
	_, err = u.repo.UpdateUser(user)
	return err
}
func (u *userUsecase) DeactivateUser(id string) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	user.IsActive = false
	_, err = u.repo.UpdateUser(user)
	return err
}

func (u *userUsecase) PromoteUser(id string) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	user.IsAdmin = true
	_, err = u.repo.UpdateUser(user)
	return err
}
func (u *userUsecase) DemoteUser(id string) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	user.IsAdmin = false
	_, err = u.repo.UpdateUser(user)
	return err
}

func (u *userUsecase) GetUserByEmail(email string) (domain.User, error) {
	users, err := u.repo.GetUsers(domain.UserFilter{Email: email})
	if err != nil {
		return domain.User{}, err
	}
	if len(users) == 0 {
		return domain.User{}, &domain.ErrUserNotFound{}
	}
	return users[0], nil
}
func (u *userUsecase) GetUserByTgUserName(username string) (domain.User, error) {
	users, err := u.repo.GetUsers(domain.UserFilter{TelegramUsername: username})
	if err != nil {
		return domain.User{}, err
	}
	return users[0], nil
}

func (u *userUsecase) Login(email, password string) (domain.User, error) {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	if !pkg.ComparePassword(user.Password, password) {
		return domain.User{}, &domain.ErrInvalidCredentials{}
	}
	return user, nil
}

func (u *userUsecase) ForgetPassword(email string) error {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return err
	}
	token, err := pkg.GenerateResetToken()
	if err = u.tokenRepo.CreateToken(domain.Token{Token: token, Email: email}); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return pkg.SendResetEmail(user.Email, token)
}

func (u *userUsecase) ResetPassword(email, token, newPassword string) error {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return err
	}
	t, err := u.tokenRepo.GetTokenByEmail(email)
	if err != nil {
		return err
	}
	if t.Token != token {
		return errors.New("invalid token")
	}
	if err = u.tokenRepo.DeleteToken(email); err != nil {
		return err
	}
	hashedPassword, err := pkg.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	_, err = u.repo.UpdateUser(user)
	return err
}
