package scheduler

import (
	"time"
)

type Task struct {
	timer *time.Timer
}

type TasksTypes string

const (
	ReminderTask TasksTypes = "ReminderTask"
)

func (s *Scheduler) reminderTask(timer *time.Timer) {
	<-timer.C
	defer timer.Stop()

	// get list of slots having invoke time
}
