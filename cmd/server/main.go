package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/chromedp/chromedp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	logger = slog.Make(sloghuman.Sink(os.Stdout))
	router = chi.NewRouter()
)

type chiLogger struct{}

func (l *chiLogger) Print(v ...interface{}) {
	logger.Info(context.Background(), v[0].(string))
}

func init() {
	// create background program
	// cmd := exec.Command(
	// 	"./clientportal.gw/bin/run.sh", "./clientportal.gw/root/conf.yaml",
	// )
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// if err := cmd.Start(); err != nil {
	// 	logger.Fatal(context.Background(), "cmd.Start", slog.Error(err))
	// }

	// setup router
	router.Use(middleware.Timeout(conf.Server.Timeout))
	router.Use(middleware.Throttle(conf.Server.Throttle))
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{Logger: &chiLogger{}},
	)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get(
		"/login", func(writer http.ResponseWriter, request *http.Request) {
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
						"#xyz-field-username", conf.Account.Username, chromedp.ByID,
					),
					chromedp.SendKeys(
						"#xyz-field-password", conf.Account.Password, chromedp.ByID,
					),
					// click the login button
					chromedp.Click(".btn-primary", chromedp.ByQuery),
					// wait for class=xyzblock-notification to appear
					chromedp.WaitVisible(".xyzblock-notification", chromedp.ByQuery),
					chromedp.ActionFunc(
						func(ctx context.Context) error {
							logger.Info(ctx, "Waiting user confirm")
							return nil
						},
					),
					// wait for class=xyzblock-notification to disappear
					chromedp.WaitNotPresent(".xyzblock-notification", chromedp.ByQuery),
				},
			)
			if err != nil {
				logger.Fatal(context.Background(), "chromedp.Run", slog.Error(err))
			}
		},
	)
	router.HandleFunc(
		"/*", func(writer http.ResponseWriter, request *http.Request) {
			req, err := http.NewRequestWithContext(
				request.Context(), request.Method,
				fmt.Sprintf("https://localhost:5000%s", request.URL.Path), request.Body,
			)
			if err != nil {
				logger.Error(request.Context(), "new request", slog.Error(err))
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			// resend all other requests to localhost:5000
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}, Timeout: conf.Server.Timeout,
			}
			resp, err := client.Do(req)
			if err != nil {
				logger.Error(request.Context(), "resend request", slog.Error(err))
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				raw, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Error(
						request.Context(), "read response body", slog.Error(err),
					)
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
				if _, err = writer.Write(raw); err != nil {
					logger.Error(
						request.Context(), "write response body", slog.Error(err),
					)
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			writer.WriteHeader(resp.StatusCode)
			writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		},
	)
}

func main() {

	// start server
	logger.Info(
		context.Background(),
		fmt.Sprintf("Listening on %s:%v", conf.Server.Host, conf.Server.Port),
	)
	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%v", conf.Server.Host, conf.Server.Port), router,
	); err != nil {
		logger.Fatal(context.Background(), "http.ListenAndServe", slog.Error(err))
	}
}
