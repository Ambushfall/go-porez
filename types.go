package main

type Time struct {
	CurrentTime string `json:"current_time"`
}

type ErrResponse struct {
	Resonse string `json:"response"`
	Code    int    `json:"code"`
	Client  string `json:"client"`
}

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

type QRBody struct {
	K  string `json:"K"`
	V  string `json:"V"`
	C  string `json:"C"`
	R  string `json:"R"`
	N  string `json:"N"`
	I  string `json:"I"`
	Sf string `json:"SF"`
	S  string `json:"S"`
	Ro string `json:"RO"`
}

type RouteParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
