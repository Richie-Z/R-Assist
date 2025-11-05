package service

import (
	"github.com/richie-z/R-Assist/internal/storage"
)

func LogMessage(sender, message string, reply string) error {
	_, err := storage.DB.Exec(
		"INSERT INTO logs (sender, message, reply) VALUES (?, ?, ?)",
		sender, message, reply,
	)
	return err
}
