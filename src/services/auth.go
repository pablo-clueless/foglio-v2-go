package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"foglio/v2/src/config"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"log"
	"net/http"
	"strings"
	"time"

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

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

func (s *AuthService) getGoogleOAuthConfig() OAuthConfig {
	return OAuthConfig{
		ClientID:     config.AppConfig.GoogleClientId,
		ClientSecret: config.AppConfig.GoogleClientSecret,
		RedirectURL:  config.AppConfig.GoogleRedirectUrl,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
		UserInfoURL:  "https://www.googleapis.com/oauth2/v3/userinfo",
	}
}

func (s *AuthService) getGitHubOAuthConfig() OAuthConfig {
	return OAuthConfig{
		ClientID:     config.AppConfig.GithubClientId,
		ClientSecret: config.AppConfig.GithubClientSecret,
		RedirectURL:  config.AppConfig.GithubRedirectUrl,
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		UserInfoURL:  "https://api.github.com/user",
	}
}

type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type SigninResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func (s *AuthService) CreateUser(payload dto.CreateUserDto) (*models.User, error) {
	fmt.Println(payload)
	emailExists, err := s.FindUserByEmail(payload.Email)
	if err == nil && emailExists != nil {
		return nil, errors.New("this email has been used")
	}
	usernameExists, err := s.FindUserByUsername(payload.Username)
	if err == nil && usernameExists != nil {
		return nil, errors.New("this username has been used")
	}
	if !lib.ValidatePassword(payload.Password) {
		return nil, errors.New("invalid password")
	}

	hashed, err := lib.HashPassword(payload.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	otp := lib.GenerateOtp()

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Otp:      otp,
		Password: hashed,
		Username: payload.Username,
	}

	if err := s.database.Create(&user).Error; err != nil {
		return nil, err
	}

	go func() {
		err := lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Welcome",
			Template: "welcome",
			Data: map[string]interface{}{
				"Name":  user.Name,
				"Email": user.Email,
				"Otp":   user.Otp,
			},
		})
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent to: %v", user.Email)
		}
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

	err = lib.ComparePassword(payload.Password, user.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	if !user.Verified {
		otp := lib.GenerateOtp()
		user.Otp = otp
		if err = s.database.Save(&user).Error; err != nil {
			return nil, err
		}

		go func() {
			err = lib.SendEmail(lib.EmailDto{
				To:       []string{user.Email},
				Subject:  "Verification",
				Template: "verification",
				Data: map[string]interface{}{
					"Name":  user.Name,
					"Email": user.Email,
					"Otp":   otp,
				},
			})
			if err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email sent to: %v", user.Email)
			}
		}()
		return nil, errors.New("user not verified. a verification mail has been sent to your email.")
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

func (s *AuthService) Verification(otp string) (*SigninResponse, error) {
	var user models.User

	if err := s.database.Where("otp = ?", otp).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Verified = true
	user.Otp = ""
	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	go func() {
		err := lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Account Verified",
			Template: "verified",
			Data: map[string]interface{}{
				"Name":  user.Name,
				"Email": user.Email,
			},
		})
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent to: %v", user.Email)
		}
	}()

	token, err := lib.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.Otp = ""

	return &SigninResponse{
		User:  user,
		Token: token,
	}, nil
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

	client := config.AppConfig.ClientUrl + "/reset-password"
	url := lib.GenerateUrl(client, token)

	go func() {
		err := lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Forgot Password",
			Template: "forgot-password",
			Data: map[string]interface{}{
				"Name":  user.Username,
				"Email": user.Email,
				"Url":   url,
			},
		})
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent to: %v", user.Email)
		}
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
		err := lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Reset Password",
			Template: "reset-password",
			Data: map[string]interface{}{
				"Name":  user.Name,
				"Email": user.Email,
			},
		})
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent to: %v", user.Email)
		}
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

func (s *AuthService) FindUserById(id string) (*models.User, error) {
	var user models.User

	if err := s.database.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) FindUserByUsername(username string) (*models.User, error) {
	var user models.User

	if err := s.database.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) GetOAuthURL(provider string) (string, error) {
	var config OAuthConfig
	var scope string

	switch provider {
	case "google":
		config = s.getGoogleOAuthConfig()
		scope = "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email"
	case "github":
		config = s.getGitHubOAuthConfig()
		scope = "user:email"
	default:
		return "", errors.New("unsupported OAuth provider")
	}

	state := lib.GenerateRandomString(32)

	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		config.AuthURL,
		config.ClientID,
		config.RedirectURL,
		scope,
		state,
	)

	return authURL, nil
}

func (s *AuthService) HandleOAuthCallback(provider string, payload dto.OAuthCallbackDto) (*SigninResponse, error) {
	accessToken, err := s.exchangeCodeForToken(provider, payload.Code)
	if err != nil {
		return nil, err
	}

	oauthUser, err := s.getUserInfo(provider, accessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.findOrCreateOAuthUser(oauthUser)
	if err != nil {
		return nil, err
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

func (s *AuthService) exchangeCodeForToken(provider string, code string) (string, error) {
	var config OAuthConfig
	switch provider {
	case "google":
		config = s.getGoogleOAuthConfig()
	case "github":
		config = s.getGitHubOAuthConfig()
	default:
		return "", errors.New("unsupported OAuth provider")
	}

	body := strings.NewReader(fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s",
		config.ClientID,
		config.ClientSecret,
		code,
		config.RedirectURL,
	))

	req, err := http.NewRequest("POST", config.TokenURL, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	var tokenResp OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", errors.New("failed to get access token")
	}

	return tokenResp.AccessToken, nil
}

func (s *AuthService) getUserInfo(provider string, accessToken string) (*dto.OAuthUserDto, error) {
	var config OAuthConfig
	switch provider {
	case "google":
		config = s.getGoogleOAuthConfig()
	case "github":
		config = s.getGitHubOAuthConfig()
	default:
		return nil, errors.New("unsupported OAuth provider")
	}

	req, err := http.NewRequest("GET", config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	var oauthUser dto.OAuthUserDto
	oauthUser.Provider = provider

	switch provider {
	case "google":
		var googleUser struct {
			Sub     string `json:"sub"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
			return nil, err
		}
		oauthUser.ID = googleUser.Sub
		oauthUser.Email = googleUser.Email
		oauthUser.Name = googleUser.Name
		oauthUser.Avatar = googleUser.Picture

	case "github":
		var githubUser struct {
			ID        int    `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Login     string `json:"login"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			return nil, err
		}
		oauthUser.ID = fmt.Sprintf("%d", githubUser.ID)
		oauthUser.Email = githubUser.Email
		oauthUser.Name = githubUser.Name
		oauthUser.Avatar = githubUser.AvatarURL

		if oauthUser.Email == "" {
			email, err := s.getGitHubPrimaryEmail(accessToken)
			if err == nil && email != "" {
				oauthUser.Email = email
			}
		}

		if oauthUser.Name == "" {
			oauthUser.Name = githubUser.Login
		}
	}

	if oauthUser.Email == "" {
		return nil, errors.New("could not get user email from OAuth provider")
	}

	return &oauthUser, nil
}

func (s *AuthService) getGitHubPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", errors.New("no primary email found")
}

func (s *AuthService) findOrCreateOAuthUser(oauthUser *dto.OAuthUserDto) (*models.User, error) {
	var user models.User

	err := s.database.Where("provider = ? AND provider_id = ?", oauthUser.Provider, oauthUser.ID).First(&user).Error
	if err == nil {
		return &user, nil
	}

	err = s.database.Where("email = ?", oauthUser.Email).First(&user).Error
	if err == nil {
		user.Provider = oauthUser.Provider
		user.ProviderID = oauthUser.ID
		if user.Image == nil || *user.Image == "" {
			user.Image = &oauthUser.Avatar
		}
		if err := s.database.Save(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	user = models.User{
		Name:       oauthUser.Name,
		Email:      oauthUser.Email,
		Username:   lib.GenerateUsername(oauthUser.Name),
		Provider:   oauthUser.Provider,
		ProviderID: oauthUser.ID,
		Verified:   true,
	}

	if oauthUser.Avatar != "" {
		user.Image = &oauthUser.Avatar
	}

	if err := s.database.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
