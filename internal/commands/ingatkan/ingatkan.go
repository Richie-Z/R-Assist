package ingatkan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/richie-z/R-Assist/internal/service"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	reDur   = regexp.MustCompile(`(?i)^(?P<msg>.+?)\s+(?P<num>\d+)\s*(?P<unit>menit|detik)\s+lagi$`)
	reJam   = regexp.MustCompile(`(?i)^(?P<msg>.+?)\s+jam\s+(?P<h>\d{1,2}):(?P<m>\d{2})$`)
	reBesok = regexp.MustCompile(`(?i)^(?P<msg>.+?)\s+besok(?:\s+jam\s+(?P<h>\d{1,2}):(?P<m>\d{2}))?$`)
)

func Run(client *whatsmeow.Client, evt *events.Message, args string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	remText, when, err := parseWhen(args, now, loc)
	if err != nil {
		service.SendReply(client, evt.Info.Sender, "Format tidak dikenali.\nContoh:\n"+
			"!ingatkan makan 20 menit lagi\n!ingatkan makan 5 detik lagi\n!ingatkan makan jam 10:00\n!ingatkan makan besok\n!ingatkan makan besok jam 18:00")
		return
	}

	// feedback to user
	confirm := fmt.Sprintf("Sip! Aku akan ingatkan: %s\nWaktu: %s", remText,
		when.In(loc).Format("Mon, 02 Jan 2006 15:04:05 MST"))
	service.SendReply(client, evt.Info.Sender, confirm)

	// schedule to DB
	if err := service.ScheduleReminder(evt.Info.Sender, remText, when); err != nil {
		service.SendReply(client, evt.Info.Sender, "Gagal menyimpan pengingat 😥")
		return
	}

	// _ = service.LogMessage(evt.Info.Sender.String(), "!ingatkan "+args)
}

func parseWhen(input string, now time.Time, loc *time.Location) (msg string, when time.Time, err error) {
	in := strings.TrimSpace(input)

	// 1) Duration: "<msg> 20menit lagi" or "5 detik lagi"
	if m := reDur.FindStringSubmatch(in); m != nil {
		msg = strings.TrimSpace(m[groupIndex(reDur, "msg")])
		n, _ := strconv.Atoi(m[groupIndex(reDur, "num")])
		unit := strings.ToLower(m[groupIndex(reDur, "unit")])
		d := time.Duration(n)
		if unit == "menit" {
			when = now.Add(d * time.Minute)
		} else { // detik
			when = now.Add(d * time.Second)
		}
		return
	}

	// 2) Absolute today: "<msg> jam 10:00" (if lewat, geser besok)
	if m := reJam.FindStringSubmatch(in); m != nil {
		msg = strings.TrimSpace(m[groupIndex(reJam, "msg")])
		h, _ := strconv.Atoi(m[groupIndex(reJam, "h")])
		min, _ := strconv.Atoi(m[groupIndex(reJam, "m")])
		target := time.Date(now.Year(), now.Month(), now.Day(), h, min, 0, 0, loc)
		if !target.After(now) {
			target = target.Add(24 * time.Hour)
		}
		when = target
		return
	}

	// 3) Besok: "<msg> besok" or "<msg> besok jam 18:00"
	if m := reBesok.FindStringSubmatch(in); m != nil {
		msg = strings.TrimSpace(m[groupIndex(reBesok, "msg")])
		hStr := m[groupIndex(reBesok, "h")]
		mStr := m[groupIndex(reBesok, "m")]

		h := now.Hour()
		min := now.Minute()
		if hStr != "" && mStr != "" {
			h, _ = strconv.Atoi(hStr)
			min, _ = strconv.Atoi(mStr)
		}
		target := time.Date(now.Year(), now.Month(), now.Day(), h, min, 0, 0, loc).Add(24 * time.Hour)
		when = target
		return
	}

	return "", time.Time{}, fmt.Errorf("unknown format")
}

func groupIndex(r *regexp.Regexp, name string) int {
	for i, n := range r.SubexpNames() {
		if n == name {
			return i
		}
	}
	return -1
}
