package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"
	"strings"

	"gorm.io/gorm"
)

type UserService struct {
	database *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		database: database,
	}
}

func (s *UserService) GetUsers(params dto.UserPagination) (*dto.PaginatedResponse[models.User], error) {
	q := normalizeUserQuery(params)

	var users []models.User
	var totalItems int64
	query := s.database.Model(&models.User{})

	if q.Query != nil && strings.TrimSpace(*q.Query) != "" {
		searchTerm := strings.ToLower(strings.TrimSpace(*q.Query))
		query = query.Where(
			s.database.Where("LOWER(location) LIKE ?", "%"+searchTerm+"%").
				Or("LOWER(?) = ANY (SELECT LOWER(unnest(experiences)))", searchTerm).
				Or("LOWER(?) = ANY (SELECT LOWER(unnest(languages)))", searchTerm).
				Or("LOWER(?) = ANY (SELECT LOWER(unnest(skills)))", searchTerm).
				Or("LOWER(?) = ANY (SELECT LOWER(role))", searchTerm).
				Or("LOWER(role) LIKE ?", "%"+searchTerm+"%").
				Or("LOWER(headline) LIKE ?", "%"+searchTerm+"%").
				Or("LOWER(summary) LIKE ?", "%"+searchTerm+"%"),
		)
	}

	if q.UserType != nil && strings.TrimSpace(*q.UserType) != "" {
		userType := strings.ToLower(strings.TrimSpace(*q.UserType))
		if userType != "all" {
			switch userType {
			case "recruiter":
				query = query.Where("is_recruiter = ?", true)
			case "talent", "talents":
				query = query.Where("is_recruiter = ?", false)
			}
		}
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.User]{
			Data:       []models.User{},
			Limit:      q.Limit,
			Page:       q.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (q.Page - 1) * q.Limit
	if err := query.
		Preload("Projects").
		Preload("Experiences").
		Preload("Education").
		Preload("Certifications").
		Preload("Languages").
		Offset(offset).
		Order("created_at DESC").
		Limit(q.Limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int64(0)
	if q.Limit > 0 {
		totalPages = (totalItems + int64(q.Limit) - 1) / int64(q.Limit)
	}

	return &dto.PaginatedResponse[models.User]{
		Data:       users,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       q.Page,
		Limit:      q.Limit,
	}, nil
}

func (s *UserService) GetUser(idOrUsername string) (*models.User, error) {
	var user *models.User

	if err := s.database.Preload("Projects").Preload("Experiences").Preload("Education").Preload("Certifications").Preload("Languages").Preload("Company").
		Where("id = ? OR LOWER(username) = LOWER(?)", idOrUsername, idOrUsername).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(id string, payload dto.UpdateUserDto) (*models.User, error) {
	auth := NewAuthService(s.database)
	user, err := auth.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if payload.Headline != nil {
		user.Headline = payload.Headline
	}
	if payload.Location != nil {
		user.Location = payload.Location
	}
	if payload.Phone != nil {
		user.Phone = payload.Phone
	}
	if payload.Summary != nil {
		user.Summary = payload.Summary
	}
	if payload.SocialMedia != nil {
		user.SocialMedia = payload.SocialMedia
	}
	if payload.Skills != nil {
		user.Skills = payload.Skills
	}

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	// Handle associations using GORM's Association API
	if payload.Projects != nil {
		for i := range payload.Projects {
			payload.Projects[i].UserID = user.ID
		}
		if err := s.database.Model(&user).Association("Projects").Replace(payload.Projects); err != nil {
			return nil, err
		}
	}
	if payload.Experiences != nil {
		for i := range payload.Experiences {
			payload.Experiences[i].UserID = user.ID
		}
		if err := s.database.Model(&user).Association("Experiences").Replace(payload.Experiences); err != nil {
			return nil, err
		}
	}
	if payload.Education != nil {
		for i := range payload.Education {
			payload.Education[i].UserID = user.ID
		}
		if err := s.database.Model(&user).Association("Education").Replace(payload.Education); err != nil {
			return nil, err
		}
	}
	if payload.Certifications != nil {
		for i := range payload.Certifications {
			payload.Certifications[i].UserID = user.ID
		}
		if err := s.database.Model(&user).Association("Certifications").Replace(payload.Certifications); err != nil {
			return nil, err
		}
	}
	if payload.Languages != nil {
		for i := range payload.Languages {
			payload.Languages[i].UserID = user.ID
		}
		if err := s.database.Model(&user).Association("Languages").Replace(payload.Languages); err != nil {
			return nil, err
		}
	}

	return s.GetUser(id)
}

func (s *UserService) UpdateAvatar(id, imageUrl string) (*models.User, error) {
	auth := NewAuthService(s.database)
	user, err := auth.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Image = &imageUrl

	if err := s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	user, err := NewAuthService(s.database).FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := s.database.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func normalizeUserQuery(q dto.UserPagination) dto.UserPagination {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	return q
}
