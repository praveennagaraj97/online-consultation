package scheduler

import (
	"errors"
	"time"

	appointmentrepository "github.com/praveennagaraj97/online-consultation/repository/appointment"
)

type TasksTypes string

type Scheduler struct {
	appointmentReminderTasks map[time.Time]*Task
	apptScheduledRepo        *appointmentrepository.AppointmentScheduleReminderRepository

	jobs map[string]*time.Ticker
}

// Starts when the app starts
func (s *Scheduler) Initialize() {

	// Appointment Schedule reminder tasks
	s.appointmentReminderTasks = make(map[time.Time]*Task, 0)
	s.jobs = make(map[string]*time.Ticker, 0)

	// Start Remider tasks when app starts
	everyDayTaskChan := make(chan bool, 1)
	go s.runEveryDay(everyDayTaskChan)

	select {
	case <-everyDayTaskChan:
		close(everyDayTaskChan)
	default:
		close(everyDayTaskChan)
	}

}

func (s *Scheduler) InitializeAppointmentRemainderPersistRepo(
	apptScheduledRepo *appointmentrepository.AppointmentScheduleReminderRepository) {
	s.apptScheduledRepo = apptScheduledRepo
}

// Used to schedule tasks manually.
func (s *Scheduler) NewSchedule(invokeTime time.Time, name TasksTypes) error {

	if invokeTime.Unix() < time.Now().Unix() {
		return errors.New("scheduling time is invalid")
	}

	// Create a invoke timer.
	timer := time.NewTimer(time.Until(invokeTime))

	switch name {
	case AppointmentReminderTask:
		// If an reminder is already assigned at invoke time stop the created timer.
		if s.appointmentReminderTasks[invokeTime] != nil {
			timer.Stop()
			return nil
		}

		// Create a new task with given time
		s.appointmentReminderTasks[invokeTime] = &Task{
			timer: timer,
		}
		// Fire the go routine which will end once invoke time is over
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
