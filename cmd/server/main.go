package main

import (
	"context"

	"cdr.dev/slog"
	"github.com/chromedp/chromedp"

	"github.com/ib-gambler/ib-cp-server/pkg"
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

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(
		taskCtx,
		chromedp.Navigate("https://localhost:5000"),
	)
	if err != nil {
		pkg.Logger.Fatal(context.Background(), "chromedp.Run", slog.Error(err))
	}
}
