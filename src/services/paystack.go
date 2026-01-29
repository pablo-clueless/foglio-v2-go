package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"foglio/v2/src/config"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	paystackBaseURL = "https://api.paystack.co"
)

type PaystackService struct {
	database   *gorm.DB
	secretKey  string
	httpClient *http.Client
}

func NewPaystackService(database *gorm.DB) *PaystackService {
	return &PaystackService{
		database:  database,
		secretKey: config.AppConfig.PaystackSecretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *PaystackService) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, paystackBaseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("paystack API error: %s", string(respBody))
	}

	return respBody, nil
}

func (s *PaystackService) InitializeTransaction(userID, tierID string, callbackURL string) (*dto.InitiatePaymentResponse, error) {
	var user models.User
	if err := s.database.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var tier models.Subscription
	if err := s.database.First(&tier, "id = ?", tierID).Error; err != nil {
		return nil, errors.New("subscription tier not found")
	}

	var existingSub models.UserSubscription
	err := s.database.Where("user_id = ? AND status = ?", userID, "active").First(&existingSub).Error
	if err == nil {
		return nil, errors.New("user already has an active subscription")
	}

	planCode, err := s.getOrCreatePlan(&tier)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create plan: %w", err)
	}

	reference := fmt.Sprintf("sub_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	amount := int(tier.Price * 100)

	reqBody := dto.InitializeTransactionRequest{
		Email:       user.Email,
		Amount:      amount,
		Currency:    tier.Currency,
		Reference:   reference,
		CallbackURL: callbackURL,
		Plan:        planCode,
		Metadata: map[string]string{
			"user_id":         userID,
			"subscription_id": tierID,
		},
	}

	respBody, err := s.makeRequest("POST", "/transaction/initialize", reqBody)
	if err != nil {
		return nil, err
	}

	var result dto.PaystackResponse[dto.InitializeTransactionData]
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New(result.Message)
	}

	return &dto.InitiatePaymentResponse{
		AuthorizationURL: result.Data.AuthorizationURL,
		AccessCode:       result.Data.AccessCode,
		Reference:        result.Data.Reference,
	}, nil
}

func (s *PaystackService) VerifyTransaction(reference string) (*dto.VerifyTransactionData, error) {
	respBody, err := s.makeRequest("GET", "/transaction/verify/"+reference, nil)
	if err != nil {
		return nil, err
	}

	var result dto.PaystackResponse[dto.VerifyTransactionData]
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New(result.Message)
	}

	return &result.Data, nil
}

func (s *PaystackService) getOrCreatePlan(tier *models.Subscription) (string, error) {
	var existingPlan models.PaystackPlan
	err := s.database.Where("subscription_id = ? AND is_active = ?", tier.ID, true).First(&existingPlan).Error
	if err == nil {
		return existingPlan.PlanCode, nil
	}

	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Where("subscription_id = ? AND is_active = ?", tier.ID, true).First(&existingPlan).Error
	if err == nil {
		tx.Rollback()
		return existingPlan.PlanCode, nil
	}

	interval := "monthly"
	if tier.Type == models.SubscriptionYearly {
		interval = "annually"
	}

	reqBody := dto.CreatePlanRequest{
		Name:     tier.Name,
		Amount:   int(tier.Price * 100), // Convert to kobo
		Interval: interval,
		Currency: tier.Currency,
		Description: func() string {
			if tier.Description != nil {
				return *tier.Description
			}
			return ""
		}(),
	}

	respBody, err := s.makeRequest("POST", "/plan", reqBody)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	var result dto.PaystackResponse[dto.PaystackPlanData]
	if err := json.Unmarshal(respBody, &result); err != nil {
		tx.Rollback()
		return "", err
	}

	if !result.Status {
		tx.Rollback()
		return "", errors.New(result.Message)
	}

	plan := models.PaystackPlan{
		SubscriptionID: tier.ID,
		PlanCode:       result.Data.PlanCode,
		PaystackPlanID: result.Data.ID,
		Interval:       interval,
		IsActive:       true,
	}

	if err := tx.Create(&plan).Error; err != nil {
		tx.Rollback()
		if isDuplicateKeyError(err) {
			var createdPlan models.PaystackPlan
			if err := s.database.Where("subscription_id = ? AND is_active = ?", tier.ID, true).First(&createdPlan).Error; err == nil {
				return createdPlan.PlanCode, nil
			}
		}
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return result.Data.PlanCode, nil
}

func (s *PaystackService) CreateSubscription(customerCode, planCode string) (*dto.PaystackSubscriptionData, error) {
	reqBody := dto.CreateSubscriptionRequest{
		Customer: customerCode,
		Plan:     planCode,
	}

	respBody, err := s.makeRequest("POST", "/subscription", reqBody)
	if err != nil {
		return nil, err
	}

	var result dto.PaystackResponse[dto.PaystackSubscriptionData]
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New(result.Message)
	}

	return &result.Data, nil
}

func (s *PaystackService) DisableSubscription(subscriptionCode, emailToken string) error {
	reqBody := map[string]string{
		"code":  subscriptionCode,
		"token": emailToken,
	}

	respBody, err := s.makeRequest("POST", "/subscription/disable", reqBody)
	if err != nil {
		return err
	}

	var result dto.PaystackResponse[map[string]interface{}]
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	if !result.Status {
		return errors.New(result.Message)
	}

	return nil
}

func (s *PaystackService) GetSubscription(subscriptionCode string) (*dto.PaystackSubscriptionData, error) {
	respBody, err := s.makeRequest("GET", "/subscription/"+subscriptionCode, nil)
	if err != nil {
		return nil, err
	}

	var result dto.PaystackResponse[dto.PaystackSubscriptionData]
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New(result.Message)
	}

	return &result.Data, nil
}

func (s *PaystackService) ProcessSuccessfulPayment(data *dto.VerifyTransactionData) error {
	metadata := data.Metadata
	userIDStr, ok := metadata["user_id"].(string)
	if !ok {
		return errors.New("user_id not found in metadata")
	}
	subscriptionIDStr, ok := metadata["subscription_id"].(string)
	if !ok {
		return errors.New("subscription_id not found in metadata")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user_id")
	}
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		return errors.New("invalid subscription_id")
	}

	// idempotency check
	var existingInvoice models.SubscriptionInvoice
	if err := s.database.Where("paystack_reference = ?", data.Reference).First(&existingInvoice).Error; err == nil {
		return nil
	}

	var tier models.Subscription
	if err := s.database.First(&tier, "id = ?", subscriptionID).Error; err != nil {
		return err
	}

	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// double-check idempotency within transaction
	if err := tx.Where("paystack_reference = ?", data.Reference).First(&existingInvoice).Error; err == nil {
		tx.Rollback()
		return nil
	}

	// row-level locking to prevent race condition
	var existingSub models.UserSubscription
	err = tx.Set("gorm:query_option", "FOR UPDATE SKIP LOCKED").
		Where("user_id = ? AND status = ?", userID, "active").
		First(&existingSub).Error
	if err == nil {
		tx.Rollback()
		return errors.New("user already has an active subscription")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	now := time.Now()
	periodEnd := now.AddDate(0, 0, tier.BillingCycleDays)

	customerCode := data.Customer.CustomerCode
	amountPaid := float64(data.Amount) / 100 // Convert from kobo

	userSub := models.UserSubscription{
		UserID:             userID,
		SubscriptionID:     subscriptionID,
		PaystackCustomerID: &customerCode,
		LastPaymentAmount:  &amountPaid,
		LastPaymentDate:    &now,
		Status:             "active",
		IsActive:           true,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   periodEnd,
	}

	if err := tx.Create(&userSub).Error; err != nil {
		tx.Rollback()
		if isDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	paidAt := now
	invoice := models.SubscriptionInvoice{
		UserSubscriptionID: userSub.ID,
		PaystackReference:  data.Reference,
		AmountPaid:         amountPaid,
		Currency:           data.Currency,
		Status:             "paid",
		PeriodStart:        now,
		PeriodEnd:          periodEnd,
		PaidAt:             &paidAt,
	}

	if err := tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
		if isDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	return tx.Commit().Error
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "SQLSTATE 23505")
}

func (s *PaystackService) VerifyWebhookSignature(body []byte, signature string) bool {
	mac := hmac.New(sha512.New, []byte(config.AppConfig.PaystackWebhookSecret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

func (s *PaystackService) HandleWebhook(event *dto.PaystackWebhookEvent) error {
	switch event.Event {
	case "charge.success":
		return s.handleChargeSuccess(event.Data)
	case "subscription.create":
		return s.handleSubscriptionCreate(event.Data)
	case "subscription.disable":
		return s.handleSubscriptionDisable(event.Data)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(event.Data)
	default:
		return nil
	}
}

func (s *PaystackService) handleChargeSuccess(data map[string]interface{}) error {
	reference, ok := data["reference"].(string)
	if !ok {
		return errors.New("reference not found in webhook data")
	}

	txData, err := s.VerifyTransaction(reference)
	if err != nil {
		return err
	}

	if txData.Status != "success" {
		return nil
	}

	return s.ProcessSuccessfulPayment(txData)
}

func (s *PaystackService) handleSubscriptionCreate(data map[string]interface{}) error {
	subscriptionCode, ok := data["subscription_code"].(string)
	if !ok {
		return nil
	}

	customerData, ok := data["customer"].(map[string]interface{})
	if !ok {
		return nil
	}

	customerCode, ok := customerData["customer_code"].(string)
	if !ok {
		return nil
	}

	return s.database.Model(&models.UserSubscription{}).
		Where("paystack_customer_id = ?", customerCode).
		Update("paystack_subscription_id", subscriptionCode).Error
}

func (s *PaystackService) handleSubscriptionDisable(data map[string]interface{}) error {
	subscriptionCode, ok := data["subscription_code"].(string)
	if !ok {
		return nil
	}

	now := time.Now()
	return s.database.Model(&models.UserSubscription{}).
		Where("paystack_subscription_id = ?", subscriptionCode).
		Updates(map[string]interface{}{
			"status":       "cancelled",
			"is_active":    false,
			"cancelled_at": now,
		}).Error
}

func (s *PaystackService) handlePaymentFailed(data map[string]interface{}) error {
	subscriptionData, ok := data["subscription"].(map[string]interface{})
	if !ok {
		return nil
	}

	subscriptionCode, ok := subscriptionData["subscription_code"].(string)
	if !ok {
		return nil
	}

	return s.database.Model(&models.UserSubscription{}).
		Where("paystack_subscription_id = ?", subscriptionCode).
		Update("status", "past_due").Error
}

func (s *PaystackService) CancelUserSubscription(userID string) error {
	var userSub models.UserSubscription
	err := s.database.Where("user_id = ? AND status = ?", userID, "active").First(&userSub).Error
	if err != nil {
		return errors.New("no active subscription found")
	}

	if userSub.PaystackSubscriptionID != nil && *userSub.PaystackSubscriptionID != "" {
		subData, err := s.GetSubscription(*userSub.PaystackSubscriptionID)
		if err == nil && subData.EmailToken != "" {
			_ = s.DisableSubscription(*userSub.PaystackSubscriptionID, subData.EmailToken)
		}
	}

	now := time.Now()
	userSub.Status = "cancelled"
	userSub.IsActive = false
	userSub.CancelledAt = &now
	userSub.CancelAtPeriodEnd = true

	return s.database.Save(&userSub).Error
}
