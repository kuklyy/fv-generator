package main

import (
	"bytes"
	"context"
	"fmt"
	"fv-generator/prompt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	fv := prompt.PopulateOptions()

	baseTemplate, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
	}
	var parsedTemplate bytes.Buffer

	if err := baseTemplate.Execute(&parsedTemplate, fv); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("temp.html", parsedTemplate.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove("temp.html")
	}()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("file://%s/%s", cwd, "temp.html")
	parsedTemplate.Reset()

	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	if err := chromedp.Run(taskCtx, grabPdf(filePath, &parsedTemplate)); err != nil {
		log.Fatal(err)
	}

	outputName := strings.Replace(fmt.Sprintf("fv-%s.pdf", fv.NO), "/", "-", -1)

	if err := ioutil.WriteFile("output/"+outputName, parsedTemplate.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File saved to output/%s at %v\n", outputName, time.Now().String())
}

func grabPdf(url string, res *bytes.Buffer) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
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
}
