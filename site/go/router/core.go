package router

import (
	"site/go/app"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/sjc5/river/kit/middleware/healthcheck"
	"github.com/sjc5/river/kit/middleware/robotstxt"
	"github.com/sjc5/river/kit/mux"
	"github.com/sjc5/river/kit/tasks"
)

var sharedTasksRegistry = tasks.NewRegistry()

func Core() *mux.Router {
	r := mux.NewRouter(nil)

	mux.SetGlobalHTTPMiddleware(r, chimw.Logger)
	mux.SetGlobalHTTPMiddleware(r, chimw.Recoverer)
	mux.SetGlobalHTTPMiddleware(r, healthcheck.Healthz)
	mux.SetGlobalHTTPMiddleware(r, robotstxt.Allow)
	mux.SetGlobalHTTPMiddleware(r, app.Kiruna.FaviconRedirect())

	// static public assets
	mux.RegisterHandler(r, "GET", app.Kiruna.GetPublicPathPrefix()+"*", app.Kiruna.MustGetServeStaticHandler(true))

	// river UI routes
	mux.RegisterHandler(r, "GET", "/*", app.River.GetUIHandler(UIRouter))

	// river API routes
	actionsHandler := app.River.GetActionsHandler(ActionsRouter)
	mux.RegisterHandler(r, "GET", ActionsRouter.MountRoot("*"), actionsHandler)
	mux.RegisterHandler(r, "POST", ActionsRouter.MountRoot("*"), actionsHandler)

	return r
}
