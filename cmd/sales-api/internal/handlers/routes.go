package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/joshakeman/service-practice/internal/mid"
	"github.com/joshakeman/service-practice/internal/platform/web"
)

// API ...
func API(build string, shutdown chan os.Signal, log *log.Logger) *web.App {
	// tm := httptreemux.New()
	app := web.NewApp(shutdown, mid.Logger(log))

	app.Handle(http.MethodGet, "/test", health)

	return app
}
