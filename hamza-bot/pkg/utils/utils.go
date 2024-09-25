package utils

import (
	"fmt"
	"strings"
	"time"
)

var loc *time.Location
var err error

func init() {
	loc, err = time.LoadLocation("Europe/Istanbul")
	if err != nil {
		panic(err)
	}
}

func GetTarih(t time.Time) string {
	return t.Format("02.01.2006")
}
func GetSaat(t time.Time) string {
	return t.Format("15:04:05")
}

func ParseCommandArguments(params string) []string {
	if params == "" {
		return nil
	}
	return strings.Split(params, " ")
}

func GetNow() time.Time {
	return time.Now().In(loc)
}

func ParseStrTime(str string) time.Time {
	s := fmt.Sprintf("%s %s", time.Now().In(loc).Format("02-Jan-2006"), str)
	t1, err := time.ParseInLocation("02-Jan-2006 15:04:05", s, loc)
	if err != nil {
		panic(err)
	}
	return t1
}
