package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	BaseURL  string `env:"BASE_URL,required"`
	DataPath string `env:"DATA_PATH,required"`

	IncidentToken string `env:"INCIDENT_TOKEN,required"`
	IncidentRottenDuration time.Duration `env:"INCIDENT_ROTTEN_DURATION" envDefault:"4h"`

	BasicAuthUser string `env:"BASIC_AUTH_USER,required"`
	BasicAuthPass string `env:"BASIC_AUTH_PASS,required"`


	TwilSID   string `env:"TWIL_SID,required"`
	TwilToken string `env:"TWIL_TOKEN,required"`
	TwilPhone string `env:"TWIL_PHONE,required"`
}

var once sync.Once
var config Config

func Get() Config {
	once.Do(func() {
		err := env.Parse(&config)

		if err != nil {
			if strings.HasPrefix(err.Error(), "env: could not load content of file \"\"") {
				fmt.Printf("%+v\n\n", config)
			}
			os.Exit(1)
		}
	})

	return config
}
