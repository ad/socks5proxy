package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/ad/socks5proxy/config"
	"github.com/armon/go-socks5"
	"github.com/getlantern/systray"
	"github.com/kardianos/osext"
)

var (
	version = `0.0.2`
	ctx     context.Context
	cancel  context.CancelFunc

	cfg *config.Config
	err error

	listener net.Listener
	server   *socks5.Server
)

var banner = `
                 __              
  _ _  _  /_  _ /_   _  _ _      
_\ /_//_ /\ _\ ._/  /_// /_/></_/
                   /          _/                             
`

func main() {
	ctx, cancel = context.WithCancel(context.Background())

	log.Println(banner)

	defer func() {
		cancel()
		systray.Quit()
	}()

	// Init config
	cfg, err = config.GetConfig()
	if err != nil {
		log.Fatal(err)
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

	server, err = socks5.New(socks5config)
	if err != nil {
		panic(err)
	}

	log.Println("Started")

	listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.ProxyPort))
	if err != nil {
		panic(err)
	}

	log.Printf("Running on port %d\n", cfg.ProxyPort)
	if cfg.ProxyUser != "" {
		log.Printf("User: %s, Password: %s\n", cfg.ProxyUser, cfg.ProxyPassword)
	}

	if !cfg.HideSystrayIcon {
		go func() {
			if err = server.Serve(listener); err != nil {
				log.Println(err)
			}
		}()

		systray.Run(onReady, onExit)
	} else {
		if err = server.Serve(listener); err != nil {
			log.Println(err)
		}
	}
}

func onReady() {
	systray.SetTitle("🧦")
	systray.SetTooltip("Socks5 Proxy")
	mTitle := systray.AddMenuItem(fmt.Sprintf("Socks5Proxy v%s", version), "App title")
	mTitle.Disable()
	mPort := systray.AddMenuItem(fmt.Sprintf("Running on port %d", cfg.ProxyPort), "Proxy port")
	mPort.Disable()
	if cfg.ProxyUser != "" {
		mCreds := systray.AddMenuItem(fmt.Sprintf("User: %s, Password: %s", cfg.ProxyUser, cfg.ProxyPassword), "Proxy credentials")
		mCreds.Disable()
	}
	started := true
	mStop := systray.AddMenuItem("Stop", "Stop proxy")
	mStart := systray.AddMenuItem("Start", "Start proxy")
	mStart.Hide()

	toggle := func() {
		if started {
			mStart.Hide()
			mStop.Hide()
			mPort.Hide()

			err = listener.Close()
			if err != nil {
				panic(err)
			}

			mStart.Show()
		} else {
			mStop.Hide()
			mStart.Hide()

			listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.ProxyPort))
			if err != nil {
				log.Println(err)
			}

			log.Printf("Running on port %d\n", cfg.ProxyPort)
			if cfg.ProxyUser != "" {
				log.Printf("User: %s, Password: %s\n", cfg.ProxyUser, cfg.ProxyPassword)
			}

			go func() {
				if err = server.Serve(listener); err != nil {
					log.Println(err)
				}
			}()

			mStop.Show()
			mPort.Show()
		}
		started = !started
	}

	mRestart := systray.AddMenuItem("Restart", "Restart app")
	mQuit := systray.AddMenuItem("Quit", "Quit app")

	for {
		select {
		case <-mRestart.ClickedCh:
			fmt.Println("Requesting restart")
			cancel()
			Restart()
			return
		case <-mQuit.ClickedCh:
			fmt.Println("Requesting quit")
			cancel()
			systray.Quit()
			return
		case <-mStart.ClickedCh:
			toggle()
		case <-mStop.ClickedCh:
			toggle()
		}
	}
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
