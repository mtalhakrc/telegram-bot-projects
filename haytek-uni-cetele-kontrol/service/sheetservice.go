package service

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"haytekuni-cetele-kontrol/config"
	"haytekuni-cetele-kontrol/logx"
	"io"
	"log"
	"os"
)

var srv *sheets.Service
var cfg *config.Config

func Get() *sheets.Service {
	return srv
}

func getCredentialsAsBytes() ([]byte, error) {
	f, err := os.Open(cfg.Credentials.CredentialsPath)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func InitSheetService() {
	cfg = config.Get()
	b, err := getCredentialsAsBytes()
	if err != nil {
		logx.SendLog(err.Error())
		log.Fatalf("unable to read client secret file: %v", err)
	}
	srv, err = newSheetService(b)
	if err != nil {
		logx.SendLog(err.Error())
		log.Fatalf("unable to retrieve sheets client: %v", err)
	}

	if err = isSheetServiceGranted(); err != nil {
		logx.SendLog(err.Error())
		log.Fatalf("unable to retrieve sheets client: %v", err)
	}
	log.Println("Sheet service has initialized")
}

func newSheetService(b []byte) (*sheets.Service, error) {
	var err error
	ctx := context.Background()
	srv, err = sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func isSheetServiceGranted() error {
	_, err := srv.Spreadsheets.Values.Get(cfg.Cetele.SpreadSheetID, "Dün Özet").Do()
	//auth yoksa
	if err != nil {
		return err
	}
	return nil
}
