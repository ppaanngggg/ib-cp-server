package main

import (
	"context"

	"cdr.dev/slog"
	"github.com/chromedp/chromedp"

	"github.com/ib-gambler/ib-cp-server/pkg/util"
)

func main() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.IgnoreCertErrors,
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(
		taskCtx,
		chromedp.Tasks{
			chromedp.Navigate("https://localhost:5000"),
			// wait for id=xyz-field-username and id=xyz-field-password to appear
			chromedp.WaitVisible("#xyz-field-username", chromedp.ByID),
			chromedp.WaitVisible("#xyz-field-password", chromedp.ByID),
			// wait for button class=btn-primary to appear
			chromedp.WaitVisible(".btn-primary", chromedp.ByQuery),
			// enter username and password
			chromedp.SendKeys(
				"#xyz-field-username", util.Conf.Account.Username, chromedp.ByID,
			),
			chromedp.SendKeys(
				"#xyz-field-password", util.Conf.Account.Password, chromedp.ByID,
			),
			// click the login button
			chromedp.Click(".btn-primary", chromedp.ByQuery),
			// wait for class=xyzblock-notification to appear
			chromedp.WaitVisible(".xyzblock-notification", chromedp.ByQuery),
			chromedp.ActionFunc(
				func(ctx context.Context) error {
					util.Logger.Info(ctx, "Waiting user confirm")
					return nil
				},
			),
			// wait for class=xyzblock-notification to disappear
			chromedp.WaitNotPresent(".xyzblock-notification", chromedp.ByQuery),
		},
	)
	if err != nil {
		util.Logger.Fatal(context.Background(), "chromedp.Run", slog.Error(err))
	}
}
