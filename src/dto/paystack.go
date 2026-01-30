package dto

type InitializeTransactionRequest struct {
	Email       string            `json:"email"`
	Amount      int               `json:"amount"` // Amount in kobo (smallest currency unit)
	Currency    string            `json:"currency,omitempty"`
	Reference   string            `json:"reference,omitempty"`
	CallbackURL string            `json:"callback_url,omitempty"`
	Plan        string            `json:"plan,omitempty"`
	Channels    []string          `json:"channels,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type CreatePlanRequest struct {
	Name         string `json:"name"`
	Amount       int    `json:"amount"`   // Amount in kobo
	Interval     string `json:"interval"` // daily, weekly, monthly, annually
	Currency     string `json:"currency,omitempty"`
	Description  string `json:"description,omitempty"`
	InvoiceLimit int    `json:"invoice_limit,omitempty"`
}

type CreateSubscriptionRequest struct {
	Customer      string `json:"customer"` // Customer email or code
	Plan          string `json:"plan"`     // Plan code
	Authorization string `json:"authorization,omitempty"`
	StartDate     string `json:"start_date,omitempty"`
}

type CreateCustomerRequest struct {
	Email     string            `json:"email"`
	FirstName string            `json:"first_name,omitempty"`
	LastName  string            `json:"last_name,omitempty"`
	Phone     string            `json:"phone,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type PaystackResponse[T any] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type InitializeTransactionData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

type VerifyTransactionData struct {
	ID              int                    `json:"id"`
	Status          string                 `json:"status"`
	Reference       string                 `json:"reference"`
	Amount          int                    `json:"amount"`
	Currency        string                 `json:"currency"`
	Channel         string                 `json:"channel"`
	GatewayResponse string                 `json:"gateway_response"`
	PaidAt          string                 `json:"paid_at"`
	Customer        PaystackCustomer       `json:"customer"`
	Authorization   PaystackAuthorization  `json:"authorization"`
	Plan            *PaystackPlanData      `json:"plan"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type PaystackCustomer struct {
	ID           int    `json:"id"`
	CustomerCode string `json:"customer_code"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
}

type PaystackAuthorization struct {
	AuthorizationCode string `json:"authorization_code"`
	Bin               string `json:"bin"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	Channel           string `json:"channel"`
	CardType          string `json:"card_type"`
	Bank              string `json:"bank"`
	CountryCode       string `json:"country_code"`
	Brand             string `json:"brand"`
	Reusable          bool   `json:"reusable"`
}

type PaystackPlanData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PlanCode    string `json:"plan_code"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Interval    string `json:"interval"`
	Currency    string `json:"currency"`
	IsDeleted   bool   `json:"is_deleted"`
	IsArchived  bool   `json:"is_archived"`
}

type PaystackSubscriptionData struct {
	ID               int              `json:"id"`
	SubscriptionCode string           `json:"subscription_code"`
	EmailToken       string           `json:"email_token"`
	Status           string           `json:"status"`
	Amount           int              `json:"amount"`
	NextPaymentDate  string           `json:"next_payment_date"`
	Plan             PaystackPlanData `json:"plan"`
	Customer         PaystackCustomer `json:"customer"`
}

// Webhook Event DTOs

type PaystackWebhookEvent struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// Client-facing DTOs

type InitiatePaymentDto struct {
	SubscriptionTierID string `json:"subscription_tier_id" binding:"required"`
	CallbackURL        string `json:"callback_url,omitempty"`
}

type InitiatePaymentResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

type VerifyPaymentDto struct {
	Reference string `json:"reference" binding:"required"`
}

// Payment Method DTOs

type PaymentMethodResponse struct {
	ID                string `json:"id"`
	AuthorizationCode string `json:"authorization_code"`
	CardType          string `json:"card_type"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	Bank              string `json:"bank"`
	Brand             string `json:"brand"`
	IsDefault         bool   `json:"is_default"`
	Reusable          bool   `json:"reusable"`
}

type AddPaymentMethodDto struct {
	CallbackURL string `json:"callback_url,omitempty"`
}

type SetDefaultPaymentMethodDto struct {
	AuthorizationCode string `json:"authorization_code" binding:"required"`
}

// Invoice DTOs

type InvoiceResponse struct {
	ID          string  `json:"id"`
	Reference   string  `json:"reference"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	PeriodStart string  `json:"period_start"`
	PeriodEnd   string  `json:"period_end"`
	PaidAt      *string `json:"paid_at,omitempty"`
	InvoicePDF  *string `json:"invoice_pdf,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// Paystack Customer Authorization List Response

type PaystackCustomerData struct {
	ID             int                     `json:"id"`
	CustomerCode   string                  `json:"customer_code"`
	Email          string                  `json:"email"`
	FirstName      string                  `json:"first_name"`
	LastName       string                  `json:"last_name"`
	Authorizations []PaystackAuthorization `json:"authorizations"`
}
