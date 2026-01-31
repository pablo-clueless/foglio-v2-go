package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"
	"net"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrSubdomainTaken       = errors.New("subdomain is already taken")
	ErrSubdomainInvalid     = errors.New("subdomain contains invalid characters")
	ErrSubdomainReserved    = errors.New("subdomain is reserved")
	ErrCustomDomainRequired = errors.New("custom domain requires a paid subscription")
	ErrDomainAlreadySet     = errors.New("you already have a subdomain set")
	ErrNoDomainConfigured   = errors.New("no domain configured for this user")
	ErrNoCustomDomain       = errors.New("no custom domain configured")
)

var reservedSubdomains = map[string]bool{
	"www":        true,
	"api":        true,
	"app":        true,
	"admin":      true,
	"dashboard":  true,
	"mail":       true,
	"email":      true,
	"help":       true,
	"support":    true,
	"blog":       true,
	"status":     true,
	"docs":       true,
	"cdn":        true,
	"static":     true,
	"assets":     true,
	"media":      true,
	"images":     true,
	"files":      true,
	"download":   true,
	"downloads":  true,
	"upload":     true,
	"uploads":    true,
	"account":    true,
	"accounts":   true,
	"login":      true,
	"logout":     true,
	"signup":     true,
	"signin":     true,
	"register":   true,
	"auth":       true,
	"oauth":      true,
	"sso":        true,
	"test":       true,
	"dev":        true,
	"staging":    true,
	"prod":       true,
	"production": true,
}

type DomainService struct {
	database *gorm.DB
}

func NewDomainService(database *gorm.DB) *DomainService {
	return &DomainService{
		database: database,
	}
}

func (s *DomainService) CheckSubdomainAvailability(subdomain string) (*dto.SubdomainAvailabilityResponse, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))

	if reservedSubdomains[subdomain] {
		return &dto.SubdomainAvailabilityResponse{
			Subdomain: subdomain,
			Available: false,
		}, nil
	}

	if !isValidSubdomain(subdomain) {
		return &dto.SubdomainAvailabilityResponse{
			Subdomain: subdomain,
			Available: false,
		}, nil
	}

	var count int64
	if err := s.database.Model(&models.User{}).
		Where("domain->>'subdomain' = ?", subdomain).
		Count(&count).Error; err != nil {
		return nil, err
	}

	return &dto.SubdomainAvailabilityResponse{
		Subdomain: subdomain,
		Available: count == 0,
	}, nil
}

func (s *DomainService) ClaimSubdomain(userId string, payload dto.ClaimSubdomainDto) (*dto.DomainResponse, error) {
	subdomain := strings.ToLower(strings.TrimSpace(payload.Subdomain))
	if !isValidSubdomain(subdomain) {
		return nil, ErrSubdomainInvalid
	}

	if reservedSubdomains[subdomain] {
		return nil, ErrSubdomainReserved
	}

	var user models.User
	if err := s.database.Preload("CurrentSubscription.Subscription").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	if user.Domain != nil && user.Domain.Subdomain != "" {
		return nil, ErrDomainAlreadySet
	}

	availability, err := s.CheckSubdomainAvailability(subdomain)
	if err != nil {
		return nil, err
	}
	if !availability.Available {
		return nil, ErrSubdomainTaken
	}

	domain := &models.Domain{
		Subdomain: subdomain,
	}
	user.Domain = domain

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return &dto.DomainResponse{
		Subdomain:          domain.Subdomain,
		CanUseCustomDomain: user.CanUseCustomDomain(),
	}, nil
}
func (s *DomainService) GetDomain(userId string) (*dto.DomainResponse, error) {
	var user models.User
	if err := s.database.Preload("CurrentSubscription.Subscription").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	if user.Domain == nil {
		return &dto.DomainResponse{
			CanUseCustomDomain: user.CanUseCustomDomain(),
		}, nil
	}

	return &dto.DomainResponse{
		Subdomain:              user.Domain.Subdomain,
		CustomDomain:           user.Domain.CustomDomain,
		CustomDomainStatus:     user.Domain.CustomDomainStatus,
		CustomDomainVerifiedAt: user.Domain.CustomDomainVerifiedAt,
		DnsRecords:             user.Domain.DnsRecords,
		CanUseCustomDomain:     user.CanUseCustomDomain(),
	}, nil
}
func (s *DomainService) SetCustomDomain(userId string, payload dto.SetCustomDomainDto) (*dto.DomainResponse, error) {
	customDomain := strings.ToLower(strings.TrimSpace(payload.CustomDomain))

	var user models.User
	if err := s.database.Preload("CurrentSubscription.Subscription").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	if !user.CanUseCustomDomain() {
		return nil, ErrCustomDomainRequired
	}

	if user.Domain == nil {
		user.Domain = &models.Domain{}
	}

	verificationToken := models.GenerateVerificationToken()
	dnsRecords := models.GenerateDnsRecords(customDomain, verificationToken)

	user.Domain.CustomDomain = customDomain
	user.Domain.CustomDomainStatus = models.DomainStatusPending
	user.Domain.CustomDomainVerifiedAt = nil
	user.Domain.DnsRecords = dnsRecords

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return &dto.DomainResponse{
		Subdomain:          user.Domain.Subdomain,
		CustomDomain:       user.Domain.CustomDomain,
		CustomDomainStatus: user.Domain.CustomDomainStatus,
		DnsRecords:         user.Domain.DnsRecords,
		CanUseCustomDomain: true,
	}, nil
}
func (s *DomainService) VerifyCustomDomain(userId string) (*dto.DomainResponse, error) {
	var user models.User
	if err := s.database.Preload("CurrentSubscription.Subscription").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	if user.Domain == nil || user.Domain.CustomDomain == "" {
		return nil, ErrNoCustomDomain
	}

	allVerified := true
	for i, record := range user.Domain.DnsRecords {
		verified := verifyDnsRecord(record)
		if verified {
			user.Domain.DnsRecords[i].Status = models.DomainStatusVerified
		} else {
			user.Domain.DnsRecords[i].Status = models.DomainStatusPending
			allVerified = false
		}
	}

	if allVerified {
		user.Domain.CustomDomainStatus = models.DomainStatusVerified
		now := time.Now()
		user.Domain.CustomDomainVerifiedAt = &now
	} else {
		user.Domain.CustomDomainStatus = models.DomainStatusPending
	}

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return &dto.DomainResponse{
		Subdomain:              user.Domain.Subdomain,
		CustomDomain:           user.Domain.CustomDomain,
		CustomDomainStatus:     user.Domain.CustomDomainStatus,
		CustomDomainVerifiedAt: user.Domain.CustomDomainVerifiedAt,
		DnsRecords:             user.Domain.DnsRecords,
		CanUseCustomDomain:     user.CanUseCustomDomain(),
	}, nil
}
func (s *DomainService) RemoveCustomDomain(userId string) (*dto.DomainResponse, error) {
	var user models.User
	if err := s.database.Preload("CurrentSubscription.Subscription").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	if user.Domain == nil || user.Domain.CustomDomain == "" {
		return nil, ErrNoCustomDomain
	}

	user.Domain.CustomDomain = ""
	user.Domain.CustomDomainStatus = ""
	user.Domain.CustomDomainVerifiedAt = nil
	user.Domain.DnsRecords = nil

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return &dto.DomainResponse{
		Subdomain:          user.Domain.Subdomain,
		CanUseCustomDomain: user.CanUseCustomDomain(),
	}, nil
}
func (s *DomainService) UpdateSubdomain(userId string, payload dto.ClaimSubdomainDto) (*dto.DomainResponse, error) {
	subdomain := strings.ToLower(strings.TrimSpace(payload.Subdomain))
	if !isValidSubdomain(subdomain) {
		return nil, ErrSubdomainInvalid
	}
	if reservedSubdomains[subdomain] {
		return nil, ErrSubdomainReserved
	}

	var user models.User
	if err := s.database.Preload("CurrentSubscription.Subscription").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	var count int64
	if err := s.database.Model(&models.User{}).
		Where("domain->>'subdomain' = ? AND id != ?", subdomain, userId).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrSubdomainTaken
	}

	if user.Domain == nil {
		user.Domain = &models.Domain{}
	}

	user.Domain.Subdomain = subdomain

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return &dto.DomainResponse{
		Subdomain:              user.Domain.Subdomain,
		CustomDomain:           user.Domain.CustomDomain,
		CustomDomainStatus:     user.Domain.CustomDomainStatus,
		CustomDomainVerifiedAt: user.Domain.CustomDomainVerifiedAt,
		DnsRecords:             user.Domain.DnsRecords,
		CanUseCustomDomain:     user.CanUseCustomDomain(),
	}, nil
}
func isValidSubdomain(subdomain string) bool {
	if len(subdomain) < 3 || len(subdomain) > 32 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]{1,2}$`, subdomain)
	return matched
}
func verifyDnsRecord(record models.DnsRecord) bool {
	switch record.Type {
	case "TXT":
		return verifyTxtRecord(record.Name, record.Value)
	case "CNAME":
		return verifyCnameRecord(record.Name, record.Value)
	default:
		return false
	}
}
func verifyTxtRecord(name, expectedValue string) bool {
	records, err := net.LookupTXT(name)
	if err != nil {
		return false
	}
	for _, record := range records {
		if record == expectedValue {
			return true
		}
	}
	return false
}
func verifyCnameRecord(name, expectedTarget string) bool {
	cname, err := net.LookupCNAME(name)
	if err != nil {
		return false
	}
	cname = strings.TrimSuffix(cname, ".")
	expectedTarget = strings.TrimSuffix(expectedTarget, ".")
	return strings.EqualFold(cname, expectedTarget)
}
