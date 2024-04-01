package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

func try_robo_Login() {

	browser := startNativeOrPORT()
	// startMonitor(browser)
	defer browser.MustClose()
	router := browser.HijackRequests()
	defer router.MustStop()

	router.MustAdd("https://prijava.eid.gov.rs/oauth2/userinfo", RoboHandlerResponse(func(ctx *rod.Hijack) {

		if len(upitStanja.UpitStanjaSaldoOpstList) >= 1 {
			qrbody := TransformerUpitaStanja(upitStanja)
			response := GenerateQR(qrbody, ctx.Request.Header("Authorization"))
			fmt.Println(response)
		} else {
			fmt.Println("Porez se ne moze proveriti trenutno, nema povratnih informacija")
		}

	}))

	go router.Run()

	page :=
		browser.MustPage("https://lpa.gov.rs/jisportal/homepage") // .MustWait(`() => document.title === "hi"`)

	fmt.Println("done")

	page.MustElement("div.sistem-login> a").MustClick()
	page.MustElement("#username1").MustClick().MustInput(email)
	page.MustElement("#password1").MustClick().MustInput(password)
	page.MustElement("#aetButtonUP1").MustClick()

	utils.Pause() // pause goroutine
}

func RoboHandlerResponse(handler func(ctx *rod.Hijack)) func(*rod.Hijack) {
	return func(ctx *rod.Hijack) {
		_ = ctx.LoadResponse(http.DefaultClient, true)
		p := ctx.Response.Body()

		if (len(p) > 5) && (!hasGotData) {
			if len(ctx.Request.Header("Authorization")) > 5 {
				hasGotData = true
			}

			parseJSON(p, &jsonRes)
			upitData.Pib = jsonRes.HTTPSchemaIDRsClaimsUmcn
			responseBody, err := JSONPostHeader(requestURL, &upitData, map[string]string{
				"Content-Type":  "application/json",
				"authorization": ctx.Request.Header("Authorization"),
			})
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(responseBody)
			parseJSON(responseBody, &upitStanja)
			handler(ctx)
		}
	}
}

func instantiateBrowser(u string) *rod.Browser {
	return rod.New().
		ControlURL(u).
		// SlowMotion(2 * time.Second).
		// Trace(true).
		MustConnect()

}

func startMonitor(browser *rod.Browser) {
	launcher.Open(browser.ServeMonitor(""))
}

func startNativeOrPORT() *rod.Browser {
	u, err := launcher.ResolveURL("")

	if err == nil {
		return instantiateBrowser(u)
	} else {
		path, exists := launcher.LookPath()

		if exists {
			u, err := launcher.New().
				Bin(path).
				Headless(false).
				Devtools(false).
				Launch()

			if err == nil {
				return instantiateBrowser(u)
			} else {
				return rod.New().MustConnect()
			}
		} else {
			return rod.New().MustConnect()
		}

	}
}
