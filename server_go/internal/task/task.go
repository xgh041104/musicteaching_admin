package task

import (
	"ai_summary_project/internal/repository"
	"ai_summary_project/pkg/jwt"
	"ai_summary_project/pkg/log"
	"ai_summary_project/pkg/sid"
)

type Task struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewTask(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Task {
	return &Task{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
