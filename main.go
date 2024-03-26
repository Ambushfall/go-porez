package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/go-rod/rod"
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

func Example_hijack_requests() {
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
	ports := ":80"

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

	str := fmt.Sprintf("Server is running at %q", ports)
	fmt.Println(str)
	log.Fatal(http.ListenAndServe(ports, nil))
}
