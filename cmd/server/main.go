package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.IgnoreCertErrors,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	err := chromedp.Run(
		taskCtx,
		chromedp.Navigate("https://localhost:5000"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
