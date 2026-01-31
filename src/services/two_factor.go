package services

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type TwoFactorService struct {
	database *gorm.DB
}

func NewTwoFactorService(database *gorm.DB) *TwoFactorService {
	return &TwoFactorService{database: database}
}

func (s *TwoFactorService) GenerateSecret(user *models.User) (*dto.Enable2FAResponse, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Foglio",
		AccountName: user.Email,
		SecretSize:  32,
		Algorithm:   otp.AlgorithmSHA1,
		Digits:      otp.DigitsSix,
	})
	if err != nil {
		return nil, err
	}

	backupCodes, err := s.generateBackupCodes()
	if err != nil {
		return nil, err
	}

	secret := key.Secret()
	user.TwoFactorSecret = &secret
	user.TwoFactorBackupCodes = hashBackupCodes(backupCodes)

	if err := s.database.Save(user).Error; err != nil {
		return nil, err
	}

	return &dto.Enable2FAResponse{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

func (s *TwoFactorService) VerifyAndEnable(userID string, code string) error {
	var user models.User
	if err := s.database.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		return errors.New("2FA setup not initiated")
	}

	valid := totp.Validate(code, *user.TwoFactorSecret)
	if !valid {
		return errors.New("invalid verification code")
	}

	user.IsTwoFactorEnabled = true
	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *TwoFactorService) VerifyCode(userID string, code string) (*models.User, error) {
	var user models.User
	if err := s.database.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsTwoFactorEnabled || user.TwoFactorSecret == nil {
		return nil, errors.New("2FA is not enabled for this user")
	}

	if totp.Validate(code, *user.TwoFactorSecret) {
		return &user, nil
	}

	if s.verifyAndConsumeBackupCode(&user, code) {
		return &user, nil
	}

	return nil, errors.New("invalid verification code")
}

func (s *TwoFactorService) Disable2FA(userID string, password string) error {
	var user models.User
	if err := s.database.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if err := lib.ComparePassword(password, user.Password); err != nil {
		return errors.New("invalid password")
	}

	user.IsTwoFactorEnabled = false
	user.TwoFactorSecret = nil
	user.TwoFactorBackupCodes = nil

	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *TwoFactorService) RegenerateBackupCodes(userID string, password string) ([]string, error) {
	var user models.User
	if err := s.database.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsTwoFactorEnabled {
		return nil, errors.New("2FA is not enabled")
	}

	if err := lib.ComparePassword(password, user.Password); err != nil {
		return nil, errors.New("invalid password")
	}

	backupCodes, err := s.generateBackupCodes()
	if err != nil {
		return nil, err
	}

	user.TwoFactorBackupCodes = hashBackupCodes(backupCodes)
	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return backupCodes, nil
}

func (s *TwoFactorService) GetStatus(userID string) (*dto.TwoFactorStatusResponse, error) {
	var user models.User
	if err := s.database.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	backupCodesLeft := 0
	if user.TwoFactorBackupCodes != nil {
		for _, code := range user.TwoFactorBackupCodes {
			if code != "" && code != "USED" {
				backupCodesLeft++
			}
		}
	}

	return &dto.TwoFactorStatusResponse{
		Enabled:         user.IsTwoFactorEnabled,
		BackupCodesLeft: backupCodesLeft,
	}, nil
}

func (s *TwoFactorService) generateBackupCodes() ([]string, error) {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code, err := generateRandomCode(8)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

func (s *TwoFactorService) verifyAndConsumeBackupCode(user *models.User, code string) bool {
	if user.TwoFactorBackupCodes == nil {
		return false
	}

	normalizedCode := strings.ToUpper(strings.ReplaceAll(code, "-", ""))

	for i, hashedCode := range user.TwoFactorBackupCodes {
		if hashedCode != "" && hashedCode != "USED" {
			// Use bcrypt comparison for secure backup code verification
			if lib.ComparePassword(normalizedCode, hashedCode) == nil {
				user.TwoFactorBackupCodes[i] = "USED"
				s.database.Save(user)
				return true
			}
		}
	}
	return false
}

func generateRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	encoded := base32.StdEncoding.EncodeToString(bytes)
	code := encoded[:length]
	if length >= 8 {
		return code[:4] + "-" + code[4:], nil
	}
	return code, nil
}

func hashBackupCodes(codes []string) []string {
	hashed := make([]string, len(codes))
	for i, code := range codes {
		normalized := strings.ToUpper(strings.ReplaceAll(code, "-", ""))
		hashed[i] = hashBackupCode(normalized)
	}
	return hashed
}

func hashBackupCode(code string) string {
	hashedPassword, err := lib.HashPassword(code)
	if err != nil {
		return ""
	}
	return hashedPassword
}
