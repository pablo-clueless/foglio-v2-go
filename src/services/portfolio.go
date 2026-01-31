package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPortfolioNotFound    = errors.New("portfolio not found")
	ErrPortfolioExists      = errors.New("you already have a portfolio")
	ErrSlugTaken            = errors.New("this slug is already taken")
	ErrSlugInvalid          = errors.New("slug contains invalid characters")
	ErrSectionNotFound      = errors.New("section not found")
	ErrUnauthorized         = errors.New("you are not authorized to perform this action")
)

type PortfolioService struct {
	database *gorm.DB
}

func NewPortfolioService(database *gorm.DB) *PortfolioService {
	return &PortfolioService{
		database: database,
	}
}

// CreatePortfolio creates a new portfolio for the user
func (s *PortfolioService) CreatePortfolio(userId string, payload dto.CreatePortfolioDto) (*models.Portfolio, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if user already has a portfolio
	var existingPortfolio models.Portfolio
	if err := s.database.Where("user_id = ?", userUUID).First(&existingPortfolio).Error; err == nil {
		return nil, ErrPortfolioExists
	}

	// Validate and normalize slug
	slug := strings.ToLower(strings.TrimSpace(payload.Slug))
	if !isValidSlug(slug) {
		return nil, ErrSlugInvalid
	}

	// Check if slug is taken
	var slugCount int64
	s.database.Model(&models.Portfolio{}).Where("slug = ?", slug).Count(&slugCount)
	if slugCount > 0 {
		return nil, ErrSlugTaken
	}

	template := "default"
	if payload.Template != "" {
		template = payload.Template
	}

	// Set default settings if not provided
	settings := payload.Settings
	if settings == nil {
		settings = &models.PortfolioSettings{
			ShowProjects:       true,
			ShowExperiences:    true,
			ShowEducation:      true,
			ShowSkills:         true,
			ShowCertifications: true,
			ShowContact:        true,
			ShowSocialLinks:    true,
			EnableAnalytics:    false,
			EnableComments:     false,
		}
	}

	portfolio := &models.Portfolio{
		UserID:   userUUID,
		Title:    payload.Title,
		Slug:     slug,
		Tagline:  payload.Tagline,
		Bio:      payload.Bio,
		Template: template,
		Theme:    payload.Theme,
		Status:   models.PortfolioStatusDraft,
		IsPublic: true,
		SEO:      payload.SEO,
		Settings: settings,
	}

	if err := s.database.Create(portfolio).Error; err != nil {
		return nil, err
	}

	return portfolio, nil
}

// GetPortfolio retrieves a user's portfolio
func (s *PortfolioService) GetPortfolio(userId string) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	if err := s.database.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Where("user_id = ?", userId).First(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPortfolioNotFound
		}
		return nil, err
	}

	return &portfolio, nil
}

// GetPortfolioBySlug retrieves a portfolio by slug (public access)
func (s *PortfolioService) GetPortfolioBySlug(slug string) (*models.Portfolio, *models.User, error) {
	var portfolio models.Portfolio
	if err := s.database.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_visible = ?", true).Order("sort_order ASC")
	}).Where("slug = ? AND status = ? AND is_public = ?", slug, models.PortfolioStatusPublished, true).First(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrPortfolioNotFound
		}
		return nil, nil, err
	}

	// Increment view count
	s.database.Model(&portfolio).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	// Get user info
	var user models.User
	if err := s.database.First(&user, "id = ?", portfolio.UserID).Error; err != nil {
		return &portfolio, nil, nil
	}

	return &portfolio, &user, nil
}

// UpdatePortfolio updates a user's portfolio
func (s *PortfolioService) UpdatePortfolio(userId string, payload dto.UpdatePortfolioDto) (*models.Portfolio, error) {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return nil, err
	}

	if payload.Title != nil {
		portfolio.Title = *payload.Title
	}
	if payload.Slug != nil {
		slug := strings.ToLower(strings.TrimSpace(*payload.Slug))
		if !isValidSlug(slug) {
			return nil, ErrSlugInvalid
		}
		// Check if slug is taken by another portfolio
		var slugCount int64
		s.database.Model(&models.Portfolio{}).Where("slug = ? AND id != ?", slug, portfolio.ID).Count(&slugCount)
		if slugCount > 0 {
			return nil, ErrSlugTaken
		}
		portfolio.Slug = slug
	}
	if payload.Tagline != nil {
		portfolio.Tagline = payload.Tagline
	}
	if payload.Bio != nil {
		portfolio.Bio = payload.Bio
	}
	if payload.CoverImage != nil {
		portfolio.CoverImage = payload.CoverImage
	}
	if payload.Logo != nil {
		portfolio.Logo = payload.Logo
	}
	if payload.Template != nil {
		portfolio.Template = *payload.Template
	}
	if payload.Theme != nil {
		portfolio.Theme = payload.Theme
	}
	if payload.CustomCSS != nil {
		portfolio.CustomCSS = payload.CustomCSS
	}
	if payload.Status != nil {
		portfolio.Status = *payload.Status
	}
	if payload.IsPublic != nil {
		portfolio.IsPublic = *payload.IsPublic
	}
	if payload.SEO != nil {
		portfolio.SEO = payload.SEO
	}
	if payload.Settings != nil {
		portfolio.Settings = payload.Settings
	}

	if err := s.database.Save(portfolio).Error; err != nil {
		return nil, err
	}

	return s.GetPortfolio(userId)
}

// DeletePortfolio deletes a user's portfolio
func (s *PortfolioService) DeletePortfolio(userId string) error {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return err
	}

	return s.database.Delete(portfolio).Error
}

// PublishPortfolio publishes a portfolio
func (s *PortfolioService) PublishPortfolio(userId string) (*models.Portfolio, error) {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return nil, err
	}

	portfolio.Status = models.PortfolioStatusPublished
	if err := s.database.Save(portfolio).Error; err != nil {
		return nil, err
	}

	return portfolio, nil
}

// UnpublishPortfolio unpublishes a portfolio (sets to draft)
func (s *PortfolioService) UnpublishPortfolio(userId string) (*models.Portfolio, error) {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return nil, err
	}

	portfolio.Status = models.PortfolioStatusDraft
	if err := s.database.Save(portfolio).Error; err != nil {
		return nil, err
	}

	return portfolio, nil
}

// CreateSection creates a new section in the portfolio
func (s *PortfolioService) CreateSection(userId string, payload dto.CreatePortfolioSectionDto) (*models.PortfolioSection, error) {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return nil, err
	}

	isVisible := true
	if payload.IsVisible != nil {
		isVisible = *payload.IsVisible
	}

	section := &models.PortfolioSection{
		PortfolioID: portfolio.ID,
		Title:       payload.Title,
		Type:        payload.Type,
		Content:     payload.Content,
		Settings:    payload.Settings,
		SortOrder:   payload.SortOrder,
		IsVisible:   isVisible,
	}

	if err := s.database.Create(section).Error; err != nil {
		return nil, err
	}

	return section, nil
}

// UpdateSection updates a section
func (s *PortfolioService) UpdateSection(userId string, sectionId string, payload dto.UpdatePortfolioSectionDto) (*models.PortfolioSection, error) {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return nil, err
	}

	var section models.PortfolioSection
	if err := s.database.Where("id = ? AND portfolio_id = ?", sectionId, portfolio.ID).First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSectionNotFound
		}
		return nil, err
	}

	if payload.Title != nil {
		section.Title = *payload.Title
	}
	if payload.Type != nil {
		section.Type = *payload.Type
	}
	if payload.Content != nil {
		section.Content = payload.Content
	}
	if payload.Settings != nil {
		section.Settings = payload.Settings
	}
	if payload.SortOrder != nil {
		section.SortOrder = *payload.SortOrder
	}
	if payload.IsVisible != nil {
		section.IsVisible = *payload.IsVisible
	}

	if err := s.database.Save(&section).Error; err != nil {
		return nil, err
	}

	return &section, nil
}

// DeleteSection deletes a section
func (s *PortfolioService) DeleteSection(userId string, sectionId string) error {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return err
	}

	result := s.database.Where("id = ? AND portfolio_id = ?", sectionId, portfolio.ID).Delete(&models.PortfolioSection{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSectionNotFound
	}

	return nil
}

// ReorderSections reorders sections
func (s *PortfolioService) ReorderSections(userId string, payload dto.ReorderSectionsDto) error {
	portfolio, err := s.GetPortfolio(userId)
	if err != nil {
		return err
	}

	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, sectionId := range payload.SectionIDs {
		if err := tx.Model(&models.PortfolioSection{}).
			Where("id = ? AND portfolio_id = ?", sectionId, portfolio.ID).
			Update("sort_order", i).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// ToResponse converts a portfolio model to a response DTO
func (s *PortfolioService) ToResponse(portfolio *models.Portfolio) *dto.PortfolioResponse {
	sections := make([]dto.PortfolioSectionResponse, len(portfolio.Sections))
	for i, section := range portfolio.Sections {
		sections[i] = dto.PortfolioSectionResponse{
			ID:          section.ID.String(),
			PortfolioID: section.PortfolioID.String(),
			Title:       section.Title,
			Type:        section.Type,
			Content:     section.Content,
			Settings:    section.Settings,
			SortOrder:   section.SortOrder,
			IsVisible:   section.IsVisible,
			CreatedAt:   section.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   section.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &dto.PortfolioResponse{
		ID:         portfolio.ID.String(),
		UserID:     portfolio.UserID.String(),
		Title:      portfolio.Title,
		Slug:       portfolio.Slug,
		Tagline:    portfolio.Tagline,
		Bio:        portfolio.Bio,
		CoverImage: portfolio.CoverImage,
		Logo:       portfolio.Logo,
		Template:   portfolio.Template,
		Theme:      portfolio.Theme,
		CustomCSS:  portfolio.CustomCSS,
		Status:     portfolio.Status,
		IsPublic:   portfolio.IsPublic,
		ViewCount:  portfolio.ViewCount,
		SEO:        portfolio.SEO,
		Settings:   portfolio.Settings,
		Sections:   sections,
		CreatedAt:  portfolio.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  portfolio.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPublicResponse converts a portfolio to a public response
func (s *PortfolioService) ToPublicResponse(portfolio *models.Portfolio, user *models.User) *dto.PublicPortfolioResponse {
	sections := make([]dto.PortfolioSectionResponse, len(portfolio.Sections))
	for i, section := range portfolio.Sections {
		sections[i] = dto.PortfolioSectionResponse{
			ID:          section.ID.String(),
			PortfolioID: section.PortfolioID.String(),
			Title:       section.Title,
			Type:        section.Type,
			Content:     section.Content,
			Settings:    section.Settings,
			SortOrder:   section.SortOrder,
			IsVisible:   section.IsVisible,
			CreatedAt:   section.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   section.UpdatedAt.Format(time.RFC3339),
		}
	}

	var userInfo *dto.PublicUserInfo
	if user != nil {
		userInfo = &dto.PublicUserInfo{
			Name:        user.Name,
			Username:    user.Username,
			Headline:    user.Headline,
			Image:       user.Image,
			Location:    user.Location,
			SocialMedia: user.SocialMedia,
		}
	}

	return &dto.PublicPortfolioResponse{
		Title:      portfolio.Title,
		Slug:       portfolio.Slug,
		Tagline:    portfolio.Tagline,
		Bio:        portfolio.Bio,
		CoverImage: portfolio.CoverImage,
		Logo:       portfolio.Logo,
		Template:   portfolio.Template,
		Theme:      portfolio.Theme,
		CustomCSS:  portfolio.CustomCSS,
		SEO:        portfolio.SEO,
		Settings:   portfolio.Settings,
		Sections:   sections,
		User:       userInfo,
	}
}

func isValidSlug(slug string) bool {
	if len(slug) < 3 || len(slug) > 50 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, slug)
	return matched
}
