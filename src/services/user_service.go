package services

import (
	"RestGoTest/src/config"
	"RestGoTest/src/constant"
	"RestGoTest/src/domain/repository"
	"RestGoTest/src/dto"
	"RestGoTest/src/helper/service_errors"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	service "RestGoTest/src/services/dto"
	"RestGoTest/src/util"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger     logging.Logger
	cfg        *config.Config
	OtpService *OtpUsecase
	//database         *gorm.DB
	tokenUserService *TokenUsecase
	repository       repository.UserRepository
}

func NewUserService(cfg *config.Config, repository repository.UserRepository) *UserService {
	logger := logging.NewLogger(cfg)
	return &UserService{
		cfg:              cfg,
		repository:       repository,
		logger:           logger,
		OtpService:       NewOtpUsecase(cfg),
		tokenUserService: NewTokenUsecase(cfg),
	}
}

// Login by username
func (u *UserService) LoginByUsername(ctx context.Context, username string, password string) (*dto.TokenDetail, error) {
	user, err := u.repository.FetchUserInfo(ctx, username, password)

	if err != nil {
		return nil, err
	}
	tokenDto := tokenDto{UserId: user.Id, FirstName: user.FullName, LastName: user.LastName,
		Email: user.Email, MobileNumber: user.MobileNumber}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := u.tokenUserService.GenerateToken(tokenDto)

	if err != nil {
		return nil, err
	}
	return token, nil

}

// Register by username
func (u *UserService) RegisterByUsername(ctx context.Context, req service.RegisterUserByUsername) error {
	user := service.ToUserModel(req)

	exists, err := u.repository.ExistsEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: constant.EmailExists}
	}
	exists, err = u.repository.ExistsUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: constant.UsernameExists}
	}

	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	user.Password = string(hp)
	_, err = u.repository.CreateUser(ctx, user)
	return err

}

// Register/login by mobile number
func (u *UserService) RegisterAndLoginByMobileNumber(ctx context.Context, mobileNumber string, otp string) (*dto.TokenDetail, error) {
	err := u.OtpService.ValidateOtp(mobileNumber, otp)
	if err != nil {
		return nil, err
	}
	exists, err := u.repository.ExistsMobileNumber(ctx, mobileNumber)
	if err != nil {
		return nil, err
	}

	user := model.User{MobileNumber: mobileNumber, Username: mobileNumber}

	if exists {
		user, err = u.repository.FetchUserInfo(ctx, user.Username, user.Password)
		if err != nil {
			return nil, err
		}

		return u.generateToken(user)
	}

	// Register and login
	bp := []byte(util.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	user.Password = string(hp)

	user, err = u.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return u.generateToken(user)

}

func (u *UserService) generateToken(user model.User) (*dto.TokenDetail, error) {
	tokenDto := tokenDto{UserId: user.Id, FirstName: user.FullName, LastName: user.LastName,
		Email: user.Email, MobileNumber: user.MobileNumber}

	if user.UserRoles != nil && len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := u.tokenUserService.GenerateToken(tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}
