package settings

import (
	"os"
	"time"
	"github.com/astaxie/beego/config"
	"fmt"
	"strconv"
	"io/ioutil"
	"path/filepath"
	"log"
)

const (
	CONF_NAME string = "/conf/node.conf"
	permMode os.FileMode = 0666
)

// Singleton
var instantiated *Settings = nil
var iter uint8

func SettingsPtr() *Settings {
	if instantiated == nil {
		instantiated = new(Settings)
	}
	return instantiated
}

type Settings struct {
	IP		string
	Port		int
	DeviceList	[]string
	StartTime	time.Time
	UpTime		time.Duration
	Baud		int
	Timeout		time.Duration
	Cachetime	time.Duration
	StopBits	int
	cfg 		config.IniConfigContainer
	dir		string
}

func (s *Settings) Init() *Settings {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	s.StartTime = time.Now()
	s.dir = fmt.Sprintf("%s/", dir)
	s.IP = "127.0.0.1"
	s.Port = 3001
	s.Baud = 19200
	s.Timeout = 2
	s.Cachetime = 3600
	s.StopBits = 2

	s.Load()

	return s
}

func (s *Settings) Load() (*Settings, error) {

	iter++

	if _, err := os.Stat(s.dir + CONF_NAME); os.IsNotExist(err) {
		return s.Save()
	}

	// read config file
	cfg, err := config.NewConfig("ini", s.dir + CONF_NAME)
	if err != nil {
		return s, err
	}

	s.IP = cfg.String("ip")
	s.Port, _ = cfg.Int("port")
	s.Baud, _ = cfg.Int("baud")
	timeout, _ := cfg.Int("timeout")
	s.Timeout = time.Duration(timeout)
	cachetime, _ := cfg.Int("cachetime")
	s.Cachetime = time.Duration(cachetime)
	s.StopBits, _ = cfg.Int("stopbits")

	return s, err
}

func (s *Settings) Save() (*Settings, error) {

	if _, err := os.Stat(s.dir + CONF_NAME); os.IsNotExist(err) {
		ioutil.WriteFile(s.dir + CONF_NAME, []byte{}, permMode)
	}

	cfg, err := config.NewConfig("ini", s.dir + CONF_NAME)
	if err != nil {
		return s, err
	}

	cfg.Set("app_version", s.AppVresion())
	cfg.Set("ip", s.IP)
	cfg.Set("port", strconv.Itoa(s.Port))
	cfg.Set("baud", strconv.Itoa(s.Baud))
	cfg.Set("timeout", strconv.Itoa(int(s.Timeout)))
	cfg.Set("cachetime", strconv.Itoa(int(s.Cachetime)))
	cfg.Set("stopbits", strconv.Itoa(s.StopBits))

	if err := cfg.SaveConfigFile(s.dir + CONF_NAME); err != nil {
		fmt.Printf("err with create conf file: %s\n", s.dir + CONF_NAME)
		return s, err
	}

	return s, nil
}

func (s *Settings) AppVresion() string {
	return os.Getenv("VERSION")
}

func init() {
	instantiated = new(Settings)
}