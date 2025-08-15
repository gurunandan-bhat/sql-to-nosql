package reldb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	defaultConfigFileName = ".mario.json"
)

type DbType string
type DbAuth struct {
	Driver   string `json:"driver,omitempty"`
	Host     string `json:"host,omitempty"`
	Name     string `json:"name,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
type Config struct {
	InProduction bool              `json:"inProduction,omitempty"`
	LayoutRoot   string            `json:"layoutRoot"`
	DynDocRoot   string            `json:"dynDocRoot"`
	EmailSender  string            `json:"emailSender"`
	Db           map[DbType]DbAuth `json:"db,omitempty"`
	Remote       struct {
		Host      string `json:"host,omitempty"`
		Port      int    `json:"port,omitempty"`
		User      string `json:"user,omitempty"`
		ImageRoot string `json:"remote_image_root,omitempty"`
		KeyPath   string `json:"key_path,omitempty"`
		HostKey   string `json:"host_key,omitempty"`
	} `json:"remote,omitempty"`
	Security struct {
		CSRFKey        string   `json:"csrfKey,omitempty"`
		AllowedOrigins []string `json:"allowedOrigins,omitempty"`
		JWTSecret      string   `json:"jwtSecret,omitempty"`
	} `json:"security,omitempty"`
	SMTP struct {
		Server    string `json:"server,omitempty"`
		Port      int    `json:"port,omitempty"`
		User      string `json:"user,omitempty"`
		Password  string `json:"password,omitempty"`
		TestEmail string `json:"test_email,omitempty"`
	} `json:"smtp,omitempty"`
}

var c *Config
var onceConfig sync.Once

func Configuration(configFileName ...string) (*Config, error) {

	onceConfig.Do(func() {

		var cfname string

		switch len(configFileName) {
		case 0:
			dirname, err := os.UserHomeDir()
			if err != nil {
				panic(fmt.Sprintf("Cannot get home dir: %s", err))
			}
			cfname = fmt.Sprintf("%s/%s", dirname, defaultConfigFileName)
		case 1:
			cfname = configFileName[0]
		default:
			panic("incorrect arguments for configuration file name")
		}

		configFile, err := os.Open(cfname)
		if err != nil {
			panic(fmt.Sprintf("failed to open config file %s: %s", cfname, err))
		}
		defer func() {
			if err := configFile.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(&c)
		if err != nil {
			panic(fmt.Sprintf("failed to decode configuration: %s", err))
		}
	})

	return c, nil

}
