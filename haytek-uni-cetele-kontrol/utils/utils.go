package utils

import (
	"errors"
	"strings"
	"time"
)

func GetTarih(t time.Time) string {
	return t.Format("02.01.2006")
}
func GetSaat(t time.Time) string {
	return t.Format("15:04:05")
}

func ParseTarihFromCommandArguments(argumentstr string) (string, error) {
	//arguments := strings.Split(update.Message.CommandArguments(), " ")
	var err error
	if argumentstr == "" {
		err = errors.New("komuttan sonra bir tarih belirtin (örn: /gunlukozet 30.10.2022)")
		return "", err
	}
	arguments := strings.Split(strings.TrimSpace(argumentstr), " ")
	if arguments[0] == "" || len(arguments) > 1 {
		err = errors.New("komuttan sonra bir tarih belirtin (örn: /gunlukozet 30.10.2022)")
		return "", err
	}
	parts := strings.Split(strings.TrimSpace(arguments[0]), ".")
	if len(parts) != 3 {
		err = errors.New("tarih formatı hatalı (örn: /gunlukozet 30.10.2022)")
		return "", err
	}

	tarih := parts[2] + "-" + parts[1] + "-" + parts[0]
	return tarih, nil
}
