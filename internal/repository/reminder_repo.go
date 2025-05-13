package repository

import (
	"time"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
)

type ReminderRepo interface {
	CreateReminder(model.Reminder) (string, error)
	GetRemindersByNote(string) ([]model.Reminder, error)
	UpdateReminder(string, string, string, time.Time, *bool, string) (int, error)
	DeleteReminder(string) (int, error)
}
