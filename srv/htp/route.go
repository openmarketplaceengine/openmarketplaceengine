package htp

import (
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"
)

type (
	Handler     = http.Handler
	HandlerFunc = http.HandlerFunc
	Middleware  = func(h http.Handler) http.Handler
)

type Routes struct {
	path       string
	uses       []Middleware
	hand       []*patternHandler
	meth       []*methodHandler
	with       []*Routes
	mount      routeArray
	notFound   HandlerFunc
	notAllowed HandlerFunc
}

//-----------------------------------------------------------------------------

type routeArray struct {
	routes []*Routes
}

type patternHandler struct {
	pattern string
	handler Handler
}

type methodHandler struct {
	method  string
	pattern string
	handler Handler
}

//-----------------------------------------------------------------------------
// Handlers
//-----------------------------------------------------------------------------

// Handle adds the route `pattern` that matches any http method to
// execute the `handler` http.Handler.
func (r *Routes) Handle(pattern string, handler Handler) {
	r.hand = append(r.hand, &patternHandler{pattern, handler})
}

// HandleFunc adds the route `pattern` that matches any http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) HandleFunc(pattern string, handlerFn HandlerFunc) {
	r.hand = append(r.hand, &patternHandler{pattern, handlerFn})
}

// Method adds the route `pattern` that matches `method` http method to
// execute the `handler` http.Handler.
func (r *Routes) Method(method, pattern string, handler Handler) {
	r.meth = append(r.meth, &methodHandler{method, pattern, handler})
}

// MethodFunc adds the route `pattern` that matches `method` http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) MethodFunc(method, pattern string, handlerFn HandlerFunc) {
	r.meth = append(r.meth, &methodHandler{method, pattern, handlerFn})
}

//-----------------------------------------------------------------------------
// Sub Routing
//-----------------------------------------------------------------------------

// Mount attaches a subrouter along a routing path.
// It's very useful to split up a large API as many independent routers
// and compose them as a single service using Mount.
func (r *Routes) Mount(pattern string) *Routes {
	return r.mount.getOrAdd(pattern)
}

// Use appends a middleware handler to the middleware stack.
//
// The middleware stack for any path will execute before searching for a matching
// route to a specific handler, which provides opportunity to respond early,
// change the course of the request execution, or set request-scoped values for
// the next http.Handler.
func (r *Routes) Use(middlewares ...Middleware) {
	r.uses = append(r.uses, middlewares...)
}

// With adds inline middlewares for an endpoint handler.
func (r *Routes) With(middlewares ...Middleware) *Routes {
	with := &Routes{uses: middlewares}
	r.with = append(r.with, with)
	return with
}

//-----------------------------------------------------------------------------
// HTTP Errors
//-----------------------------------------------------------------------------

// NotFound sets a custom http.HandlerFunc for routing paths that could
// not be found. The default 404 handler is `http.NotFound`.
func (r *Routes) NotFound(hfunc HandlerFunc) {
	r.notFound = hfunc
}

// MethodNotAllowed sets a custom http.HandlerFunc for routing paths where the
// method is unresolved. The default handler returns a 405 with an empty body.
func (r *Routes) MethodNotAllowed(hfunc HandlerFunc) {
	r.notAllowed = hfunc
}

//-----------------------------------------------------------------------------
// Method Shortcuts
//-----------------------------------------------------------------------------

// Connect adds the route `pattern` that matches a CONNECT http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Connect(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("CONNECT", pattern, handlerFn)
}

// Delete adds the route `pattern` that matches a DELETE http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Delete(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("DELETE", pattern, handlerFn)
}

// Get adds the route `pattern` that matches a GET http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Get(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("GET", pattern, handlerFn)
}

// Head adds the route `pattern` that matches a HEAD http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Head(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("HEAD", pattern, handlerFn)
}

// Options adds the route `pattern` that matches a OPTIONS http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Options(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("OPTIONS", pattern, handlerFn)
}

// Patch adds the route `pattern` that matches a PATCH http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Patch(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("PATCH", pattern, handlerFn)
}

// Post adds the route `pattern` that matches a POST http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Post(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("POST", pattern, handlerFn)
}

// Put adds the route `pattern` that matches a PUT http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Put(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("PUT", pattern, handlerFn)
}

// Trace adds the route `pattern` that matches a TRACE http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *Routes) Trace(pattern string, handlerFn HandlerFunc) {
	r.MethodFunc("TRACE", pattern, handlerFn)
}

//-----------------------------------------------------------------------------
// Build
//-----------------------------------------------------------------------------

// Build creates actual router implementation with chi.Router.
func (r *Routes) Build() Handler {
	return chiMux(r, chi.NewMux())
}

//-----------------------------------------------------------------------------
// ChiMux
//-----------------------------------------------------------------------------

func chiMux(r *Routes, c *chi.Mux) *chi.Mux {
	if len(r.uses) > 0 {
		c.Use(r.uses...)
	}
	for _, w := range r.with {
		chiMux(w, c.With().(*chi.Mux))
	}
	for _, m := range r.mount.routes {
		c.Mount(m.path, chiMux(m, chi.NewMux()))
	}
	for _, h := range r.hand {
		c.Handle(h.pattern, h.handler)
	}
	for _, m := range r.meth {
		c.Method(m.method, m.pattern, m.handler)
	}
	if r.notFound != nil {
		c.NotFound(r.notFound)
	}
	if r.notAllowed != nil {
		c.MethodNotAllowed(r.notAllowed)
	}
	return c
}

//-----------------------------------------------------------------------------
// Routes Array
//-----------------------------------------------------------------------------

func (r *routeArray) Len() int {
	return len(r.routes)
}

func (r *routeArray) Less(i, j int) bool {
	return r.routes[i].path < r.routes[j].path
}

func (r *routeArray) Swap(i, j int) {
	r.routes[i], r.routes[j] = r.routes[j], r.routes[i]
}

func (r *routeArray) sort() {
	sort.Sort(r)
}

//-----------------------------------------------------------------------------

func (r *routeArray) find(path string) int {
	a := r.routes
	n := len(a)
	x := sort.Search(n, func(i int) bool {
		return a[i].path >= path
	})
	if x < n && a[x].path == path {
		return x
	}
	return -1
}

//-----------------------------------------------------------------------------

func (r *routeArray) appendRoutes(path string) *Routes {
	rts := &Routes{path: path}
	r.routes = append(r.routes, rts)
	r.sort()
	return rts
}

//-----------------------------------------------------------------------------

func (r *routeArray) getOrAdd(path string) *Routes {
	if x := r.find(path); x != -1 {
		return r.routes[x]
	}
	return r.appendRoutes(path)
}
