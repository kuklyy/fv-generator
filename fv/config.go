package fv

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

//go:embed templates/*.html
var classicTemplate embed.FS

type Entry struct {
	Amount int     `yaml:"amount"`
	Unit   string  `yaml:"unit"`
	Price  float64 `yaml:"price"`
	Name   string  `yaml:"name"`
	Vat    int     `yaml:"vat"`
}

func (e Entry) GetNetPrice() float64 {
	return float64(e.Amount) * e.Price
}
func (e Entry) GetVatPrice() float64 {
	return e.GetNetPrice() * (float64(e.Vat) / 100)
}
func (e Entry) GetGrossPrice() float64 {
	return e.GetNetPrice() + e.GetVatPrice()
}

type Config struct {
	Recipient struct {
		Name      string `yaml:"name"`
		Address   string `yaml:"address"`
		VatNumber string `yaml:"vatNumber"`
	} `yaml:"recipient"`
	Seller struct {
		City       string `yaml:"city"`
		Name       string `yaml:"name"`
		Address    string `yaml:"address"`
		VatNumber  string `yaml:"vatNumber"`
		BankNumber string `yaml:"bankNumber"`
	} `yaml:"seller"`
	FV struct {
		NO           string  `yaml:"no"`
		Entries      []Entry `yaml:"entries"`
		HeaderPrefix string  `yaml:"headerPrefix"`
	} `yaml:"fv"`
	Now FVTime `yaml:"now"`
}

const FV_CONFIG_PATH_NAME = "FV_CONFIG_PATH"
const dateFormat = time.DateOnly

// NewConfig returns a new Config
// Config path could be overwritten
// by setting the FV_CONFIG_PATH environment variable.
// By default it uses the config.yaml file in the current working directory.
func NewConfig() (Config, error) {
	path := "config.yaml"
	if pathFromEnv, ok := os.LookupEnv(FV_CONFIG_PATH_NAME); ok {
		path = pathFromEnv
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	t := Config{}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return Config{}, err
	}
	if t.FV.NO == "" {
		t.FV.NO = fmt.Sprintf("FV/01/%s", t.Now.Time().Format("01/2006"))
	}

	return t, nil
}

func (c Config) GetCreatedAt() string {
	return c.Now.Time().Format(dateFormat)
}

func (c Config) GetPayday() string {
	return c.Now.Time().Add(time.Hour * 24 * 14).Format(dateFormat)
}

func (c Config) GetAmountStr() string {
	res, err := http.Get(fmt.Sprintf("https://slownie.pl/%.2f", c.GetTotalGrossAmount()))
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

func (c Config) GetTotalGrossAmount() float64 {
	sum := 0.0
	for _, e := range c.FV.Entries {
		sum += e.GetGrossPrice()
	}
	return sum
}

func (c Config) GetTotalVatAmount() float64 {
	sum := 0.0
	for _, e := range c.FV.Entries {
		sum += e.GetVatPrice()
	}
	return sum
}

func (c Config) GetTotalNetAmount() float64 {
	sum := 0.0
	for _, e := range c.FV.Entries {
		sum += e.GetNetPrice()
	}
	return sum
}

func (c Config) getFileName() string {
	return strings.Replace(fmt.Sprintf("%s.pdf", fmt.Sprintf("FV/%s", c.Now.Time().Format("2006/01"))), "/", "-", -1)
}

func (c Config) parseTemplate() ([]byte, error) {
	baseTemplate, err := template.ParseFS(classicTemplate, "*/classic.html")
	if err != nil {
		return nil, err
	}
	parsedTemplate := bytes.NewBuffer(nil)

	if err := baseTemplate.Execute(parsedTemplate, c); err != nil {
		return nil, err
	}

	return parsedTemplate.Bytes(), nil

}

func (c Config) SavePDF() error {
	parsedTemplate, err := c.parseTemplate()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	tempPath := path.Join(cwd, "temp.html")
	if err := os.WriteFile(tempPath, parsedTemplate, 0644); err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tempPath)
	}()

	pdf, err := c.grabPDF(tempPath)
	if err != nil {
		return err
	}

	outputPath := path.Join("output", c.getFileName())

	if err := os.WriteFile(outputPath, pdf, 0644); err != nil {
		return err
	}

	log.Printf("file saved to %s", outputPath)
	return nil

}

func (c Config) grabPDF(path string) ([]byte, error) {
	res := bytes.NewBuffer(nil)

	task := chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("file://%s", path)),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			res.Write(buf)
			return nil
		}),
	}

	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	if err := chromedp.Run(taskCtx, task); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}
