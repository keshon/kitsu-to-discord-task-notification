package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID               uint `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	TaskID           string
	TaskUpdatedAt    string
	TaskStatus       string
	CommentID        string
	CommentUpdatedAt string
}

func CreateTask(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	db.Create(&Task{TaskID: taskID, TaskUpdatedAt: taskUpdatedAt, TaskStatus: taskStatus, CommentUpdatedAt: commentUpdatedAt, CommentID: commentID})
}

func UpdateTask(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	var rec Task
	db.Where("task_id=?", taskID).Find(&rec)
	rec.TaskUpdatedAt = taskUpdatedAt
	rec.TaskStatus = taskStatus
	rec.CommentUpdatedAt = commentUpdatedAt
	rec.CommentID = commentID
	db.Save(&rec)
}

func FindTask(db *gorm.DB, taskID string) Task {
	var Task Task
	db.First(&Task, "task_id = ?", taskID)
	return Task
}
