package job

import (
	"ai_summary_project/internal/repository"
	"ai_summary_project/pkg/jwt"
	"ai_summary_project/pkg/log"
	"ai_summary_project/pkg/sid"
)

type Job struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
