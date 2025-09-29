package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"

	"gorm.io/gorm"
)

type AuthService struct {
	database *gorm.DB
}

func NewAuthService(database *gorm.DB) *AuthService {
	return &AuthService{
		database: database,
	}
}

type SigninResponse struct {
	User  models.User
	Token string
}

func (s *AuthService) CreateUser(payload dto.CreateUserDto) (*models.User, error) {
	exists, err := s.FindUserByEmail(payload.Email)
	if err != nil {
		return nil, lib.ErrBadRequest
	}
	if exists != nil {
		return nil, errors.New("this email has been used")
	}
	if !lib.ValidatePassword(payload.Password) {
		return nil, errors.New("invalid password")
	}

	hashed, err := lib.HashPassword(payload.Password)
	if err != nil {
		return nil, errors.New("")
	}

	otp := lib.GenerateOtp()

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashed,
		Otp:      otp,
	}

	if err := s.database.Create(&user).Error; err != nil {
		return nil, err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Welcome",
			Template: "welcome",
			Data: map[string]interface{}{
				"Name": []string{user.Username},
				"Otp":  user.Otp,
			},
		})
	}()

	return &user, nil
}

func (s *AuthService) Signin(payload dto.SigninDto) (*SigninResponse, error) {
	user, err := s.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	otp := lib.GenerateOtp()

	user.Otp = otp
	if err = s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	if !user.Verified {
		go func() {
			lib.SendEmail(lib.EmailDto{
				To:       []string{user.Email},
				Subject:  "Verification",
				Template: "verification",
				Data: map[string]interface{}{
					"Name": []string{user.Username},
					"Otp":  user.Otp,
				},
			})
		}()
		return nil, errors.New("user not verified")
	}

	err = lib.ComparePassword(payload.Password, user.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := lib.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.Otp = ""

	return &SigninResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *AuthService) Verification(otp string) error {
	user, err := s.FindUserByOtp(otp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	user.Verified = true

	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Account Verified",
			Template: "verified",
			Data: map[string]interface{}{
				"Name": []string{user.Username},
			},
		})
	}()

	return nil
}

func (s *AuthService) ChangePassword(id string, payload dto.ChangePasswordDto) error {
	user, err := s.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	err = lib.ComparePassword(payload.CurrentPassword, user.Password)
	if err != nil {
		return err
	}

	if !lib.ValidatePassword(payload.NewPassword) {
		return errors.New("invalid password")
	}

	hashed, err := lib.HashPassword(payload.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	token, err := lib.GenerateToken(user.ID)
	if err != nil {
		return err
	}

	url, err := lib.GenerateUrl(token)
	if err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Forgot Password",
			Template: "forgot-password",
			Data: map[string]interface{}{
				"Name": user.Username,
				"Url":  url,
			},
		})
	}()

	return nil
}

func (s *AuthService) ResetPassword(payload dto.ResetPasswordDto) error {
	id, err := lib.ExtractUserID(payload.Token)
	if err != nil {
		return err
	}

	user, err := s.FindUserById(id.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if !lib.ValidatePassword(payload.NewPassword) {
		return errors.New("invalid password")
	}

	hashed, err := lib.HashPassword(payload.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Reset Password",
			Template: "reset-password",
			Data: map[string]interface{}{
				"Name": user.Username,
			},
		})
	}()

	return nil
}

func (s *AuthService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := s.database.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) FindUserById(id string) (models.User, error) {
	var user models.User

	if err := s.database.Where("id = ?", id).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) FindUserByUsername(username string) (models.User, error) {
	var user models.User

	if err := s.database.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
func (s *AuthService) FindUserByOtp(otp string) (models.User, error) {
	var user models.User

	if err := s.database.Where("otp = ?", otp).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) OAuthSignin() {}
