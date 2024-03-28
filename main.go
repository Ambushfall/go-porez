package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
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

// type UserInfoResponse struct

type UserInfo struct {
	HTTPSchemaIDRsClaimsMail       string `json:"http://schema.id.rs/claims/mail"`
	HTTPSchemaIDRsClaimsCountry    string `json:"http://schema.id.rs/claims/country"`
	Sub                            string `json:"sub"`
	HTTPSchemaIDRsClaimsAal        string `json:"http://schema.id.rs/claims/aal"`
	HTTPSchemaIDRsClaimsGivenname  string `json:"http://schema.id.rs/claims/givenname"`
	HTTPSchemaIDRsClaimsUmcn       string `json:"http://schema.id.rs/claims/umcn"`
	HTTPSchemaIDRsClaimsFamilyname string `json:"http://schema.id.rs/claims/familyname"`
	HTTPSchemaIDRsClaimsCity       string `json:"http://schema.id.rs/claims/city"`
	HTTPSchemaIDRsClaimsIal        string `json:"http://schema.id.rs/claims/ial"`
}

type UpitStanja struct {
	Pib                     string                    `json:"pib"`
	DatumZaduzenjaDo        string                    `json:"datumZaduzenjaDo"`
	DatumUplateDo           string                    `json:"datumUplateDo"`
	IsObveznik              bool                      `json:"isObveznik"`
	UpitStanjaSaldoOpstList []UpitStanjaSaldoOpstList `json:"upitStanjaSaldoOpstList"`
}
type ListaPromena struct {
	KnjPromSifra     string  `json:"knjPromSifra"`
	KnjPromSifraIP   string  `json:"knjPromSifraIp"`
	KnjPromSifraZp   string  `json:"knjPromSifraZp"`
	KnjPromDISSifra  string  `json:"knjPromDISSifra"`
	KnjPromDISOpis   string  `json:"knjPromDISOpis"`
	BrojNaloga       string  `json:"brojNaloga"`
	DisDokument      any     `json:"disDokument"`
	Datum            string  `json:"datum"`
	PrometDuguje     float64 `json:"prometDuguje"`
	PrometPotrazuje  float64 `json:"prometPotrazuje"`
	KamataZaduzenje  float64 `json:"kamataZaduzenje"`
	KamataObracunata float64 `json:"kamataObracunata"`
	KamataNaplacena  float64 `json:"kamataNaplacena"`
	SaldoGlavnica    float64 `json:"saldoGlavnica"`
	SaldoKamata      float64 `json:"saldoKamata"`
	KnjPromOpis      string  `json:"knjPromOpis"`
	KnjPromOpisPu    string  `json:"knjPromOpisPu"`
}
type UpitStanjaSaldoList struct {
	Racun            string         `json:"racun"`
	RacunCeo         string         `json:"racunCeo"`
	RacunOpis        string         `json:"racunOpis"`
	SaldoDuguje      float64        `json:"saldoDuguje"`
	SaldoPotrazuje   float64        `json:"saldoPotrazuje"`
	KamataZaduzenje  float64        `json:"kamataZaduzenje"`
	KamataObracunata float64        `json:"kamataObracunata"`
	KamataNaplacena  float64        `json:"kamataNaplacena"`
	SaldoGlavnica    float64        `json:"saldoGlavnica"`
	SaldoKamata      float64        `json:"saldoKamata"`
	SaldoUkupan      float64        `json:"saldoUkupan"`
	ListaPromena     []ListaPromena `json:"listaPromena"`
}
type UpitStanjaSaldoOpstList struct {
	UpitStanjaSaldoList []UpitStanjaSaldoList `json:"upitStanjaSaldoList"`
	SifraOpstine        string                `json:"sifraOpstine"`
	NazivOpstine        string                `json:"nazivOpstine"`
	PozivNaBroj         string                `json:"pozivNaBroj"`
	ObveznikIdent       string                `json:"obveznikIdent"`
	DatumUpita          string                `json:"datumUpita"`
	VremeObrade         string                `json:"vremeObrade"`
}

type Upit struct {
	DatumZaduzenjaDo any    `json:"datumZaduzenjaDo"`
	DatumUplateDo    any    `json:"datumUplateDo"`
	Pib              string `json:"pib"`
	Racun            any    `json:"racun"`
	Detail           any    `json:"detail"`
}

var jsonRes UserInfo
var email string
var password string
var hasGotData bool
var upitData Upit
var upitStanja UpitStanja

const requestURL string = "https://lpa.gov.rs/upitstanja/upit"

// Example:
// parseJSON(jsonstring, &v)
func parseJSON(jsonstring string, v any) {
	resBytes := []byte(jsonstring)

	err := json.Unmarshal(resBytes, &v)
	if err != nil {
		fmt.Println(err)
	}
}

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

func JSONPostHeader(url string, body any, header map[string]string) (string, error) {
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		return "", nil
	}

	reader := bytes.NewReader(bodyBytes)

	// Make HTTP POST request
	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return "", nil
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		return "", nil
	}

	// Close response body
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		return string(responseBody), errors.New("400/500 status code error")
	}

	return string(responseBody), nil
}

func try_robo_Login() {

	browser := startNativeOrPORT()
	// startMonitor(browser)
	defer browser.MustClose()
	router := browser.HijackRequests()
	defer router.MustStop()

	router.MustAdd("https://prijava.eid.gov.rs/oauth2/userinfo", RoboHandlerResponse(func(ctx *rod.Hijack) {
		fmt.Println(upitStanja)
	}))

	go router.Run()

	page :=
		browser.MustPage("https://lpa.gov.rs/jisportal/homepage") // .MustWait(`() => document.title === 'hi'`)

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
