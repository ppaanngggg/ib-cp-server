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
	"github.com/go-chi/render"
)

var (
	logger = slog.Make(sloghuman.Sink(os.Stdout))
	router = chi.NewRouter()
)

type chiLogger struct{}

func (l *chiLogger) Print(v ...interface{}) {
	logger.Info(context.Background(), v[0].(string))
}

func renderErr(
	ctx context.Context, w http.ResponseWriter, r *http.Request, msg string, err error,
) {
	logger.Error(ctx, msg, slog.Error(err))
	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(
		w, r, map[string]string{
			"msg": msg,
			"err": err.Error(),
		},
	)
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

	router.Post(
		"/v1/api/login", func(writer http.ResponseWriter, request *http.Request) {
			// chromedp exec options
			opts := append(
				chromedp.DefaultExecAllocatorOptions[:],
				chromedp.DisableGPU,       // disable GPU
				chromedp.IgnoreCertErrors, // ignore certificate errors
			)
			ctx, cancel := chromedp.NewExecAllocator(request.Context(), opts...)
			defer cancel()
			ctx, cancel = chromedp.NewContext(ctx)
			defer cancel()

			err := chromedp.Run(
				ctx,
				chromedp.Tasks{
					chromedp.Navigate(conf.IB.Url),
					// wait for id=xyz-field-username and id=xyz-field-password to appear
					chromedp.WaitVisible("#xyz-field-username", chromedp.ByID),
					chromedp.WaitVisible("#xyz-field-password", chromedp.ByID),
					// wait for button class=btn-primary to appear
					chromedp.WaitVisible(".btn-primary", chromedp.ByQuery),
					// enter username and password
					chromedp.SendKeys(
						"#xyz-field-username", conf.IB.Username, chromedp.ByID,
					),
					chromedp.SendKeys(
						"#xyz-field-password", conf.IB.Password, chromedp.ByID,
					),
					// click the login button
					chromedp.Click(".btn-primary", chromedp.ByQuery),
					// wait for class=xyzblock-notification to appear
					chromedp.WaitVisible(".xyzblock-notification", chromedp.ByQuery),
					chromedp.ActionFunc(
						func(ctx context.Context) error {
							logger.Debug(ctx, "Waiting user confirm")
							return nil
						},
					),
					// wait for class=xyzblock-notification to disappear
					chromedp.WaitNotPresent(".xyzblock-notification", chromedp.ByQuery),
					// TODO: check succeed or failed, it's not easy to confirm
				},
			)
			if err != nil {
				renderErr(ctx, writer, request, "chromedp.Run error", err)
				return
			}
			writer.WriteHeader(http.StatusOK)
		},
	)
	router.HandleFunc(
		"/*", func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			req, err := http.NewRequestWithContext(
				ctx, request.Method, conf.IB.Url+request.URL.Path, request.Body,
			)
			if err != nil {
				renderErr(ctx, writer, request, "http.NewRequestWithContext error", err)
				return
			}
			// resend all other requests to IB
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}, Timeout: conf.Server.Timeout,
			}
			resp, err := client.Do(req)
			if err != nil {
				renderErr(ctx, writer, request, "client.Do error", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
				raw, err := io.ReadAll(resp.Body)
				if err != nil {
					renderErr(ctx, writer, request, "io.ReadAll error", err)
					return
				}
				if _, err = writer.Write(raw); err != nil {
					renderErr(ctx, writer, request, "writer.Write error", err)
					return
				}
			}
			writer.WriteHeader(resp.StatusCode)
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
