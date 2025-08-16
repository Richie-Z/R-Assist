package bot

import (
	"context"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/richie-z/whatsapp-bot/internal/storage"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

func NewConnection() *whatsmeow.Client {
	dbLog := waLog.Stdout("DB", "INFO", true)
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:store.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}

	err = storage.UseExistingDB("store.db")
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	return client
}
