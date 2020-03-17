package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/ad/socks5proxy/config"
	"github.com/armon/go-socks5"
	"github.com/getlantern/systray"
	"github.com/kardianos/osext"
)

var version = `0.0.1`

var banner = `
                 __              
  _ _  _  /_  _ /_   _  _ _      
_\ /_//_ /\ _\ ._/  /_// /_/></_/
                   /          _/                             
`

func main() {
	log.Println(banner)

	// Init config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	if !cfg.HideSystrayIcon {
		systray.Run(onReady, onExit)
	}

	socks5config := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if cfg.ProxyUser != "" {
		creadentials := socks5.StaticCredentials{
			cfg.ProxyUser: cfg.ProxyPassword,
		}

		authenticator := socks5.UserPassAuthenticator{Credentials: creadentials}

		socks5config = &socks5.Config{
			AuthMethods: []socks5.Authenticator{authenticator},
			Logger:      log.New(os.Stdout, "", log.LstdFlags),
		}
	}

	server, err := socks5.New(socks5config)
	if err != nil {
		panic(err)
	}

	log.Println("Started")

	if err := server.ListenAndServe("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.ProxyPort)); err != nil {
		panic(err)
	}
}

func onReady() {
	systray.SetTitle("🧦")
	systray.SetTooltip("Socks5 Proxy")
	mTitle := systray.AddMenuItem("Socks5Proxy", "App title")
	mTitle.Disable()
	mVersion := systray.AddMenuItem(fmt.Sprintf("Version %s", version), "App version")
	mVersion.Disable()
	mRestart := systray.AddMenuItem("Restart", "Restart app")
	mQuit := systray.AddMenuItem("Quit", "Quit app")
	go func() {
		<-mRestart.ClickedCh
		fmt.Println("Requesting restart")
		Restart()
	}()
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}

// Restart app
func Restart() error {
	file, error := osext.Executable()
	if error != nil {
		return error
	}

	error = syscall.Exec(file, os.Args, os.Environ())
	if error != nil {
		return error
	}

	return nil
}
