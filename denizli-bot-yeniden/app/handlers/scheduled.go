package handlers

import (
	"fmt"
	"github.com/go-bot-template/pkg/app"
	"github.com/go-bot-template/pkg/utils"
	"google.golang.org/api/sheets/v4"
	"time"
)

type ScheduledHandler struct {
	srv *sheets.Service
}

func NewScheduled(s *sheets.Service) ScheduledHandler {
	return ScheduledHandler{
		srv: s,
	}
}

type GunlukRapor struct {
	Tarih          time.Time     `json:"tarih"`
	MedreseOzetler []MedreseOzet `json:"medrese_ozetler"`
}
type MedreseOzet struct {
	Isim    string     `json:"isim"`
	Kisiler []KisiOzet `json:"kisiler"`
}
type KisiOzet struct {
	Isim         string `json:"isim"`
	DoldurulduMu bool   `json:"dolduruldu_mu"`
}

func (s ScheduledHandler) GunSonuOzetMesaji() []app.ScheduledResponse {
	rapor := s.GunlukRaporHazirla()
	str := GunlukOzetFormat(rapor)
	res := []app.ScheduledResponse{
		{
			UserID: 243296460,
			Result: str,
			Error:  nil,
		},
		{
			UserID: 952363491, //952363491 ben
			Result: str,
			Error:  nil,
		},
	}
	return res
}

func (s ScheduledHandler) GunlukRaporHazirla() GunlukRapor {
	var rapor GunlukRapor
	medreseSpreadsheetIDs := map[string]string{
		"HAYRAT":      "17zep31hMu1r8jH4tc2yFJsK_YA1eL5Dave5DWux4wes",
		"SUFFE":       "1kuu5wfNidRGaGam9nodA574hWV3T4K8jf-UUsE41X38",
		"TEFANİ":      "1XCDq8A_vOmqVhvqp77y5ElB31VNydNdDuWtuuGnTymw",
		"HAFIZ ALİ":   "1djTVXtvbvZSNH1XxYtfV9Py49VJ-KVJ9AVnfd-rq8Oc",
		"BAĞBAŞI":     "1xo6YYqWscQLOcLTaRl5MqAQjVYprSDCCgBU11OAa3BA",
		"SELİMİYE":    "1umJumhPNRqI7wHl5C4j2m8Z3hUZO9ZpqFFTJLdRhZak",
		"AYASOFYA":    "1JTKdkBAhnMDzTC-NEwaDgOgcBR2bwmS2F9ZpK0uvYYc",
		"RAVZA":       "13NaHw35uJZNCs02TpCit_AyCGDa62pmcr9upKKLA-F0",
		"HASAN FEYZİ": "1gU2kQ8ytM00jRV0X2PKmadFp8vY4dPWK8KJdhVGX7Yo",
	}

	for isim, id := range medreseSpreadsheetIDs {
		medrese, err := s.GetMedreseProgramOzetFromSheet(isim, id)
		if err != nil {
			panic(err)
		}
		rapor.MedreseOzetler = append(rapor.MedreseOzetler, medrese)
	}
	rapor.Tarih = time.Now()
	return rapor
}

func (s ScheduledHandler) GetMedreseProgramOzetFromSheet(medreseisim string, spreadsheetid string) (MedreseOzet, error) {
	var medrese MedreseOzet
	resp, err := s.srv.Spreadsheets.Values.Get(spreadsheetid, "ÖZET").Do()
	if err != nil {
		return medrese, err
	}
	var kisilerOzet []KisiOzet

	var kisiOzet KisiOzet
	for i := 1; i < len(resp.Values); i++ {
		isim := resp.Values[i][0].(string)
		kisiOzet.Isim = isim
		kisiOzet.DoldurulduMu = !(resp.Values[i][1].(string) == "#N/A")
		kisilerOzet = append(kisilerOzet, kisiOzet)

	}

	medrese.Isim = medreseisim
	medrese.Kisiler = kisilerOzet

	return medrese, nil
}

func GunlukOzetFormat(rapor GunlukRapor) string {
	var str string
	var eksikMi = false
	loc, _ := time.LoadLocation("Europe/Istanbul")
	str += fmt.Sprintf("Denizli Talebeleri %s Tarihli Gun Sonu Ozeti\n", utils.GetTarih(rapor.Tarih))
	str += fmt.Sprintf("Saat: %s \n", utils.GetSaat(rapor.Tarih.In(loc)))
	str += "------------------------- \n"

	for _, medrese := range rapor.MedreseOzetler {
		str += fmt.Sprintf("%s Medresesi:\n", medrese.Isim)
		for _, kisi := range medrese.Kisiler {
			if !kisi.DoldurulduMu {
				eksikMi = true
				str += fmt.Sprintf("%s\n", kisi.Isim)
			}
		}
		if !eksikMi {
			str += fmt.Sprintf("Bu medrese için eksiklik yok! 🎉🎉🎉\n")
		}
		str += "------------------------- \n"
		eksikMi = false
	}
	return str
}
