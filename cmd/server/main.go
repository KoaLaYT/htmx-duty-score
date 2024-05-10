package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tempates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tempates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Use(logMiddleware())

	e.Static("/public", "./public")
	e.Renderer = &Template{
		tempates: template.Must(template.ParseGlob("./ui/*.html")),
	}

	server := &server{}
	server.InitRoutes(e)

	go func() {
		if err := e.Start(":41988"); err != nil && err != http.ErrServerClosed {
			log.Fatal("echo start: ", err)
		}
	}()

	go func() {
		url := fmt.Sprintf("http://127.0.0.1:41987/_dev/key/%d", time.Now().UnixNano())
		_, err := http.Get(url)
		if err != nil {
			log.Fatal("notify: ", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	<-ctx.Done()

	ctx, timeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeout()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown: ", err)
	}
}

func logMiddleware() echo.MiddlewareFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	})
}
