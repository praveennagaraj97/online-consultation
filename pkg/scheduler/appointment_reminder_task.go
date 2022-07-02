package scheduler

import (
	"fmt"
	"time"
)

type Task struct {
	timer *time.Timer
}

type TasksTypes string

const (
	AppointmentReminderTask TasksTypes = "AppointmentReminderTask"
)

func (s *Scheduler) appointmenrReminderTask(timer *time.Timer) {
	fmt.Println("lol")
	<-timer.C
	defer timer.Stop()

	// get list of slots having invoke time

}
