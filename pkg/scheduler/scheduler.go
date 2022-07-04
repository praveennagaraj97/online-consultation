package scheduler

import (
	"errors"
	"fmt"
	"time"

	appointmentrepository "github.com/praveennagaraj97/online-consultation/repository/appointment"
)

type Scheduler struct {
	appointmentReminderTasks map[time.Time]*Task

	apptScheduledRepo *appointmentrepository.AppointmentScheduleReminderRepository
}

func (s *Scheduler) Initialize() {

	s.appointmentReminderTasks = make(map[time.Time]*Task, 0)
}

func (s *Scheduler) InitializeAppointmentRemainderPersistRepo(apptScheduledRepo *appointmentrepository.AppointmentScheduleReminderRepository) {
	s.apptScheduledRepo = apptScheduledRepo
}

// Used to schedule tasks manually.
func (s *Scheduler) NewSchedule(invokeTime time.Time, name TasksTypes) error {

	if invokeTime.Unix() < time.Now().Unix() {
		return errors.New("scheduling time is invalid")
	}

	timer := time.NewTimer(time.Until(invokeTime))

	switch name {
	case AppointmentReminderTask:
		if s.appointmentReminderTasks[invokeTime] != nil {
			timer.Stop()
			return nil
		}

		s.appointmentReminderTasks[invokeTime] = &Task{
			timer: timer,
		}

		go s.appointmentReminderTask(timer, s.appointmentReminderTasks, invokeTime)
	default:
		timer.Stop()
	}

	return nil
}

func (s *Scheduler) Shutdown() {
	for _, value := range s.appointmentReminderTasks {
		value.timer.Stop()
	}
}

func (s *Scheduler) StartScheduler(loopEvery time.Duration) {
	t := time.NewTicker(loopEvery)

	select {
	case <-t.C:
		fmt.Println("Get Todays reminders and assign to cron job")
	default:
		t.Stop()
	}
}
