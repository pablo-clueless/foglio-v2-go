package models

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

type DomainStatus string

const (
	DomainStatusPending  DomainStatus = "PENDING"
	DomainStatusVerified DomainStatus = "VERIFIED"
	DomainStatusFailed   DomainStatus = "FAILED"
)

type Domain struct {
	Subdomain              string       `json:"subdomain"`
	CustomDomain           string       `json:"custom_domain,omitempty"`
	CustomDomainStatus     DomainStatus `json:"custom_domain_status,omitempty"`
	CustomDomainVerifiedAt *time.Time   `json:"custom_domain_verified_at,omitempty"`
	DnsRecords             []DnsRecord  `json:"dns_records,omitempty"`
}

type DnsRecord struct {
	Type   string       `json:"type"`
	Name   string       `json:"name"`
	Value  string       `json:"value"`
	Status DomainStatus `json:"status"`
}

func (d *Domain) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), d)
}

func (d Domain) Value() (driver.Value, error) {
	return json.Marshal(d)
}
func GenerateVerificationToken() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println("error generating token", err)
	}
	return hex.EncodeToString(bytes)
}
func GenerateDnsRecords(customDomain, verificationToken string) []DnsRecord {
	return []DnsRecord{
		{
			Type:   "TXT",
			Name:   "_foglio-verification." + customDomain,
			Value:  verificationToken,
			Status: DomainStatusPending,
		},
		{
			Type:   "CNAME",
			Name:   customDomain,
			Value:  "cname.foglio.app",
			Status: DomainStatusPending,
		},
	}
}
