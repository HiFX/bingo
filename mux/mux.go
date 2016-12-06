/*
Package mux is a wrapper over goji.Mux

It adds a predefined list of middlewares to muxes and sub-muxes that would be created via the New() and Sub() functions

e.g. usage

	acslog := log.NewJSONLogger(conf.Log.Access)
	errlog := log.NewJSONLogger(conf.Log.Err)
	mux.Init(acslog, errlog)

To set custom middleware like logger, 404, metrics and 404 handlers to all muxes, use

	acslog := log.NewJSONLogger(conf.Log.Access)
	errlog := log.NewJSONLogger(conf.Log.Err)
	mux.SetLogs(acslog, errlog)

	mux.SetSubMware(
			middleware.ApplySubStats,
			middleware.Apply404)

	mux.SetMware(
		middleware.ApplyLog(acsslog),
		middleware.Apply404,
		middleware.ApplyStats,
		middleware.ApplyReqID)

*/
package mux

import (
	"net/http"

	"github.com/hifx/bingo/infra/log"
	"github.com/hifx/bingo/middleware"
	"github.com/hifx/errgo"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

var mlist []func(goji.Handler) goji.Handler
var submlist []func(goji.Handler) goji.Handler
var acslog, errlog log.Logger

//Mux is a wrapper over Goji's mux
type Mux struct {
	*goji.Mux
}

// Get dispatches to the given handler when the pattern matches and the HTTP
// method is GET.
func (m *Mux) Get(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Get(pattern), wrap(h))
}

// Post dispatches to the given handler when the pattern matches and the HTTP
// method is POST.
func (m *Mux) Post(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Post(pattern), wrap(h))
}

// Put dispatches to the given handler when the pattern matches and the HTTP
// method is PUT.
func (m *Mux) Put(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Put(pattern), wrap(h))
}

// Patch dispatches to the given handler when the pattern matches and the HTTP
// method is PATCH.
func (m *Mux) Patch(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Patch(pattern), wrap(h))
}

// Delete dispatches to the given handler when the pattern matches and the HTTP
// method is DELETE.
func (m *Mux) Delete(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Delete(pattern), wrap(h))
}

// Options dispatches to the given handler when the pattern matches and the HTTP
// method is OPTIONS.
func (m *Mux) Options(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Options(pattern), wrap(h))
}

// Head dispatches to the given handler when the pattern matches and the HTTP
// method is HEAD.
func (m *Mux) Head(pattern string, h func(context.Context, http.ResponseWriter, *http.Request) error) {
	m.HandleFuncC(pat.Head(pattern), wrap(h))
}

//wrap helps make application handlers  satisfy goji's type HandlerFunc.
//Any error returned by bingo's app handler's would be logged to the error log
//TODO: Log "req_id", "method", "uri", "remote", "err", "stack"
func wrap(h func(context.Context, http.ResponseWriter, *http.Request) error) func(context.Context, http.ResponseWriter, *http.Request) {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := h(ctx, w, r)
		if err != nil {
			reqid := middleware.GetReqID(ctx)
			switch e := err.(type) {
			case *errgo.Err:
				errlog.Error(
					"req_id", reqid,
					"uri", r.RequestURI,
					"method", r.Method,
					"remote", r.RemoteAddr,
					"error", e.Error(),
					"stack", e.Stack(),
				)
				w.Header().Set("Content-Type", e.ContentType())
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(code)
				w.Write([]byte(e.Message()))
			default:
				errlog.Error(
					"req_id", reqid,
					"uri", r.RequestURI,
					"method", r.Method,
					"remote", r.RemoteAddr,
					"error", err.Error(),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
	return fn
}

/*
New is a wrapper over goji.NewMux(). It adds
the pre-defined list of middlewares to the mux
*/
func New() *Mux {
	m := goji.NewMux()
	for _, mware := range mlist {
		m.UseC(mware)
	}

	return &Mux{m}
}

/*
Sub is a wrapper over goji.SubMux(). It adds
the pre-defined list of middlewares to the submux
*/
func Sub() *Mux {
	m := goji.SubMux()
	for _, mware := range submlist {
		m.UseC(mware)
	}
	return &Mux{m}
}

// SetMware sets the middlewares to be used for all muxes
func SetMware(m ...func(goji.Handler) goji.Handler) {
	mlist = m
}

// SetSubMware sets the middlewares to be used for all sub-muxes
func SetSubMware(m ...func(goji.Handler) goji.Handler) {
	submlist = m
}

//Init initializes the mux package. It initializes the middlewares to be used by Muxes & SubMuxes
//and sets the loggers. An application can overwrite the middlewares by calling SetMware & SetSubMware
func Init(acslog, errlog log.Logger) {
	SetMware(
		middleware.CrossDomainRequestAllower,
		middleware.ApplyReqID,
		middleware.ApplyRecoverer(errlog),
		middleware.ApplyLog(acslog),
		middleware.Apply404,
		middleware.ApplyStats,
	)
	SetSubMware(
		middleware.ApplySubStats,
		middleware.Apply404,
	)
	SetLogs(acslog, errlog)
}

//SetLogs sets the loggers used by the mux package
func SetLogs(a, e log.Logger) {
	acslog = a
	errlog = e
}
