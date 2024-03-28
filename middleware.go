package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func grabClient(sec_ch_ua string) string {

	reg, err := regexp.Compile(`\P{L}+`)
	if err != nil {
		fmt.Println(err)
	}
	replaced_sec := reg.ReplaceAllString(sec_ch_ua, " ")
	newStripped := strings.Split(replaced_sec, "v")

	grablast := newStripped[len(newStripped)-2]
	return grablast
}

/*
Attach an auth middleware to your http.HandleFunc as callback.
Include your own function inside of BasicAuth

	http.HandleFunc("/", BasicAuth(func(w http.ResponseWriter, r *http.Request) {

Code to be executed after authorizing

	....

	}, username, password, realm))
*/
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()
		client := grabClient(r.Header.Get("Sec-Ch-Ua"))

		unauthorized := !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1

		if unauthorized {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)

			errResp := ErrResponse{
				Resonse: "Client unsupported",
				Code:    401,
				Client:  strings.Trim(client, " "),
			}

			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				log.Println("failed", err)
			}

		} else {
			handler(w, r)
		}
	}
}

func router(route string, handler http.HandlerFunc, user string, pass string, realm string) {
	http.HandleFunc(route, BasicAuth(handler, user, pass, realm))
}
