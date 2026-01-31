package dto

import "foglio/v2/src/models"

type CreatePortfolioDto struct {
	Title    string                   `json:"title" binding:"required,min=1,max=100"`
	Slug     string                   `json:"slug" binding:"required,min=3,max=50"`
	Tagline  *string                  `json:"tagline,omitempty"`
	Bio      *string                  `json:"bio,omitempty"`
	Template string                   `json:"template,omitempty"`
	Theme    *models.PortfolioTheme   `json:"theme,omitempty"`
	SEO      *models.PortfolioSEO     `json:"seo,omitempty"`
	Settings *models.PortfolioSettings `json:"settings,omitempty"`
}

type UpdatePortfolioDto struct {
	Title      *string                   `json:"title,omitempty"`
	Slug       *string                   `json:"slug,omitempty"`
	Tagline    *string                   `json:"tagline,omitempty"`
	Bio        *string                   `json:"bio,omitempty"`
	CoverImage *string                   `json:"cover_image,omitempty"`
	Logo       *string                   `json:"logo,omitempty"`
	Template   *string                   `json:"template,omitempty"`
	Theme      *models.PortfolioTheme    `json:"theme,omitempty"`
	CustomCSS  *string                   `json:"custom_css,omitempty"`
	Status     *models.PortfolioStatus   `json:"status,omitempty"`
	IsPublic   *bool                     `json:"is_public,omitempty"`
	SEO        *models.PortfolioSEO      `json:"seo,omitempty"`
	Settings   *models.PortfolioSettings `json:"settings,omitempty"`
}

type CreatePortfolioSectionDto struct {
	Title     string  `json:"title" binding:"required"`
	Type      string  `json:"type" binding:"required,oneof=hero about projects experience skills contact custom"`
	Content   *string `json:"content,omitempty"`
	Settings  *string `json:"settings,omitempty"`
	SortOrder int     `json:"sort_order"`
	IsVisible *bool   `json:"is_visible,omitempty"`
}

type UpdatePortfolioSectionDto struct {
	Title     *string `json:"title,omitempty"`
	Type      *string `json:"type,omitempty"`
	Content   *string `json:"content,omitempty"`
	Settings  *string `json:"settings,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
	IsVisible *bool   `json:"is_visible,omitempty"`
}

type ReorderSectionsDto struct {
	SectionIDs []string `json:"section_ids" binding:"required"`
}

type PortfolioResponse struct {
	ID          string                    `json:"id"`
	UserID      string                    `json:"user_id"`
	Title       string                    `json:"title"`
	Slug        string                    `json:"slug"`
	Tagline     *string                   `json:"tagline,omitempty"`
	Bio         *string                   `json:"bio,omitempty"`
	CoverImage  *string                   `json:"cover_image,omitempty"`
	Logo        *string                   `json:"logo,omitempty"`
	Template    string                    `json:"template"`
	Theme       *models.PortfolioTheme    `json:"theme,omitempty"`
	CustomCSS   *string                   `json:"custom_css,omitempty"`
	Status      models.PortfolioStatus    `json:"status"`
	IsPublic    bool                      `json:"is_public"`
	ViewCount   int                       `json:"view_count"`
	SEO         *models.PortfolioSEO      `json:"seo,omitempty"`
	Settings    *models.PortfolioSettings `json:"settings,omitempty"`
	Sections    []PortfolioSectionResponse `json:"sections,omitempty"`
	CreatedAt   string                    `json:"created_at"`
	UpdatedAt   string                    `json:"updated_at"`
}

type PortfolioSectionResponse struct {
	ID          string `json:"id"`
	PortfolioID string `json:"portfolio_id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Content     *string `json:"content,omitempty"`
	Settings    *string `json:"settings,omitempty"`
	SortOrder   int    `json:"sort_order"`
	IsVisible   bool   `json:"is_visible"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PublicPortfolioResponse struct {
	Title       string                    `json:"title"`
	Slug        string                    `json:"slug"`
	Tagline     *string                   `json:"tagline,omitempty"`
	Bio         *string                   `json:"bio,omitempty"`
	CoverImage  *string                   `json:"cover_image,omitempty"`
	Logo        *string                   `json:"logo,omitempty"`
	Template    string                    `json:"template"`
	Theme       *models.PortfolioTheme    `json:"theme,omitempty"`
	CustomCSS   *string                   `json:"custom_css,omitempty"`
	SEO         *models.PortfolioSEO      `json:"seo,omitempty"`
	Settings    *models.PortfolioSettings `json:"settings,omitempty"`
	Sections    []PortfolioSectionResponse `json:"sections,omitempty"`
	User        *PublicUserInfo           `json:"user,omitempty"`
}

type PublicUserInfo struct {
	Name        string              `json:"name"`
	Username    string              `json:"username"`
	Headline    *string             `json:"headline,omitempty"`
	Image       *string             `json:"image,omitempty"`
	Location    *string             `json:"location,omitempty"`
	SocialMedia *models.SocialMedia `json:"social_media,omitempty"`
}
