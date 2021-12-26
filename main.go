package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type Context struct {
	NO          string
	CompanyName string
}

func main() {
	baseTemplate, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
	}
	var parsedTemplate bytes.Buffer

	if err := baseTemplate.Execute(&parsedTemplate, Context{
		NO:          "firsdtname",
		CompanyName: "lastname",
	}); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("test.html", parsedTemplate.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove("test.html")
	}()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("file://%s/%s", cwd, "test.html")
	parsedTemplate.Reset()

	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	if err := chromedp.Run(taskCtx, grabPdf(filePath, &parsedTemplate)); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("test.pdf", parsedTemplate.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
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
