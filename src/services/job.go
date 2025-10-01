package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"log"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobService struct {
	database     *gorm.DB
	notification *NotificationService
}

func NewJobService(database *gorm.DB, notification *NotificationService) *JobService {
	return &JobService{
		database:     database,
		notification: notification,
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
		return nil, errors.New("you need to add at least one requirement")
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
				"Name": user.Username,
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

	if err := s.database.Model(&job).Updates(payload).Error; err != nil {
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

	var jobs []models.Job
	var totalItems int64

	query := s.database.Model(&models.Job{})

	if q.Company != "" && strings.TrimSpace(q.Company) != "" {
		company := "%" + strings.ToLower(strings.TrimSpace(q.Company)) + "%"
		query = query.Where("LOWER(company) LIKE ?", company)
	}

	if q.EmploymentType != "" && strings.TrimSpace(q.EmploymentType) != "" {
		employmentType := "%" + strings.ToLower(strings.TrimSpace(q.EmploymentType)) + "%"
		query = query.Where("LOWER(employment_type) LIKE ?", employmentType)
	}

	if q.Location != "" && strings.TrimSpace(q.Location) != "" {
		location := "%" + strings.ToLower(strings.TrimSpace(q.Location)) + "%"
		query = query.Where("LOWER(location) LIKE ?", location)
	}

	if q.PostedDate != "" && strings.TrimSpace(q.PostedDate) != "" {
		query = query.Where("DATE(posted_date) = ?", q.PostedDate)
	}

	if q.Requirement != "" && strings.TrimSpace(q.Requirement) != "" {
		requirement := strings.ToLower(strings.TrimSpace(q.Requirement))
		query = query.Where("EXISTS (SELECT 1 FROM unnest(requirements) AS req WHERE LOWER(req) LIKE ?)", "%"+requirement+"%")
	}

	if q.Salary != "" && strings.TrimSpace(q.Salary) != "" {
		query = query.Where("? BETWEEN salary_min AND salary_max", q.Salary)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.Job]{
			Data:       []models.Job{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (q.Page - 1) * q.Limit

	if err := query.
		Preload("CreatedByUser").
		Preload("Comments").
		Preload("Comments.CreatedByUser").
		Preload("Reactions").
		Preload("Reactions.CreatedByUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(q.Limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	totalPages := (totalItems + int64(q.Limit) - 1) / int64(q.Limit)

	return &dto.PaginatedResponse[models.Job]{
		Data:       jobs,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       q.Page,
		Limit:      q.Limit,
	}, nil
}

func (s *JobService) GetJobsByUser(id string, params dto.Pagination) (*dto.PaginatedResponse[models.Job], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var jobs []models.Job
	var totalItems int64

	query := s.database.Model(&models.Job{}).Where("created_by = ?", id)

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

	if err := query.
		Preload("CreatedByUser").
		Preload("Comments").
		Preload("Comments.CreatedByUser").
		Preload("Reactions").
		Preload("Reactions.CreatedByUser").
		Preload("Applications").
		Preload("Applications.Job").
		Preload("Applications.Applicant").
		Order("created_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	totalPages := (totalItems + int64(params.Limit) - 1) / int64(params.Limit)

	return &dto.PaginatedResponse[models.Job]{
		Data:       jobs,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *JobService) GetJob(id string) (*models.Job, error) {
	var job models.Job

	if err := s.database.
		Preload("CreatedByUser").
		Preload("Comments").
		Preload("Comments.CreatedByUser").
		Preload("Reactions").
		Preload("Reactions.CreatedByUser").
		Where("id = ?", id).
		First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	return &job, nil
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
	if err := s.database.Where("applicant_id = ? AND job_id = ?", user.ID, job.ID).First(&existing).Error; err == nil {
		return errors.New("already applied to this job")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	application := &models.JobApplication{
		JobID:       job.ID,
		ApplicantID: user.ID,
		Resume:      payload.Resume,
		CoverLetter: &payload.CoverLetter,
		Notes:       payload.Notes,
	}

	if err := s.database.Create(application).Error; err != nil {
		return err
	}

	go func() {
		if err := s.notification.NotifyJobApplication(
			job.CreatedBy.String(),
			userId,
			jobId,
			job.Title,
			user.Name,
		); err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}()

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

func (s *JobService) GetApplicationsByJob(recruiterId, jobId string, params dto.Pagination) (*dto.PaginatedResponse[models.JobApplication], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	recruiterUUID, err := uuid.Parse(recruiterId)
	if err != nil {
		return nil, errors.New("invalid recruiter ID")
	}

	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		return nil, errors.New("invalid job ID")
	}

	auth := NewAuthService(s.database)
	recruiter, err := auth.FindUserById(recruiterId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !recruiter.IsRecruiter {
		return nil, errors.New("only recruiters can view applications")
	}

	var job models.Job
	if err := s.database.Where("id = ? AND created_by = ?", jobUUID, recruiterUUID).First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found or unauthorized")
		}
		return nil, err
	}

	var applications []models.JobApplication
	var totalItems int64

	query := s.database.Model(&models.JobApplication{}).Where("job_id = ?", jobUUID)

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.JobApplication]{
			Data:       []models.JobApplication{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("Applicant").
		Preload("Job").
		Order("created_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&applications).Error; err != nil {
		return nil, err
	}

	totalPages := (totalItems + int64(params.Limit) - 1) / int64(params.Limit)

	return &dto.PaginatedResponse[models.JobApplication]{
		Data:       applications,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *JobService) UpdateApplicationStatus(recruiterId, applicationId, status string, reason *string) (*models.JobApplication, error) {
	recruiterUUID, err := uuid.Parse(recruiterId)
	if err != nil {
		return nil, errors.New("invalid recruiter ID")
	}

	applicationUUID, err := uuid.Parse(applicationId)
	if err != nil {
		return nil, errors.New("invalid application ID")
	}

	auth := NewAuthService(s.database)
	recruiter, err := auth.FindUserById(recruiterId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !recruiter.IsRecruiter {
		return nil, errors.New("only recruiters can update application status")
	}

	var application models.JobApplication
	if err := s.database.Preload("Job").Preload("Applicant").Where("id = ?", applicationUUID).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	if application.Job.CreatedBy != recruiterUUID {
		return nil, errors.New("you can only update applications for your own jobs")
	}

	validStatuses := map[string]bool{
		"pending":  true,
		"reviewed": true,
		"accepted": true,
		"rejected": true,
		"hired":    true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status. must be: pending, reviewed, accepted, rejected, or hired")
	}

	application.Status = status
	if reason != nil {
		application.Notes = reason
	}

	if err := s.database.Save(&application).Error; err != nil {
		return nil, err
	}

	go func() {
		statusMessage := "updated"
		switch status {
		case "accepted":
			statusMessage = "accepted"
		case "rejected":
			statusMessage = "rejected"
		case "hired":
			statusMessage = "hired"
		case "reviewed":
			statusMessage = "reviewed"
		}

		emailData := map[string]interface{}{
			"Name":   application.Applicant.Username,
			"Job":    application.Job.Title,
			"Status": statusMessage,
		}

		if reason != nil && *reason != "" {
			emailData["Reason"] = *reason
		}

		lib.SendEmail(lib.EmailDto{
			To:       []string{application.Applicant.Email},
			Subject:  "Application Status Updated",
			Template: "application-status-update",
			Data:     emailData,
		})
	}()

	go func() {
		var err error
		switch application.Status {
		case "accepted":
			err = s.notification.NotifyApplicationAccepted(
				application.ApplicantID.String(),
				recruiterId,
				application.JobID.String(),
				application.Job.Title,
			)
		case "rejected":
			err = s.notification.NotifyApplicationRejected(
				application.ApplicantID.String(),
				recruiterId,
				application.JobID.String(),
				application.Job.Title,
			)
		}

		if err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}()

	if err := s.database.Preload("Job").Preload("Applicant").First(&application, application.ID).Error; err != nil {
		return nil, err
	}

	return &application, nil
}

func (s *JobService) AcceptApplication(recruiterId, applicationId string, reason *string) (*models.JobApplication, error) {
	return s.UpdateApplicationStatus(recruiterId, applicationId, "accepted", reason)
}

func (s *JobService) RejectApplication(recruiterId, applicationId string, reason *string) (*models.JobApplication, error) {
	return s.UpdateApplicationStatus(recruiterId, applicationId, "rejected", reason)
}

func (s *JobService) AddComment(userId, jobId string, content string) (*models.Comment, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		return nil, errors.New("invalid job ID")
	}

	if _, err := s.FindJobById(jobId); err != nil {
		return nil, errors.New("job not found")
	}

	comment := &models.Comment{
		Content:   content,
		JobID:     jobUUID,
		CreatedBy: userUUID,
	}

	if err := s.database.Create(comment).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("CreatedByUser").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *JobService) DeleteComment(commentId, userId string) error {
	commentUUID, err := uuid.Parse(commentId)
	if err != nil {
		return errors.New("invalid comment ID")
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}

	var comment models.Comment
	if err := s.database.Where("id = ? AND created_by = ?", commentUUID, userUUID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment not found or unauthorized")
		}
		return err
	}

	if err := s.database.Delete(&comment).Error; err != nil {
		return err
	}

	return nil
}

func (s *JobService) AddReaction(userId, jobId string, reactionType models.ReactionType) (*models.Reaction, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		return nil, errors.New("invalid job ID")
	}

	if _, err := s.FindJobById(jobId); err != nil {
		return nil, errors.New("job not found")
	}

	var existing models.Reaction
	if err := s.database.Where("job_id = ? AND created_by = ?", jobUUID, userUUID).First(&existing).Error; err == nil {
		if existing.Type == reactionType {
			if err = s.database.Delete(&existing).Error; err != nil {
				return nil, err
			}
			return nil, nil
		}

		existing.Type = reactionType
		if err = s.database.Save(&existing).Error; err != nil {
			return nil, err
		}

		if err = s.database.Preload("CreatedByUser").First(&existing, existing.ID).Error; err != nil {
			return nil, err
		}

		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	reaction := &models.Reaction{
		Type:      reactionType,
		JobID:     jobUUID,
		CreatedBy: userUUID,
	}

	if err := s.database.Create(reaction).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("CreatedByUser").First(reaction, reaction.ID).Error; err != nil {
		return nil, err
	}

	return reaction, nil
}

func (s *JobService) RemoveReaction(userId, jobId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}

	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		return errors.New("invalid job ID")
	}

	var reaction models.Reaction
	if err := s.database.Where("job_id = ? AND created_by = ?", jobUUID, userUUID).First(&reaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reaction not found")
		}
		return err
	}

	if err := s.database.Delete(&reaction).Error; err != nil {
		return err
	}

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
