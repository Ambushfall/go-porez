package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func HandleJSONRequestParams(w http.ResponseWriter, r *http.Request) {
	var p RouteParams

	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	ok, rez, err := try_robo_Login(p.Email, p.Password)
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		if ok {
			ImagePayloadResponse(w, rez)
		}
	}

}

// Takes a responseWriter and any Payload
// Tries to encode it as JSON and logs if failed
func JSONPayloadResponse(w http.ResponseWriter, p any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Println("failed", err)
	}
}

// server string array buffer as image response
// you don't need a fucking image writer
func ImagePayloadResponse(w http.ResponseWriter, p string) {
	buff := []byte(p)
	w.Header().Set("Content-Type", "image/png")
	_, err := w.Write(buff)
	if err != nil {
		log.Fatalf("Err: %#v", err)
	}
}

// Example:
// parseJSON(jsonstring, &v)
func parseJSON(jsonstring string, v any) {
	resBytes := []byte(jsonstring)

	err := json.Unmarshal(resBytes, &v)
	if err != nil {
		fmt.Println(err)
	}
}

// Unpack array as return values
// similar to JS const [1,2] = Array
// var 1 int
// var 2 int
// var 3 int
// unpack(array, &1, &2, &3)
func unpack(s []string, vars ...*string) {
	for i, str := range s {
		*vars[i] = str
	}
}

// actual post request with json body parsed
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

// Take query and return body for QR Code
func TransformerUpitaStanja(upit UpitStanja) QRBody {
	var prvi string
	var drugi string
	var treci string
	var amount string
	var saldoUkupan float64 = upit.UpitStanjaSaldoOpstList[0].UpitStanjaSaldoList[0].SaldoUkupan
	unpack(strings.Split(upit.UpitStanjaSaldoOpstList[0].UpitStanjaSaldoList[0].RacunCeo, "-"), &prvi, &drugi, &treci)

	if debug {
		amount = "RSD253,00"
	} else {
		amount = fmt.Sprintf("RSD%s", strings.ReplaceAll(fmt.Sprintf("%.2f", saldoUkupan), ".", ","))
		if saldoUkupan < 0 {
			log.Fatalf("Vrednost poreza je ispod 0, %f", saldoUkupan)
		}
	}

	return QRBody{
		K:  "PR",
		V:  "01",
		C:  "1",
		R:  fmt.Sprintf("%s%s%s", prvi, fmt.Sprintf("%013s", drugi), treci),
		N:  fmt.Sprintf("LPA %s", upit.UpitStanjaSaldoOpstList[0].NazivOpstine),
		I:  amount,
		Sf: "253",
		S:  "Porez na Imovinu od Fizickih Lica",
		Ro: strings.ReplaceAll(upit.UpitStanjaSaldoOpstList[0].PozivNaBroj, " ", ""),
	}
}

// POST Request and get QR Code from Body
// return str
func GenerateQR(qrbody QRBody, auth string) string {
	responseBody, err := JSONPostHeader(qrCodeURL, &qrbody, map[string]string{
		"Content-Type":  "application/json",
		"authorization": auth,
	})
	if err != nil {
		log.Fatal(err)
	}
	return responseBody

}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

/*
Usage:

	var p RouteParams

	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
*/
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
