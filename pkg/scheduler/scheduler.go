package scheduler

import (
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/app"
)

type Scheduler struct {
	conf                     *app.ApplicationConfig
	appointmentReminderTasks map[time.Time]*Task
}

func (s *Scheduler) Initialize(conf *app.ApplicationConfig) {
	s.conf = conf

	s.appointmentReminderTasks = make(map[time.Time]*Task, 0)
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

		go s.appointmenrReminderTask(timer)
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
