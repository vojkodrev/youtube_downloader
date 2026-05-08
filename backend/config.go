package main

type Config struct {
	StreamsDir      string
	SplitDuration   float64  `yaml:"split_duration"`
	RecodeQualities []string `yaml:"recode_qualities"`
	AllowedEmails   []string `yaml:"allowed_emails"`

	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
	BackendURL         string
	FrontendURL        string
	PublicHosts        []string
}
