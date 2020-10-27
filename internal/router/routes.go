package router

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/NodeFactoryIo/vedran/internal/auth"
	"github.com/NodeFactoryIo/vedran/internal/controllers"
	"github.com/gorilla/mux"
)

func createRoute(route string, method string, handler http.HandlerFunc, router *mux.Router, authorized bool) {
	var r *mux.Route
	if authorized {
		r = router.Handle(route, auth.AuthMiddleware(handler))
	} else {
		r = router.Handle(route, handler)
	}
	r.Methods(method)
	r.Name(route)
	log.Debugf("Created route %s\t%s", method, route)
}

func createRoutes(apiController *controllers.ApiController, router *mux.Router) {
	createRoute("/", "POST", apiController.RPCHandler, router, false)

	createRoute("/api/v1/nodes", "POST", apiController.RegisterHandler, router, false)
	createRoute("/api/v1/nodes/pings", "POST", apiController.PingHandler, router, true)
	createRoute("/api/v1/nodes/metrics", "PUT", apiController.SaveMetricsHandler, router, true)
}
