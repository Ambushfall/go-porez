package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const debug bool = true

var jsonRes UserInfo
var email string
var password string
var hasGotData bool
var upitData Upit
var upitStanja UpitStanja

const requestURL string = "https://lpa.gov.rs/upitstanja/upit"

const (
	Large  int = 1000
	Medium int = 750
	Small  int = 500
)

var qrCodeURL string = fmt.Sprintf("https://nbs.rs/QRcode/api/qr/v1/gen/%d", Large)

// Parse flags
func init() {
	const (
		defaultEmail     = ""
		usageEmail       = "email address"
		defaultNameEmail = "email"
		shortNameEmail   = "e"

		defaultPass     = ""
		usagePass       = "password"
		defaultNamePass = "password"
		shortNamePass   = "p"
	)
	flag.StringVar(&email, defaultNameEmail, defaultEmail, usageEmail)
	flag.StringVar(&email, shortNameEmail, defaultEmail, usageEmail+" (shorthand)")
	flag.StringVar(&password, defaultNamePass, defaultPass, usagePass)
	flag.StringVar(&password, shortNamePass, defaultPass, usagePass+" (shorthand)")
}

func main() {

	flag.Parse()
	ports := ":80"

	if flag.Parsed() {
		if len(email) > 0 && len(password) > 0 {

			msg := fmt.Sprintf("email: %s password: %s", email, password)
			fmt.Println(msg)
			try_robo_Login()
			// os.Exit(0)
		}
	}

	router("/", HandleJSONRequestParams, "admin", "123", "browser")

	str := fmt.Sprintf("Server is running at %q", ports)
	fmt.Println(str)
	log.Fatal(http.ListenAndServe(ports, nil))
}
