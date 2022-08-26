package prompt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

func PopulateOptions() FV {
	no := getNo()
	recipient := getRecipient()
	entries := getEntries()
	return NewFV(entries, recipient, no)
}

func getNo() string {
	now := time.Now()

	return fmt.Sprintf("01/0%d/%d", now.Month(), now.Year())
}

func getRecipient() Recipient {
	recipients := PopulateList()
	switch len(recipients) {
	case 0:
		panic("recipient.json is empty")
	case 1:
		return recipients[0]
	}

	prompt := promptui.Select{
		HideHelp: true,
		Label:    "Select Recipient",
		Items:    PopulateList(),
		Templates: &promptui.SelectTemplates{
			Inactive: "{{ .Name }}",
			Active:   "\U00002705 {{ .Name | cyan }}",
			Selected: "Recipient: {{ .Name }}"},
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	return recipients[i]
}

func getEntries() []Entry {
	done := false
	entries := make([]Entry, 0, 16)
	useDefault := os.Getenv("FV_USE_DEFAULT")
	var entry Entry
	if useDefault != "" {
		entry = defaultEntry()
	} else {
		entry = getEntry()
	}
	entries = append(entries, entry)

	for !done {
		for _, e := range entries {
			e.Dump()
		}

		prompt := promptui.Prompt{
			Label:     "Add another entry",
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			done = true
			break
		}
		entry = getEntry()
		entries = append(entries, entry)
	}
	return entries
}

func getEntry() Entry {
	prompt := promptui.Prompt{
		Label: "Description",
		Validate: func(s string) error {
			if len(s) < 1 {
				return errors.New("description is required")
			}
			return nil
		},
	}

	description, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	prompt = promptui.Prompt{
		Label: "Unit",
		Validate: func(s string) error {
			if len(s) < 1 {
				return errors.New("unit is required")
			}
			return nil
		},
	}

	unit, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	prompt = promptui.Prompt{
		Label: "Amount",
		Validate: func(s string) error {
			_, err := strconv.Atoi(s)
			if err != nil {
				return errors.New("amount is invalid number")
			}
			return nil
		},
	}

	amountStr, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		log.Fatal(err)
	}

	prompt = promptui.Prompt{
		Label: "Hourly price",
		Validate: func(s string) error {
			_, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return errors.New("hourly price is invalid float number")
			}

			return nil
		},
	}

	hourlyPriceStr, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	hourlyPrice, err := strconv.ParseFloat(hourlyPriceStr, 64)
	if err != nil {
		log.Fatal(err)
	}

	return Entry{
		NetHPrice:   hourlyPrice,
		Unit:        unit,
		Amount:      amount,
		Description: description,
	}

}

func defaultEntry() Entry {
	defaultDescription := os.Getenv("FV_DEFAULT_PROMPT")
	defaultPrice := os.Getenv("FV_DEFAULT_PRICE")
	defaultUnit := os.Getenv("FV_DEFAULT_UNIT")
	defaultAmount := os.Getenv("FV_DEFAULT_AMOUNT")

	price, err := strconv.ParseFloat(defaultPrice, 64)
	if err != nil {
		panic("FV_DEFAULT_PRICE is invalid float")
	}

	amount, err := strconv.Atoi(defaultAmount)
	if err != nil {
		panic("FV_DEFAULT_AMOUNT is invalid integer")
	}

	return Entry{
		Description: defaultDescription,
		NetHPrice:   price,
		Unit:        defaultUnit,
		Amount:      amount,
	}
}
