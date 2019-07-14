package routes

import (
    "net/http"
)

type Route struct {
    Method string
    Version int
    Endpoint string
    Handler func(http.ResponseWriter, *http.Request)
}

// NewRoute returns the route with default api version: 1.
func NewRoute(method string, endpoint string,
    handler func(http.ResponseWriter, *http.Request)) *Route {
    return &Route{
        Method: method,
        Version: 1,
        Endpoint: endpoint,
        Handler: handler,
    }
}

// NewRouteWithVersion returns the route with specified version.
func NewRouteWithVersion(method string, endpoint string, version int,
    handler func(http.ResponseWriter, *http.Request)) *Route {
    route := NewRoute(method, endpoint, handler)
    route.Version = version
    return route
}
