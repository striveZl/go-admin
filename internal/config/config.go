package config

import (
	"encoding/json"
	"fmt"
	"go-admin/pkg/logging"
)

type Config struct {
	Logger     logging.LoggerConfig
	General    General
	Util       Util
	Middleware Middleware
}

type General struct {
	AppName            string `default:"goadmin"`
	Version            string `default:"1.0.0"`
	WorkDir            string
	DisablePrintConfig bool
	Debug              bool
	HTTP               struct {
		Addr            string `default:":8040"`
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
