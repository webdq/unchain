package global

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	SubAddresses   string `desc:"sub addresses" example:"node1.xxx.cn:80,node2.xxx.cn:443"`
	Port           string `desc:"port" def:"80"`
	RegisterUrl    string `desc:"register url" def:"https://admin.unchain.people.from.censorship"`
	RegisterToken  string `desc:"register token" def:"unchain people from censorship and surveillance"`
	AllowUsers     string `desc:"allow users" def:"" example:"903bcd04-79e7-429c-bf0c-0456c7de9cdc,903bcd04-79e7-429c-bf0c-0456c7de9cd1"`
	LogFile        string `desc:"log file path" def:""`
	DebugLevel     string `desc:"debug level" def:"DEBUG"`
	IntervalSecond string `desc:"interval second" def:"360"` //seconds
	GitHash        string `desc:"git hash" def:""`
	BuildTime      string `desc:"build time" def:""`
	RunAt          string `desc:"run at" def:""`
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}
func (c Config) PushIntervalSecond() int {
	iv, err := strconv.ParseInt(c.IntervalSecond, 10, 32)
	if err != nil {
		log.Println("failed to parse interval second:", err)
		return 360
	}
	return int(iv)
}

func loadEnv() *Config {
	opt := Config{}
	for i := 0; i < reflect.TypeOf(opt).NumField(); i++ {
		propertyName := reflect.TypeOf(opt).Field(i).Name
		key := snakeCaseUpper(propertyName)
		def := reflect.TypeOf(opt).Field(i).Tag.Get("def")
		desc := reflect.TypeOf(opt).Field(i).Tag.Get("desc")
		vv := osEnvWithDefault(key, desc, def)
		reflect.ValueOf(&opt).Elem().Field(i).SetString(vv)
	}
	return &opt
}

func snakeCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to _[a-z]
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

func snakeCaseUpper(camel string) string {
	return strings.ToUpper(snakeCase(camel))
}

func osEnvWithDefault(key, desc, def string) string {
	if v := os.Getenv(key); v == "" {
		fmt.Printf("%s <%s> 默认:  %s\n", key, desc, def)
		return def
	} else {
		return v
	}
}

func (c Config) ListenPort() int {
	iv, err := strconv.ParseInt(c.Port, 10, 32)
	if err != nil {
		log.Println("failed to parse port:", err)
		return 80
	}
	return int(iv)
}

var (
	gitHash   string
	buildTime string
)

var cfg *Config

func Cfg(tomlFilePath string) *Config {
	if cfg != nil {
		return cfg
	}
	cfgIns, err := loadFromToml(tomlFilePath)
	if err != nil {
		fmt.Println("unable to load config file form .toml file, use env instead")
		cfg = loadEnv()
	} else {
		cfg = cfgIns
	}
	cfg.GitHash = gitHash
	cfg.BuildTime = buildTime
	cfg.RunAt = time.Now().Format("2006-01-02 15:04:05")
	return cfg
}

func loadFromToml(file string) (*Config, error) {
	opt := Config{}
	_, err := toml.DecodeFile(file, &opt)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file:%s %w", file, err)
	}
	return &opt, nil
}

func (c Config) LogLevel() slog.Level {
	l := slog.LevelDebug
	switch strings.ToUpper(c.DebugLevel) {
	case "DEBUG":
		l = slog.LevelDebug
	case "INFO":
		l = slog.LevelInfo
	case "WARN":
		l = slog.LevelWarn
	case "ERROR":
		l = slog.LevelError
	default:
		l = slog.LevelError
	}
	return l
}
func (c Config) UserIDS() []string {
	parts := strings.Split(c.AllowUsers, ",")
	ids := make([]string, 0)
	for _, uid := range parts {
		uid = strings.TrimSpace(uid)
		if uid != "" {
			ids = append(ids, uid)
		}
	}
	return ids
}

func (c Config) PushInterval() time.Duration {
	if c.PushIntervalSecond() <= 0 {
		return time.Minute * 60
	}
	return time.Second * time.Duration(c.PushIntervalSecond())
}
