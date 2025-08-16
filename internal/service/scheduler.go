package service

import (
	"context"
	"fmt"
	"time"

	"github.com/richie-z/whatsapp-bot/internal/storage"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	pb "google.golang.org/protobuf/proto"
)

type Scheduler struct {
	client *whatsmeow.Client
	loc    *time.Location
	stop   chan struct{}
}

func NewScheduler(client *whatsmeow.Client, loc *time.Location) *Scheduler {
	return &Scheduler{client: client, loc: loc, stop: make(chan struct{})}
}

func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				s.tick()
			}
		}
	}()
}

func (s *Scheduler) Stop() { close(s.stop) }

func (s *Scheduler) tick() {
	tx, err := storage.DB.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`
		SELECT id, chat_jid, message
		FROM reminders
		WHERE status='pending' AND scheduled_at <= CURRENT_TIMESTAMP
		ORDER BY scheduled_at ASC
		LIMIT 20`)
	if err != nil {
		return
	}
	defer rows.Close()

	type item struct {
		id   int64
		jid  string
		text string
	}
	var due []item
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.id, &it.jid, &it.text); err == nil {
			due = append(due, it)
		}
	}
	if err := rows.Err(); err != nil {
		return
	}

	for _, d := range due {
		jid, _ := types.ParseJID(d.jid)
		msg := &proto.Message{Conversation: pb.String("⏰ Reminder: " + d.text)}
		_, sendErr := s.client.SendMessage(context.Background(), jid, msg)
		now := time.Now().In(s.loc)

		if sendErr == nil {
			_, _ = tx.Exec(`UPDATE reminders SET status='sent', sent_at=? WHERE id=?`, now, d.id)
		} else {
			fmt.Println("send reminder error:", sendErr)
			_, _ = tx.Exec(`UPDATE reminders SET status='error' WHERE id=?`, d.id)
		}
	}

	_ = tx.Commit()
}

func ScheduleReminder(chatJID types.JID, text string, when time.Time) error {
	_, err := storage.DB.Exec(
		`INSERT INTO reminders (chat_jid, message, scheduled_at) VALUES (?, ?, ?)`,
		chatJID.String(), text, when.UTC(),
	)
	return err
}
