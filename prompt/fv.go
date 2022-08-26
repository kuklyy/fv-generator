package prompt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FV struct {
	Entries   []Entry
	Recipient Recipient
	CreatedAt time.Time
	NO        string
	CSSPath   string
}

func NewFV(entries []Entry, recipient Recipient, no string) FV {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("failed to load binPath")
	}
	return FV{
		Entries:   entries,
		Recipient: recipient,
		NO:        no,
		CreatedAt: time.Now(),
		CSSPath:   fmt.Sprintf("%s/style.css", binPath),
	}
}

func (f FV) TotalNetAmount() float64 {
	var sum float64

	for _, entry := range f.Entries {
		sum += entry.NetPrice()
	}

	return sum
}

func (f FV) TotalGrossAmount() float64 {
	var sum float64

	for _, entry := range f.Entries {
		sum += entry.GrossPrice()
	}

	return sum
}

func (f FV) TotalVatAmountStr() string {
	return fmt.Sprintf("%.2f", f.TotalGrossAmount()-f.TotalNetAmount())
}

func (f FV) TotalNetAmountStr() string {
	return fmt.Sprintf("%.2f", f.TotalNetAmount())
}

func (f FV) TotalGrossAmountStr() string {
	return fmt.Sprintf("%.2f", f.TotalGrossAmount())
}

func (f FV) GetAmountStr() string {
	res, err := http.Get(fmt.Sprintf("https://slownie.pl/%s", f.TotalGrossAmountStr()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc.Find("#dataWord").Text()
}

func (f FV) GetCreateAtDate() string {
	return f.CreatedAt.Format("02.01.2006")
}

func (f FV) GetPayday() string {
	date := f.CreatedAt.AddDate(0, 0, 14)
	return date.Format("02.01.2006")
}

func (f FV) GetServiceDate() string {
	return f.GetCreateAtDate()
}
