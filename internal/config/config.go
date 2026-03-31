package config

import (
	"encoding/json"
	"fmt"
	"go-admin/pkg/logging"
)

type Config struct {
	Logger     logging.LoggerConfig
	General    General
	Database   Database
	Util       Util
	Middleware Middleware
}


type Database struct {
	Driver          string `default:"postgres"`
	Host            string `default:"127.0.0.1"`
	Port            int    `default:"5432"`
	User            string `default:"go_admin"`
	Password        string
	DBName          string
	SSLMode         string `default:"disable"`
	Timezone        string `default:"Asia/Shanghai"`
	MaxIdleConns    int    `default:"10"`
	MaxOpenConns    int    `default:"100"`
	ConnMaxLifetime int    `default:"3600"` // seconds
}

type General struct {
	AppName            string `default:"goadmin"`
	Version            string `default:"1.0.0"`
	WorkDir            string
	DisablePrintConfig bool
	Debug              bool
	HTTP               struct {
		Addr            string `default:":3000"`
		ShutdownTimeout int    `default:"10"` // seconds
		ReadTimeout     int    `default:"60"` // seconds
		WriteTimeout    int    `default:"60"` // seconds
		IdleTimeout     int    `default:"10"` // seconds
		CertFile        string
		KeyFile         string
	}
}

type Util struct {
	Captcha struct {
		Length    int    `default:"4"`
		Width     int    `default:"400"`
		Height    int    `default:"160"`
		CacheType string `default:"memory"`
	}
}

func (c *Config) IsDebug() bool {
	return c.General.Debug
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Failed to marshal config: " + err.Error())
	}
	return string(b)
}

func (c *Config) Print() {
	if c.General.DisablePrintConfig {
		return
	}
	fmt.Println("// ----------------------- Load configurations start ------------------------")
	fmt.Println(c.String())
	fmt.Println("// ----------------------- Load configurations end --------------------------")
}
