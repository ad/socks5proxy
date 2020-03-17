package config

import (
	"flag"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v2"
)

// Config struct
type Config struct {
	ProxyUser     string
	ProxyPassword string
	ProxyPort     int

	HideSystrayIcon bool
}

var cfg = Config{}

// StringSlice is a flag.Value that collects each Set string
// into a slice, allowing for repeated flags.
type StringSlice []string

// Set implements flag.Value and appends the string to the slice.
func (ss *StringSlice) Set(s string) error {
	(*ss) = append(*ss, s)
	return nil
}

// String implements flag.Value and returns the list of
// strings, or "..." if no strings have been added.
func (ss *StringSlice) String() string {
	if len(*ss) <= 0 {
		return "..."
	}
	return strings.Join(*ss, ", ")
}

// GetConfig Gets the conf in the config file
func GetConfig() (*Config, error) {
	if err := readConfig(&cfg); err == nil {
		return &cfg, err
	}

	return &cfg, nil
}

// readConfig Reads the conf in the config file
func readConfig(cfg *Config) error {
	fs := flag.NewFlagSet("socks5proxy", flag.ContinueOnError)

	_ = fs.String("config", "config/config.json", "load config file from `path` (optional)")

	fs.StringVar(&cfg.ProxyUser, "proxy_user", "", "Socks5 proxy username")
	fs.StringVar(&cfg.ProxyPassword, "proxy_password", "", "Socks5 proxy password")
	fs.IntVar(&cfg.ProxyPort, "proxy_port", 1080, "Socks5 proxy port")
	fs.BoolVar(&cfg.HideSystrayIcon, "hide_systray_icon", false, "show or hide systray icon")

	err := ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
		ff.WithEnvVarPrefix("SOCKS5PROXY"),
	)

	return err
}
