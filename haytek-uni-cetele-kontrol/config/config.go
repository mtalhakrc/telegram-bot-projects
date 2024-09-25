package config

import "os"

var config *Config

func Get() *Config {
	return config
}

type Config struct {
	Cetele      CeteleAppConfig
	Credentials CredentialsConfig
	Database    DbConfig
}

type CeteleAppConfig struct {
	SpreadSheetID string
	AllowedUsers  map[string]bool
	Kisiler       map[string]int64
	Gruplar       map[string]int64
	RaporlarPath  string
}

type CredentialsConfig struct {
	//TokenPath        string
	//RefreshTokenPath string
	CredentialsPath string
}

type DbConfig struct {
	Name        string
	Username    string
	Password    string
	Host        string
	Port        string
	Socket      string
	Debug       bool
	MaxPoolSize int
	MaxIdleConn int
	MaxLifetime int
}

func SetupConfig() {
	config = &Config{}
	config.Cetele = CeteleAppConfig{
		SpreadSheetID: "15sToZcfyEp95WINbv1nuD_sTtTZxn1RmhgkBrlLIw9g",
		Kisiler: map[string]int64{
			"Talha Karaca":      952363491,
			"Ömer Faruk Gezer":  5669602367,
			"Eyüp Erbeyin":      0,
			"Enes Yılmaz":       1404072287,
			"Hamza Uysal":       1091982292,
			"Abdüssamet Çelik":  933834129,
			"Emin Güner":        0,
			"Mirza Şimşek":      1833713247,
			"Muzaffer":          0,
			"Abdüssamet Alişar": 1642068341,
		},

		AllowedUsers: map[string]bool{
			"mtalhakrc":   true,
			"hhuseyinpay": true,
		},

		Gruplar: map[string]int64{
			"mtalhakrc": 952363491,
		},

		RaporlarPath: "raporlar/",
	}

	config.Credentials = CredentialsConfig{
		//TokenPath:        "credentials/token.json",
		//RefreshTokenPath: "credentials/refreshToken.json",
		CredentialsPath: "credentials/fluted-ranger-364116-ea4e986f9ca1.json",
	}

	config.Database = DbConfig{
		Name:        "postgres",
		Username:    "postgres",
		Password:    "postgres",
		Host:        "localhost",
		Port:        "5430",
		Socket:      "",
		Debug:       false,
		MaxPoolSize: 5,
		MaxIdleConn: 1,
		MaxLifetime: 1800,
	}
	if os.Getenv("IS_DEVELOPMENT") == "true" {
		//config.Credentials.TokenPath = "/home/ubuntu/credentials/token.json"
		//config.Credentials.RefreshTokenPath = "/home/ubuntu/credentials/refreshToken.json"
		config.Credentials.CredentialsPath = "/home/ubuntu/credentials/fluted-ranger-364116-ea4e986f9ca1.json"
		config.Cetele.RaporlarPath = "/home/ubuntu/raporlar"
		config.Database.Port = "5432"
	}
}
