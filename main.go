package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

type Time struct {
	CurrentTime string `json:"current_time"`
}

type ErrResponse struct {
	Resonse string `json:"response"`
	Code    int    `json:"code"`
	Client  string `json:"client"`
}

func router(route string, handler http.HandlerFunc, user string, pass string, realm string) {
	http.HandleFunc(route, BasicAuth(handler, user, pass, realm))
}

var email string
var password string

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

func hijack_requests() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	router := browser.HijackRequests()
	defer router.MustStop()

	router.MustAdd("*.js", func(ctx *rod.Hijack) {
		// Here we update the request's header. Rod gives functionality to
		// change or update all parts of the request. Refer to the documentation
		// for more information.
		ctx.Request.Req().Header.Set("My-Header", "test")

		// LoadResponse runs the default request to the destination of the request.
		// Not calling this will require you to mock the entire response.
		// This can be done with the SetXxx (Status, Header, Body) functions on the
		// ctx.Response struct.
		_ = ctx.LoadResponse(http.DefaultClient, true)

		// Here we append some code to every js file.
		// The code will update the document title to "hi"
		ctx.Response.SetBody(ctx.Response.Body() + "\n document.title = 'hi' ")
	})

	go router.Run()

	browser.MustPage("https://go-rod.github.io").MustWait(`() => document.title === 'hi'`)

	fmt.Println("done")

	// Output: done
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

	router("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %q", html.EscapeString(r.URL.Path))

	}, "a", "123456", "index")

	router("/time", func(w http.ResponseWriter, r *http.Request) {
		currentTime := []Time{
			{
				CurrentTime: time.Now().Format(http.TimeFormat),
			},
		}
		if err := json.NewEncoder(w).Encode(currentTime); err != nil {
			log.Println("failed", err)
		}
	}, "admin", "123456", "time")

	router("/browser", func(w http.ResponseWriter, r *http.Request) {
		try_robo_Login()
	}, "admin", "123", "browser")

	str := fmt.Sprintf("Server is running at %q", ports)
	fmt.Println(str)
	log.Fatal(http.ListenAndServe(ports, nil))
}

func try_robo_Login() {
	// Headless runs the browser on foreground, you can also use flag "-rod=show"
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New().
		Headless(false).
		Devtools(false)

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		// Trace(true).
		// SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	page :=
		browser.MustPage("https://lpa.gov.rs/jisportal/homepage")

		// login :=
	page.MustElement("div.sistem-login> a").MustClick()
	page.MustElement("#username1").MustClick().MustInput(email)
	page.MustElement("#password1").MustClick().MustInput(password)
	page.MustElement("#aetButtonUP1").MustClick()
	// fmt.Println(err)
	// textarea.MustInput("git").MustType(input.Enter)

	// text := page.MustElement(".codesearch-results p").MustText()

	// fmt.Println(text)

	utils.Pause() // pause goroutine
}
