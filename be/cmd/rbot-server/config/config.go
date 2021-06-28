package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/olvrng/rbot/be/pkg/xconfig"
)

type Config struct {
	HTTP       HTTP      `yaml:"http"`
	Messenger  Messenger `yaml:"messenger"`
	StaticPath string    `yaml:"static_path"`
}

type Messenger struct {
	VerifyToken     string `yaml:"verify_token"`
	PageAccessToken string `yaml:"page_access_token"`
}

func (m *Messenger) MustLoadEnv(prefix string) {
	xconfig.EnvMap{
		prefix + "_VERIFY_TOKEN":      &m.VerifyToken,
		prefix + "_PAGE_ACCESS_TOKEN": &m.PageAccessToken,
	}.MustLoad()
}

type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (h *HTTP) ListeningAddress() string {
	return fmt.Sprintf("%s:%v", h.Host, h.Port)
}

func (h *HTTP) MustLoadEnv(prefix string) {
	xconfig.EnvMap{
		prefix + "_HOST": &h.Host,
		prefix + "_PORT": &h.Port,
	}.MustLoad()
}

func Default() Config {
	cfg := Config{}
	cfg.HTTP.Port = 8000
	cfg.Messenger.VerifyToken = "randomToken"

	// expect to run in rbot/be
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg.StaticPath = filepath.Join(wd, "../apps/board/build")

	return cfg
}

func Load(filename string) (Config, error) {
	cfg := Default()
	if filename != "" {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return cfg, err
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return cfg, err
		}
	}
	cfg.HTTP.MustLoadEnv("HTTP")
	cfg.Messenger.MustLoadEnv("MESSENGER")
	return cfg, nil
}
