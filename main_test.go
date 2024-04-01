package main

import (
	"reflect"
	"testing"
)

/*
	{
		"K": "PR",
		"V": "01",
		"C": "1",
		"R": "840000071312184357",
		"N": "LPA Beograd-Čukarica",
		"I": "RSD253,00",
		"SF": "253",
		"S": "Porez na Imovinu od Fizickih Lica",
		"RO": "97590111609994710212"
	}
*/
func TestTransformerUpitaStanja(t *testing.T) {
	input := UpitStanja{
		Pib:              "1609994710212",
		DatumZaduzenjaDo: "2024-03-30",
		DatumUplateDo:    "2024-03-30",
		IsObveznik:       true,
		UpitStanjaSaldoOpstList: []UpitStanjaSaldoOpstList{
			{
				UpitStanjaSaldoList: []UpitStanjaSaldoList{
					{
						Racun:            "713121",
						RacunCeo:         "840-713121843-57",
						RacunOpis:        "Porez na imovinu obveznika koji ne vode poslovne knjige",
						SaldoDuguje:      1030.49,
						SaldoPotrazuje:   1168.94,
						KamataZaduzenje:  0.0,
						KamataObracunata: 13.45,
						KamataNaplacena:  13.45,
						SaldoGlavnica:    -125.0,
						SaldoKamata:      0.0,
						SaldoUkupan:      -125.0,
						ListaPromena: []ListaPromena{
							{
								KnjPromSifra:     "20",
								KnjPromSifraIP:   "102",
								KnjPromSifraZp:   "2",
								KnjPromDISSifra:  "",
								KnjPromDISOpis:   "",
								BrojNaloga:       "0",
								DisDokument:      "null",
								Datum:            "2024-01-01",
								PrometDuguje:     0.0,
								PrometPotrazuje:  2.03,
								KamataZaduzenje:  0.0,
								KamataObracunata: 0.0,
								KamataNaplacena:  0.0,
								SaldoGlavnica:    -2.03,
								SaldoKamata:      0.0,
								KnjPromOpis:      "Saldo preplata",
								KnjPromOpisPu:    "Почетно стање - преплата",
							}}}},

				SifraOpstine:  "11",
				NazivOpstine:  "Beograd-Čukarica",
				PozivNaBroj:   "97 590111609994710212",
				ObveznikIdent: "1609994710212",
				DatumUpita:    "2024-03-30",
				VremeObrade:   "02:25",
			},
		}}

	want := QRBody{
		K:  "PR",
		V:  "01",
		C:  "1",
		R:  "840000071312184357",
		N:  "LPA Beograd-Čukarica",
		I:  "RSD-125,00",
		Sf: "253",
		S:  "Porez na Imovinu od Fizickih Lica",
		Ro: "97590111609994710212",
	}
	result := TransformerUpitaStanja(input)
	if !reflect.DeepEqual(result, want) {
		t.Log("Result was incorrect.")
		t.Logf("Got:\n %#v", result)
		t.Logf("Want:\n %#v", want)
		t.Fail()

	}

}
