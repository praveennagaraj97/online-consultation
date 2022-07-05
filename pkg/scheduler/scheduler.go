package scheduler

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	appointmentrepository "github.com/praveennagaraj97/online-consultation/repository/appointment"
)

type TasksTypes string

type Scheduler struct {
	appointmentReminderTasks map[time.Time]*Task
	apptScheduledRepo        *appointmentrepository.AppointmentScheduleReminderRepository

	jobs map[string]*time.Ticker
}

func (s *Scheduler) Initialize() {

	s.appointmentReminderTasks = make(map[time.Time]*Task, 0)
	s.jobs = make(map[string]*time.Ticker, 0)

	apptReminderChan := make(chan bool, 1)
	go s.scheduleUpcomingAppointmentReminders(apptReminderChan)

	select {
	case <-apptReminderChan:
		close(apptReminderChan)
	default:
		close(apptReminderChan)
	}

}

func (s *Scheduler) InitializeAppointmentRemainderPersistRepo(apptScheduledRepo *appointmentrepository.AppointmentScheduleReminderRepository) {
	s.apptScheduledRepo = apptScheduledRepo
}

// Used to schedule tasks manually.
func (s *Scheduler) NewSchedule(invokeTime time.Time, name TasksTypes) error {

	if invokeTime.Unix() < time.Now().Unix() {
		return errors.New("scheduling time is invalid")
	}

	fmt.Printf("Number of open channels %v\n", runtime.NumGoroutine())

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
		if !value.timer.Stop() {
			// Drain and stop
			<-value.timer.C
			value.timer.Stop()
		}
	}
}
