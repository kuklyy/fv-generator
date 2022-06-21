package prompt

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

func PopulateOptions() FV {
	no := getNo()
	recipient := getRecipient()
	entries := getEntries()
	return FV{
		Entries:   entries,
		Recipient: recipient,
		CreatedAt: time.Now(),
		NO:        no,
	}
}

func getNo() string {
	prompt := promptui.Prompt{
		Label: "NO",
		Validate: func(s string) error {
			if len(s) < 1 {
				return errors.New("NO is required")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func getRecipient() Recipient {
	recipients := PopulateList()
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
	entry := getEntry()
	entries = append(entries, entry)

	for !done {
		for _, e := range entries {
			e.Dump()
		}

		// weird implementation of confirm prompts
		// err means false/empty
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
		Label:   "Description",
		Default: "Stała współpraca w&nbsp;zakresie usług informatycznych",
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
