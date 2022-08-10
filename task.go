package main

import "time"

const (
	StatusFinished = "finished"
	StatusWaiting  = "waiting"
	StatusProgress = "progress"
)

type Task struct {
	Id          int64
	Description string
	Status      string
	CreatedAt   time.Time
}

func (task *Task) IsFinished() bool {
	return task.Status == StatusFinished
}

func (task *Task) IsValid() bool {
	return task.Description != "" &&
		len(task.Description) <= 255 &&
		(task.Status == StatusWaiting || task.Status == StatusProgress || task.Status == StatusFinished)
}

func (task Task) New(description, status string) *Task {
	task = Task{}
	task.Description = description
	task.Status = status

	return &task
}
