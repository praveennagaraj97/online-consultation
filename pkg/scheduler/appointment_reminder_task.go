package scheduler

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	timer *time.Timer
}

const (
	AppointmentReminderTask TasksTypes = "AppointmentReminderTask"
)

func (s *Scheduler) scheduleUpcomingAppointmentReminders(apptReminderChan chan bool) {

	t := time.NewTicker(time.Until(getNextDateWithZeroTime()))
	s.jobs["appointment_reminder_job"] = t

	for range t.C {
		nextDate := getNextDateWithZeroTime()
		t.Reset(time.Until(nextDate))

	}

	apptReminderChan <- true

}

// Get
func (s *Scheduler) appointmentReminderTask(timer *time.Timer, appointmentReminderTasks map[time.Time]*Task, invokeTime time.Time) {

	<-timer.C

	invokeT := primitive.NewDateTimeFromTime(invokeTime)

	results, err := s.apptScheduledRepo.FindAllByInvokeTime(&invokeT)
	if err != nil {
		log.Default().Println(err.Error())
	}

	fmt.Println(results)
	if timer.Stop() {
		delete(appointmentReminderTasks, invokeTime)
	} else {
		<-timer.C
		timer.Stop()
	}

}
