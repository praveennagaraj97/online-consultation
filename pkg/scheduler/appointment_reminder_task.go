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

type TasksTypes string

const (
	AppointmentReminderTask TasksTypes = "AppointmentReminderTask"
)

// Get
func (s *Scheduler) appointmentReminderTask(timer *time.Timer, appointmentReminderTasks map[time.Time]*Task, invokeTime time.Time) {
	defer delete(appointmentReminderTasks, invokeTime)
	defer timer.Stop()
	<-timer.C

	fmt.Println("Reminder sent")

	invokeT := primitive.NewDateTimeFromTime(invokeTime)

	results, err := s.apptScheduledRepo.FindAllByInvokeTime(&invokeT)
	if err != nil {
		log.Default().Println(err.Error())
	}

	fmt.Println(results)

}
