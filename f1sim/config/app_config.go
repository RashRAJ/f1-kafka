package config

type AppConfig struct {
	DriverNumber string `yaml:"driver,omitempty"`
	YearsStr     string `yaml:"years"`
	AllData      bool   `yaml:"all_data"`
	BaseURL      string `yaml:"base_url"`
}
