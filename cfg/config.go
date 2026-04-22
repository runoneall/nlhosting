package cfg

import (
	"encoding/json"
	"log"
	"os"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

var Config = struct {
	Serve string `json:"serve"`

	Cookie struct {
		Secret string `json:"secret"`
	} `json:"cookie"`

	APIToken struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"apitoken"`

	OAuth2 struct {
		Client struct {
			ID     string `json:"id"`
			Secret string `json:"secret"`
		} `json:"client"`

		Redirect string   `json:"redirect"`
		Scopes   []string `json:"scopes"`
	} `json:"oauth2"`

	Server struct {
		Available []string                `json:"available"`
		Config    map[string]ServerConfig `json:"config"`
		Limit     map[string]int          `json:"limit"`
	} `json:"server"`

	Mail ServerConfig `json:"mail"`
}{}

func init() {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&Config); err != nil {
		log.Fatalln(err)
	}
}
