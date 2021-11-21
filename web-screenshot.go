package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/kennygrant/sanitize"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: web-screenshot <url>")
	}
	myURL := os.Args[1]
	if myURL == "" {
		log.Fatal("No URL provided")
	} else {
		log.Printf("URL: %s", myURL)
	}
	// create contexts
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	ctx2, cancel2 := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	defer cancel2()

	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx2, fullScreenshot(myURL, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse(myURL)
	if err != nil {
		log.Fatal(err)
	}
	imageFile := sanitize.Path("screenshot-" + u.Host + u.Path + ".png")
	imageFile = strings.Replace(imageFile, "/", "-", -1)
	if err := ioutil.WriteFile(imageFile, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", imageFile)
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Reset
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Poll("document.fonts.ready", res),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.FullScreenshot(res, quality),
	}
}
