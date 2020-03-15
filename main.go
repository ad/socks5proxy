package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ad/socks5proxy/config"
	"github.com/armon/go-socks5"
)

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

	config := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if cfg.ProxyUser != "" {
		creadentials := socks5.StaticCredentials{
			cfg.ProxyUser: cfg.ProxyPassword,
		}

		authenticator := socks5.UserPassAuthenticator{Credentials: creadentials}

		config = &socks5.Config{
			AuthMethods: []socks5.Authenticator{authenticator},
			Logger:      log.New(os.Stdout, "", log.LstdFlags),
		}
	}

	server, err := socks5.New(config)
	if err != nil {
		panic(err)
	}

	log.Println("Started")

	if err := server.ListenAndServe("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.ProxyPort)); err != nil {
		panic(err)
	}
}
