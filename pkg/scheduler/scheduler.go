package scheduler

import (
	"errors"
	"fmt"
	"time"

	"github.com/praveennagaraj97/online-consultation/app"
)

type Scheduler struct {
	conf          *app.ApplicationConfig
	reminderTasks map[time.Time]*Task
}

func (s *Scheduler) Initialize(conf *app.ApplicationConfig) {
	s.conf = conf

	s.reminderTasks = make(map[time.Time]*Task, 0)
}

// Used to schedule tasks manually.
func (s *Scheduler) NewSchedule(invokeTime time.Time, name TasksTypes) error {

	timer := time.NewTimer(time.Until(invokeTime))

	fmt.Println(time.Until(invokeTime))

	if invokeTime.Unix() < time.Now().Unix() {
		return errors.New("Scheduling time is invalid")
	}

	switch name {
	case AppointmentReminderTask:
		if s.reminderTasks[invokeTime] != nil {
			timer.Stop()
			return nil
		}

		s.reminderTasks[invokeTime] = &Task{
			timer: timer,
		}

		go s.reminderTask(timer)
	default:
		timer.Stop()
	}

	return nil
}

func (s *Scheduler) Shutdown() {

	for _, value := range s.reminderTasks {
		value.timer.Stop()
	}

}
