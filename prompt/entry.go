package prompt

import "fmt"

type Entry struct {
	NetHPrice   float64
	Unit        string
	Amount      int
	Description string
}

func (e Entry) GrossHPrice() float64 {
	return e.NetHPrice * 1.23
}

func (e Entry) GrossPrice() float64 {
	return e.GrossHPrice() * float64(e.Amount)
}

func (e Entry) NetHPriceStr() string {
	return fmt.Sprintf("%.2f", e.NetHPrice)
}

func (e Entry) NetPrice() float64 {
	return e.NetHPrice * float64(e.Amount)
}

func (e Entry) NetPriceStr() string {
	return fmt.Sprintf("%.2f", e.NetPrice())
}

func (e Entry) GrossPriceStr() string {
	return fmt.Sprintf("%.2f", e.GrossPrice())
}

func (e Entry) VatPriceStr() string {
	return fmt.Sprintf("%.2f", e.GrossPrice()-e.NetPrice())
}

func (e Entry) Dump() {
	fmt.Printf("%s: %v PLN * %v%s \n", e.Description, e.NetHPrice, e.Amount, e.Unit)
}
