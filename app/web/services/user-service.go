package services

import (
	"context"

	"zen-test/app/auth"
	"zen-test/app/exceptions"
	"zen-test/app/helpers"
	"zen-test/app/web/models"
	"zen-test/app/web/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserServiceimpl struct {
	UserRepo repositories.UserRepository
	DB       *gorm.DB
	Validate *validator.Validate
}

type UserService interface {
	Register(ctx context.Context, request models.UserCreate) models.UserResponse
	Update(ctx context.Context, request models.UserUpdate, userId string) models.UserResponse
	Login(ctx context.Context, requestLogin models.UserLogin) (models.UserLoginResponse, bool)
}

func NewUserService(userRepo repositories.UserRepository, db *gorm.DB, validate *validator.Validate) UserService {
	return &UserServiceimpl{
		UserRepo: userRepo,
		DB:       db,
		Validate: validate,
	}
}

func (s *UserServiceimpl) Register(ctx context.Context, request models.UserCreate) models.UserResponse {
	err := s.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	hashPassword, _ := helpers.MakePassword(request.Password)

	user := models.User{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword,
		Phone:    request.Phone,
		Address:  request.Address,
	}

	data, err := s.UserRepo.RegisterUser(ctx, tx, user)
	helpers.PanicIfError(err)

	return models.ToUserReponse(data)
}

func (s *UserServiceimpl) Update(ctx context.Context, request models.UserUpdate, userId string) models.UserResponse {
	err := s.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	userExist, err := s.UserRepo.GetUserById(ctx, tx, userId)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	hashPassword, _ := helpers.MakePassword(request.Password)

	userExist.Password = hashPassword
	userExist.Phone = request.Password
	userExist.Address = request.Address

	data, err := s.UserRepo.UpdateUser(ctx, tx, userExist)
	helpers.PanicIfError(err)

	return models.ToUserReponse(data)
}

func (s *UserServiceimpl) Login(ctx context.Context, requestLogin models.UserLogin) (models.UserLoginResponse, bool) {
	err := s.Validate.Struct(requestLogin)
	helpers.PanicIfError(err)

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	user, err := s.UserRepo.GetUserByEmail(ctx, tx, requestLogin.Email)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	passwordSync := helpers.ComparePassword(requestLogin.Password, user.Password)
	accessToken, _ := auth.CreateToken(user.ID)
	refreshToken, _ := auth.CreateRefreshToken(user.ID)

	if !passwordSync {
		return models.UserLoginResponse{}, false
	} else {

		userLoginResponse := models.UserLoginResponse{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Token:        accessToken,
			RefreshToken: refreshToken,
		}
		return userLoginResponse, true
	}
}
