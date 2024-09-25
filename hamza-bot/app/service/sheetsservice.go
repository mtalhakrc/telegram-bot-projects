package service

import (
	"context"
	_ "embed"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
)

//go:embed fluted-ranger-364116-2db60b573b1a.json
var b []byte
var srv *sheets.Service

func InitSheetsService() {

	var err error
	srv, err = sheets.NewService(context.Background(), option.WithCredentialsJSON(b))
	if err != nil {
		panic(err)
	}
	log.Println("sheet service initialized")
}
func Get() *sheets.Service {
	return srv
}

type ISheetsService interface {
	GetFromSheet(range_ string) ([][]interface{}, error)
	UpdateSheet(range_ string, values [][]interface{}) (*sheets.UpdateValuesResponse, error)
	DeleteFromSheet(range_ string) (*sheets.ClearValuesResponse, error)

	//TestSheetExist kaydolan kullanıcı sheet servicesi için sheets servicesinsdeki kendi adını da girmeli. bunun için gireceği adda bir sheet sayfası olmalı.
	TestSheetExist(name string) bool
}

type ExSheetsService struct {
	spreadsheetID string
	service       *sheets.Service
}

func (s ExSheetsService) GetFromSheet(range_ string) ([][]interface{}, error) {
	res, err := s.service.Spreadsheets.Values.Get(s.spreadsheetID, range_).Do()
	if err != nil {
		return nil, err
	}
	return res.Values, err
}

func (s ExSheetsService) UpdateSheet(range_ string, values [][]interface{}) (*sheets.UpdateValuesResponse, error) {
	m, err := s.service.Spreadsheets.Values.Update(s.spreadsheetID, range_, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(context.Background()).Do()
	return m, err
}
func (s ExSheetsService) DeleteFromSheet(range_ string) (*sheets.ClearValuesResponse, error) {
	m, err := s.service.Spreadsheets.Values.Clear(s.spreadsheetID, range_, &sheets.ClearValuesRequest{}).Context(context.Background()).Do()
	return m, err
}
func (s ExSheetsService) TestSheetExist(name string) bool {
	_, err := s.GetFromSheet(fmt.Sprintf("%s!%s", name, "A1"))
	return err == nil
}
func NewSheetsService(s *sheets.Service, spreadsheetID string) ISheetsService {
	return ExSheetsService{
		service:       s,
		spreadsheetID: spreadsheetID,
	}
}
