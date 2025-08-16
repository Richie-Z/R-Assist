package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func UseExistingDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	query := `
    CREATE TABLE IF NOT EXISTS logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender TEXT,
        message TEXT,
  			reply TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

		CREATE TABLE IF NOT EXISTS reminders (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				chat_jid TEXT NOT NULL,
				message TEXT NOT NULL,
				scheduled_at DATETIME NOT NULL,
				status TEXT NOT NULL DEFAULT 'pending',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				sent_at TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_reminders_due ON reminders(status, scheduled_at);
  `
	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create logs table: %w", err)
	}

	return nil
}
