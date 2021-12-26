package prompt

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Recipient struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	VatNumber string `json:"vatNumber"`
}

var recipients []Recipient

func PopulateList() []Recipient {
	if len(recipients) > 0 {
		return recipients
	}

	recipientsJson, err := ioutil.ReadFile("recipients.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(recipientsJson, &recipients); err != nil {
		log.Fatal(err)
	}

	return recipients
}
