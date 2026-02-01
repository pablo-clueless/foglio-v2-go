package lib

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cron *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	c := cron.New(cron.WithSeconds())
	return &CronScheduler{cron: c}
}

func (s *CronScheduler) AddJob(schedule string, job func()) error {
	_, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return err
	}
	return nil
}

func (s *CronScheduler) Start() {
	s.cron.Start()
	log.Println("Cron scheduler started at", time.Now().Format(time.RFC3339))
}

func (s *CronScheduler) Stop() {
	s.cron.Stop()
	log.Println("Cron scheduler stopped")
}
