package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_Host     string
	DB_Port     int
	DB_User     string
	DB_Name     string
	DB_Password string

	JWT_Secret string
	SMTP       SMTPConfig

	GoogleCredentials string
	GCP_ProjectID     string
	GCP_BucketName    string

	Midtrans MidtransConfig
}

type SMTPConfig struct {
	SMTPHOST string
	SMTPPORT string
	SMTPUSER string
	SMTPPASS string
}

type MidtransConfig struct {
	ServerKey string
	ClientKey string
}

func InitConfig() *Config {
	_ = godotenv.Load(".env")
	var res = new(Config)

	res.DB_Host = os.Getenv("DB_HOST")
	res.DB_Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	res.DB_User = os.Getenv("DB_USER")
	res.DB_Password = os.Getenv("DB_PASSWORD")
	res.DB_Name = os.Getenv("DB_NAME")

	res.JWT_Secret = os.Getenv("JWT_SECRET")

	res.SMTP.SMTPHOST = os.Getenv("SMTP_HOST")
	res.SMTP.SMTPPORT = os.Getenv("SMTP_PORT")
	res.SMTP.SMTPUSER = os.Getenv("SMTP_USER")
	res.SMTP.SMTPPASS = os.Getenv("SMTP_PASSWORD")

	// Google Cloud Config
	res.GoogleCredentials = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	res.GCP_ProjectID = os.Getenv("GCP_PROJECT_ID")
	res.GCP_BucketName = os.Getenv("GCP_BUCKET_NAME")

	// Midtrans
	res.Midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	res.Midtrans.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")

	return res
}
