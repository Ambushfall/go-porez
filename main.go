package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const debug bool = true

var jsonRes UserInfo
var email string
var password string
var hasGotData bool
var upitData Upit
var upitStanja UpitStanja
var serve bool
var c *exec.Cmd

const requestURL string = "https://lpa.gov.rs/upitstanja/upit"

const (
	Large  int = 1000
	Medium int = 750
	Small  int = 500
)

var qrCodeURL string = fmt.Sprintf("https://nbs.rs/QRcode/api/qr/v1/gen/%d", Large)

// Parse flags
func init() {
	// err := os.Remove("out.png")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	const (
		defaultEmail     = ""
		usageEmail       = "Euprava Email"
		defaultNameEmail = "email"
		shortNameEmail   = "e"

		defaultPass     = ""
		usagePass       = "Euprava password"
		defaultNamePass = "password"
		shortNamePass   = "p"

		defaultServe     = false
		usageServe       = "Set whether to run server"
		defaultNameServe = "serve"
		shortNameServe   = "s"
	)
	flag.StringVar(&email, defaultNameEmail, defaultEmail, usageEmail)
	flag.StringVar(&email, shortNameEmail, defaultEmail, usageEmail+" (shorthand)")

	flag.StringVar(&password, defaultNamePass, defaultPass, usagePass)
	flag.StringVar(&password, shortNamePass, defaultPass, usagePass+" (shorthand)")

	flag.BoolVar(&serve, defaultNameServe, defaultServe, usageServe)
	flag.BoolVar(&serve, shortNameServe, defaultServe, usageServe+" (shorthand)")
}

func main() {

	flag.Parse()
	ports := ":80"

	if flag.Parsed() {
		CLI()
		// Serve via server through CLI Flag
		if serve {
			router("/", HandleJSONRequestParams, "admin", "123", "browser")

			str := fmt.Sprintf("Server is running at %q", ports)
			fmt.Println(str)
			log.Fatal(http.ListenAndServe(ports, nil))
		}
	}

	if len(email) == 0 || len(password) == 0 {
		CredentialsPrompt(email, password)
	}

	switch runtime.GOOS {
	case "windows":
		c = exec.Command("cmd", "/C", "start", "out.png")

	case "darwin":
		c = exec.Command("open", "out.png")

	default: //Mac & Linux
		c = exec.Command("echo", "Platform Unsupported")
		fmt.Println("Platform Unsupported")
		os.Exit(1)
	}

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

}

// check if args were set and use them
func CLI() {
	if len(email) > 0 && len(password) > 0 {

		fmt.Printf("email: %s password: %s", email, password)
		ok, res, err := try_robo_Login(email, password)
		if err != nil {
			log.Fatalf("Error: %v", err)
		} else {
			if ok {
				saveFile(res)
				fmt.Println("Success")
			}
		}
		// os.Exit(0)
	}
}

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	done <- true
}

func saveFile(response string) {
	if !serve {
		ers := os.WriteFile("out.png", []byte(response), 0644)
		if ers != nil {
			log.Fatalf("error: %#v", ers)

		} else {
			fmt.Println("gj")
		}
	}
}

func CredentialsPrompt(email string, password string) {
	fmt.Println("Euprava Email: ")
	fmt.Scanln(&email)
	fmt.Println("Euprava password: ")
	fmt.Scanln(&password)
	ok, res, err := try_robo_Login(email, password)
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		if ok {
			fmt.Println(email, password)
			fmt.Println("Success")
			saveFile(res)
		}
	}
}
