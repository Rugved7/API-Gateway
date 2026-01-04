package router

import (
	"net/http"
	"strings"
)

type Router struct {
	routes []Route
}

func New(routes []Route) *Router {
	return &Router{
		routes: routes,
	}
}

// Match returns the upstream for most specific matching route
func (r *Router) Match(req *http.Request) (Route, bool) {
	path := req.URL.Path

	var (
		bestMatch Route
		found     bool
	)

	for _, route := range r.routes {
		if strings.HasPrefix(path, route.Prefix) {
			if !found || len(route.Prefix) > len(bestMatch.Prefix) {
				bestMatch = route
				found = true
			}
		}
	}

	return bestMatch, found
}
