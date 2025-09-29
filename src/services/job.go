package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"strings"

	"gorm.io/gorm"
)

type JobService struct {
	database *gorm.DB
}

func NewJobService(database *gorm.DB) *JobService {
	return &JobService{
		database: database,
	}
}

func (s *JobService) CreateJob(id string, payload dto.CreateJobDto) (*models.Job, error) {
	auth := NewAuthService(s.database)
	user, err := auth.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !user.IsRecruiter {
		return nil, errors.New("you need to be a recruiter to post jobs")
	}

	if payload.Salary.Min == 0 || payload.Salary.Max == 0 {
		return nil, errors.New("zero value not allowed")
	}

	if payload.Salary.Min >= payload.Salary.Max {
		return nil, errors.New("minimum salary cannot be higher than maximum salary")
	}

	if len(payload.Requirements) == 0 {
		return nil, errors.New("you need to add a least one requirement")
	}

	job := &models.Job{
		Title:          payload.Title,
		Location:       payload.Location,
		Description:    payload.Description,
		Deadline:       payload.Deadline,
		Requirements:   payload.Requirements,
		Salary:         payload.Salary,
		IsRemote:       payload.IsRemote,
		EmploymentType: payload.EmploymentType,
		CreatedBy:      user.ID,
	}

	if err := s.database.Create(&job).Error; err != nil {
		return nil, err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Job Created",
			Template: "job-created",
			Data: map[string]interface{}{
				"Name": []string{user.Username},
				"Otp":  user.Otp,
			},
		})
	}()

	return job, nil
}

func (s *JobService) UpdateJob(id string, payload dto.UpdateJobDto) (*models.Job, error) {
	job, err := s.FindJobById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	return &job, nil
}

func (s *JobService) DeleteJob(id string) error {
	job, err := s.FindJobById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("job not found")
		}
		return err
	}

	if err := s.database.Delete(&job).Error; err != nil {
		return err
	}

	return nil
}

func (s *JobService) GetJobs(params dto.JobPagination) (*dto.PaginatedResponse[models.Job], error) {
	q := normalizeJobQuery(params)
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	var jobs []models.Job
	var totalItems int64

	query := s.database.Model(&models.Job{})

	if q.Company != "" && strings.TrimSpace(q.Company) != "" {
		company := "%" + strings.ToLower(strings.TrimSpace(q.Company)) + "%"
		query = query.Where("LOWER(company) LIKE ?", company)
	}

	if q.EmploymentType != "" && strings.TrimSpace(q.EmploymentType) != "" {
		employmentType := "%" + strings.ToLower(strings.TrimSpace(q.EmploymentType)) + "%"
		query = query.Where("LOWER(employmentType) LIKE ?", employmentType)
	}

	if q.Location != "" && strings.TrimSpace(q.Location) != "" {
		location := "%" + strings.ToLower(strings.TrimSpace(q.Location)) + "%"
		query = query.Where("LOWER(location) LIKE ?", location)
	}

	if q.PostedDate != "" && strings.TrimSpace(q.PostedDate) != "" {
		query = query.Where("DATE(posted_date) = ?", q.PostedDate)
	}

	if q.Requirement != "" && strings.TrimSpace(q.Requirement) != "" {
		requirement := "%" + strings.ToLower(strings.TrimSpace(q.Requirement)) + "%"
		query = query.Where("LOWER(?) = ANY (SELECT LOWER(lang) FROM unnest(requirements) AS lang)", requirement)
	}

	if q.Salary != "" && strings.TrimSpace(q.Salary) != "" {
		query = query.Where("? BETWEEN min_salary AND max_salary", q.Salary)
	}

	query = query.Where("isRemote = ?", q.IsRemote)

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.Job]{
			Data:       []models.Job{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.Preload("Salary").Order("createdAt DESC").Offset(offset).
		Limit(params.Limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	totalPages := int64(0)
	if params.Limit > 0 {
		totalPages = (totalItems + int64(params.Limit) - 1) / int64(params.Limit)
	}

	return &dto.PaginatedResponse[models.Job]{
		Data:       jobs,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *JobService) GetJob(id string) (*models.Job, error) {
	var job *models.Job

	if err := s.database.Preload("Salary").Where("id = ?", id).First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	return job, nil
}

func (s *JobService) ApplyToJob(userId, jobId string, payload dto.JobApplicationDto) error {
	auth := NewAuthService(s.database)
	user, err := auth.FindUserById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	job, err := s.FindJobById(jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("job not found")
		}
		return err
	}

	if user.IsRecruiter {
		return errors.New("recruiters cannot apply to jobs")
	}

	var existing models.JobApplication
	if err := s.database.Where("user_id = ? AND job_id = ?", user.ID, job.ID).First(&existing).Error; err == nil {
		return errors.New("already applied to this job")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	application := &models.JobApplication{
		JobID:       job.ID,
		ApplicantID: user.ID,
		Resume:      "",
		CoverLetter: &payload.CoverLetter,
		Notes:       payload.Notes,
	}

	if err := s.database.Create(application).Error; err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Application Submitted",
			Template: "application-submitted",
			Data: map[string]interface{}{
				"Name": user.Username,
				"Job":  job.Title,
			},
		})
	}()

	return nil
}

func (s *JobService) FindJobById(id string) (models.Job, error) {
	var job models.Job

	if err := s.database.Where("id = ?", id).First(&job).Error; err != nil {
		return models.Job{}, err
	}

	return job, nil
}

func normalizeJobQuery(q dto.JobPagination) dto.JobPagination {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	return q
}
