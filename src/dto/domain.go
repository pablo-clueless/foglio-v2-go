package dto

import (
	"foglio/v2/src/models"
	"time"
)

type ClaimSubdomainDto struct {
	Subdomain string `json:"subdomain" binding:"required,min=3,max=32,alphanum"`
}

type SetCustomDomainDto struct {
	CustomDomain string `json:"custom_domain" binding:"required,fqdn"`
}

type DomainResponse struct {
	Subdomain              string               `json:"subdomain"`
	CustomDomain           string               `json:"custom_domain,omitempty"`
	CustomDomainStatus     models.DomainStatus  `json:"custom_domain_status,omitempty"`
	CustomDomainVerifiedAt *time.Time           `json:"custom_domain_verified_at,omitempty"`
	DnsRecords             []models.DnsRecord   `json:"dns_records,omitempty"`
	CanUseCustomDomain     bool                 `json:"can_use_custom_domain"`
}

type SubdomainAvailabilityResponse struct {
	Subdomain string `json:"subdomain"`
	Available bool   `json:"available"`
}
