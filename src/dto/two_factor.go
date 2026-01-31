package dto

type Verify2FASetupRequest struct {
	Code string `json:"code" binding:"required"`
}

type Verify2FALoginRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type Disable2FARequest struct {
	Password string `json:"password" binding:"required"`
}

type Enable2FAResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

type TwoFactorRequiredResponse struct {
	RequiresTwoFactor bool   `json:"requires_two_factor"`
	UserID            string `json:"user_id"`
	Message           string `json:"message"`
}

type TwoFactorStatusResponse struct {
	Enabled         bool `json:"enabled"`
	BackupCodesLeft int  `json:"backup_codes_left"`
}

type BackupCodesResponse struct {
	BackupCodes []string `json:"backup_codes"`
	Message     string   `json:"message"`
}
