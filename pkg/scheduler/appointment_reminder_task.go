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

// Run EveryDay once
func (s *Scheduler) runEveryDay(ch chan bool) {

	// Initail set time until mid night 1 (IST)
	t := time.NewTicker(time.Until(getNextDateWithZeroTime()))
	s.jobs["appointment_reminder_job"] = t

	// Resent the timer when it reached it delivery time
	for range t.C {

		// Fire Tasks
		tim := time.NewTimer(time.Second * 3)

		go s.appointmentReminderTask(tim, s.appointmentReminderTasks, time.Now())

		t.Reset(time.Until(getNextDateWithZeroTime()))

	}

	ch <- true

}

// Send Reminders to Users
func (s *Scheduler) appointmentReminderTask(timer *time.Timer,
	appointmentReminderTasks map[time.Time]*Task, invokeTime time.Time) {

	// Wait till the timer reach given invoke time.
	<-timer.C

	// Send Email and SMS once timer is released
	invokeT := primitive.NewDateTimeFromTime(invokeTime)

	results, err := s.apptScheduledRepo.FindAllByInvokeTime(&invokeT)
	if err != nil {
		log.Default().Println(err.Error())
	}

	fmt.Println(results)

	delete(appointmentReminderTasks, invokeTime)

	defer timer.Reset(0)

}
