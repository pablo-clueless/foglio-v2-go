package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnnouncementService struct {
	database *gorm.DB
	hub      *lib.Hub
}

func NewAnnouncementService(database *gorm.DB, hub *lib.Hub) *AnnouncementService {
	return &AnnouncementService{
		database: database,
		hub:      hub,
	}
}

// CreateAnnouncement creates a new announcement (admin only)
func (s *AnnouncementService) CreateAnnouncement(adminID string, payload dto.CreateAnnouncementDto) (*models.Announcement, error) {
	announcement := &models.Announcement{
		Title:          payload.Title,
		Content:        payload.Content,
		Type:           payload.Type,
		TargetAudience: payload.TargetAudience,
		Priority:       payload.Priority,
		ShowAsBanner:   payload.ShowAsBanner,
		BannerColor:    payload.BannerColor,
		ActionURL:      payload.ActionURL,
		ActionText:     payload.ActionText,
		ScheduledAt:    payload.ScheduledAt,
		ExpiresAt:      payload.ExpiresAt,
		Metadata:       payload.Metadata,
		CreatedBy:      uuid.Must(uuid.Parse(adminID)),
	}

	// Set default priority if not specified
	if announcement.Priority == "" {
		announcement.Priority = models.PriorityNormal
	}

	// If publish_now is true, publish immediately
	if payload.PublishNow {
		now := time.Now()
		announcement.IsPublished = true
		announcement.PublishedAt = &now
	}

	if err := s.database.Create(announcement).Error; err != nil {
		return nil, err
	}

	// If published and not scheduled for future, broadcast to targeted users
	if announcement.IsPublished && (announcement.ScheduledAt == nil || announcement.ScheduledAt.Before(time.Now())) {
		go s.broadcastAnnouncement(announcement)
	}

	return announcement, nil
}

// UpdateAnnouncement updates an existing announcement
func (s *AnnouncementService) UpdateAnnouncement(id string, payload dto.UpdateAnnouncementDto) (*models.Announcement, error) {
	var announcement models.Announcement
	if err := s.database.Where("id = ?", id).First(&announcement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("announcement not found")
		}
		return nil, err
	}

	// Apply updates
	if payload.Title != nil {
		announcement.Title = *payload.Title
	}
	if payload.Content != nil {
		announcement.Content = *payload.Content
	}
	if payload.Type != nil {
		announcement.Type = *payload.Type
	}
	if payload.TargetAudience != nil {
		announcement.TargetAudience = *payload.TargetAudience
	}
	if payload.Priority != nil {
		announcement.Priority = *payload.Priority
	}
	if payload.ShowAsBanner != nil {
		announcement.ShowAsBanner = *payload.ShowAsBanner
	}
	if payload.BannerColor != nil {
		announcement.BannerColor = payload.BannerColor
	}
	if payload.ActionURL != nil {
		announcement.ActionURL = payload.ActionURL
	}
	if payload.ActionText != nil {
		announcement.ActionText = payload.ActionText
	}
	if payload.ScheduledAt != nil {
		announcement.ScheduledAt = payload.ScheduledAt
	}
	if payload.ExpiresAt != nil {
		announcement.ExpiresAt = payload.ExpiresAt
	}
	if payload.Metadata != nil {
		announcement.Metadata = payload.Metadata
	}

	if err := s.database.Save(&announcement).Error; err != nil {
		return nil, err
	}

	return &announcement, nil
}

// PublishAnnouncement publishes a draft announcement
func (s *AnnouncementService) PublishAnnouncement(id string) (*models.Announcement, error) {
	var announcement models.Announcement
	if err := s.database.Where("id = ?", id).First(&announcement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("announcement not found")
		}
		return nil, err
	}

	if announcement.IsPublished {
		return nil, errors.New("announcement is already published")
	}

	now := time.Now()
	announcement.IsPublished = true
	announcement.PublishedAt = &now

	if err := s.database.Save(&announcement).Error; err != nil {
		return nil, err
	}

	// Broadcast if not scheduled for future
	if announcement.ScheduledAt == nil || announcement.ScheduledAt.Before(now) {
		go s.broadcastAnnouncement(&announcement)
	}

	return &announcement, nil
}

// UnpublishAnnouncement unpublishes an announcement
func (s *AnnouncementService) UnpublishAnnouncement(id string) (*models.Announcement, error) {
	var announcement models.Announcement
	if err := s.database.Where("id = ?", id).First(&announcement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("announcement not found")
		}
		return nil, err
	}

	announcement.IsPublished = false
	announcement.PublishedAt = nil

	if err := s.database.Save(&announcement).Error; err != nil {
		return nil, err
	}

	return &announcement, nil
}

// DeleteAnnouncement soft-deletes an announcement
func (s *AnnouncementService) DeleteAnnouncement(id string) error {
	var announcement models.Announcement
	if err := s.database.Where("id = ?", id).First(&announcement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("announcement not found")
		}
		return err
	}

	return s.database.Delete(&announcement).Error
}

// GetAnnouncement retrieves a single announcement
func (s *AnnouncementService) GetAnnouncement(id string) (*models.Announcement, error) {
	var announcement models.Announcement
	if err := s.database.Preload("CreatedByUser").Where("id = ?", id).First(&announcement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("announcement not found")
		}
		return nil, err
	}

	return &announcement, nil
}

// GetAnnouncementsForAdmin returns all announcements for admin management
func (s *AnnouncementService) GetAnnouncementsForAdmin(params dto.AnnouncementQueryParams) (*dto.AdminAnnouncementListResponse, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var announcements []models.Announcement
	var totalItems int64

	query := s.database.Model(&models.Announcement{}).Preload("CreatedByUser")

	// Apply filters
	if params.Type != nil && *params.Type != "" {
		query = query.Where("type = ?", *params.Type)
	}
	if params.TargetAudience != nil && *params.TargetAudience != "" {
		query = query.Where("target_audience = ?", *params.TargetAudience)
	}
	if params.IsPublished != nil {
		query = query.Where("is_published = ?", *params.IsPublished)
	}
	if !params.IncludeExpired {
		query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&announcements).Error; err != nil {
		return nil, err
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))
	}

	return &dto.AdminAnnouncementListResponse{
		Data:       announcements,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

// GetAnnouncementsForUser returns announcements visible to the user
func (s *AnnouncementService) GetAnnouncementsForUser(user *models.User, page, limit int) (*dto.AnnouncementListResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	now := time.Now()
	var announcements []models.Announcement
	var totalItems int64

	// Build audience filter based on user type
	audienceFilter := []models.TargetAudience{models.TargetAllUsers}
	if user.IsAdmin {
		audienceFilter = append(audienceFilter, models.TargetAdminsOnly)
	}
	if user.IsRecruiter {
		audienceFilter = append(audienceFilter, models.TargetRecruitersOnly)
	}
	if user.IsPremium {
		audienceFilter = append(audienceFilter, models.TargetPremiumOnly)
	}

	query := s.database.Model(&models.Announcement{}).
		Where("is_published = ?", true).
		Where("target_audience IN ?", audienceFilter).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", now).
		Where("(expires_at IS NULL OR expires_at > ?)", now)

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&announcements).Error; err != nil {
		return nil, err
	}

	// Fetch user statuses for these announcements
	announcementIDs := make([]uuid.UUID, len(announcements))
	for i, a := range announcements {
		announcementIDs[i] = a.ID
	}

	var statuses []models.UserAnnouncementStatus
	if len(announcementIDs) > 0 {
		s.database.Where("user_id = ? AND announcement_id IN ?", user.ID, announcementIDs).Find(&statuses)
	}

	statusMap := make(map[uuid.UUID]*models.UserAnnouncementStatus)
	for i := range statuses {
		statusMap[statuses[i].AnnouncementID] = &statuses[i]
	}

	// Build response
	responses := make([]dto.AnnouncementResponse, len(announcements))
	for i, a := range announcements {
		resp := dto.AnnouncementResponse{Announcement: a}
		if status, ok := statusMap[a.ID]; ok {
			resp.UserStatus = &dto.UserAnnouncementStatusResponse{
				IsRead:      status.IsRead,
				ReadAt:      status.ReadAt,
				IsDismissed: status.IsDismissed,
				DismissedAt: status.DismissedAt,
			}
		}
		responses[i] = resp
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}

	return &dto.AnnouncementListResponse{
		Data:       responses,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetActiveBanners returns active banner announcements for the user
func (s *AnnouncementService) GetActiveBanners(user *models.User) ([]dto.BannerAnnouncementResponse, error) {
	now := time.Now()
	var announcements []models.Announcement

	audienceFilter := []models.TargetAudience{models.TargetAllUsers}
	if user.IsAdmin {
		audienceFilter = append(audienceFilter, models.TargetAdminsOnly)
	}
	if user.IsRecruiter {
		audienceFilter = append(audienceFilter, models.TargetRecruitersOnly)
	}
	if user.IsPremium {
		audienceFilter = append(audienceFilter, models.TargetPremiumOnly)
	}

	err := s.database.
		Where("is_published = ?", true).
		Where("show_as_banner = ?", true).
		Where("target_audience IN ?", audienceFilter).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", now).
		Where("(expires_at IS NULL OR expires_at > ?)", now).
		Order("priority DESC, created_at DESC").
		Find(&announcements).Error

	if err != nil {
		return nil, err
	}

	// Get dismissed statuses
	announcementIDs := make([]uuid.UUID, len(announcements))
	for i, a := range announcements {
		announcementIDs[i] = a.ID
	}

	var statuses []models.UserAnnouncementStatus
	if len(announcementIDs) > 0 {
		s.database.Where("user_id = ? AND announcement_id IN ? AND is_dismissed = ?", user.ID, announcementIDs, true).Find(&statuses)
	}

	dismissedMap := make(map[uuid.UUID]bool)
	for _, status := range statuses {
		dismissedMap[status.AnnouncementID] = true
	}

	// Filter out dismissed banners and build response
	var banners []dto.BannerAnnouncementResponse
	for _, a := range announcements {
		isDismissed := dismissedMap[a.ID]
		if !isDismissed {
			banners = append(banners, dto.BannerAnnouncementResponse{
				ID:          a.ID.String(),
				Title:       a.Title,
				Content:     a.Content,
				Type:        a.Type,
				Priority:    a.Priority,
				BannerColor: a.BannerColor,
				ActionURL:   a.ActionURL,
				ActionText:  a.ActionText,
				IsDismissed: false,
			})
		}
	}

	return banners, nil
}

// MarkAnnouncementAsRead marks an announcement as read for a user
func (s *AnnouncementService) MarkAnnouncementAsRead(userID, announcementID string) error {
	var status models.UserAnnouncementStatus

	err := s.database.Where("user_id = ? AND announcement_id = ?", userID, announcementID).First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new status
			now := time.Now()
			status = models.UserAnnouncementStatus{
				UserID:         uuid.Must(uuid.Parse(userID)),
				AnnouncementID: uuid.Must(uuid.Parse(announcementID)),
				IsRead:         true,
				ReadAt:         &now,
			}
			return s.database.Create(&status).Error
		}
		return err
	}

	if !status.IsRead {
		now := time.Now()
		status.IsRead = true
		status.ReadAt = &now
		return s.database.Save(&status).Error
	}

	return nil
}

// DismissAnnouncement dismisses an announcement for a user
func (s *AnnouncementService) DismissAnnouncement(userID, announcementID string) error {
	var status models.UserAnnouncementStatus

	err := s.database.Where("user_id = ? AND announcement_id = ?", userID, announcementID).First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now()
			status = models.UserAnnouncementStatus{
				UserID:         uuid.Must(uuid.Parse(userID)),
				AnnouncementID: uuid.Must(uuid.Parse(announcementID)),
				IsDismissed:    true,
				DismissedAt:    &now,
			}
			return s.database.Create(&status).Error
		}
		return err
	}

	if !status.IsDismissed {
		now := time.Now()
		status.IsDismissed = true
		status.DismissedAt = &now
		return s.database.Save(&status).Error
	}

	return nil
}

// broadcastAnnouncement sends the announcement to all targeted connected users
func (s *AnnouncementService) broadcastAnnouncement(announcement *models.Announcement) {
	// Get all users matching the target audience
	var users []models.User
	query := s.database.Model(&models.User{})

	switch announcement.TargetAudience {
	case models.TargetAllUsers:
		// No filter needed
	case models.TargetAdminsOnly:
		query = query.Where("is_admin = ?", true)
	case models.TargetRecruitersOnly:
		query = query.Where("is_recruiter = ?", true)
	case models.TargetPremiumOnly:
		query = query.Where("is_premium = ?", true)
	}

	if err := query.Find(&users).Error; err != nil {
		log.Printf("Failed to fetch users for announcement broadcast: %v", err)
		return
	}

	// Send notification to each targeted user
	for _, user := range users {
		notification := models.Notification{
			ID:      uuid.New(),
			Title:   announcement.Title,
			Content: announcement.Content,
			Type:    models.System,
			OwnerID: user.ID,
			IsRead:  false,
			Data: map[string]interface{}{
				"announcement_id":   announcement.ID.String(),
				"announcement_type": announcement.Type,
				"priority":          announcement.Priority,
				"show_as_banner":    announcement.ShowAsBanner,
			},
		}

		// Save to database
		if err := s.database.Create(&notification).Error; err != nil {
			log.Printf("Failed to create notification for user %s: %v", user.ID, err)
			continue
		}

		// Send via WebSocket if hub is available
		if s.hub != nil {
			s.hub.SendToUser(user.ID.String(), notification)
		}
	}
}

// ProcessScheduledAnnouncements publishes scheduled announcements (to be called by a cron job)
func (s *AnnouncementService) ProcessScheduledAnnouncements() error {
	now := time.Now()
	var announcements []models.Announcement

	err := s.database.
		Where("is_published = ?", false).
		Where("scheduled_at IS NOT NULL").
		Where("scheduled_at <= ?", now).
		Find(&announcements).Error

	if err != nil {
		return err
	}

	for _, a := range announcements {
		a.IsPublished = true
		a.PublishedAt = &now
		if err := s.database.Save(&a).Error; err != nil {
			log.Printf("Failed to publish scheduled announcement %s: %v", a.ID, err)
			continue
		}
		go s.broadcastAnnouncement(&a)
	}

	return nil
}
