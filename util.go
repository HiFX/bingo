package bingo

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/hifx/banner"
	"github.com/hifx/bingo/infra/log"
	"github.com/hifx/errgo"
	"github.com/hifx/graceful"
	"goji.io/pat"
	"golang.org/x/net/context"
)

//PrintName prints the app name
func PrintName(str string) {
	color.New(color.FgYellow).Add(color.Bold).Println(banner.PrintS(str))
}

// PrintBanner prints the banner
func PrintBanner(a ...interface{}) {
	hd := color.New(color.FgGreen).Add(color.Bold)
	hd.Println("----------------------------------")
	hd.Println(a...)
	hd.Println("----------------------------------")
}

// PrintError prints the error
func PrintError(a ...interface{}) {
	hd := color.New(color.FgRed).Add(color.Bold)
	hd.Println("----------------------------------")
	hd.Println(a...)
	hd.Println("----------------------------------")
}

//Run gracefully starts the http server
func Run(addr string, timeout time.Duration, h http.Handler) {
	http.Handle("/", h)
	graceful.Run(addr, timeout, http.DefaultServeMux)
}

//BoundParam returns the bound parameter with the given name. Wraps around goji's pat.Param
func BoundParam(ctx context.Context, name string) string {
	return pat.Param(ctx, name)
}

// JSONW writes JSON response to the given writer
func JSONW(w http.ResponseWriter, status int, l log.Logger, data interface{}) {
	d, err := json.MarshalIndent(data, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(d)
	if nil != err {
		l.Error("web.ioerror", err.Error())
	}
	_, err = w.Write([]byte("\n"))
	if nil != err {
		l.Error("web.ioerror", err.Error())
	}
}

//Close an io.Closer object by handling it's error.
// Example:
// 		h, err := app.Init(env)
// 		defer utils.Close(h, conf.Config.Log())
func Close(closer io.Closer, log log.Logger) {
	if err := closer.Close(); err != nil {
		log.Crit("Failed to Close Object: %#v\n Error: %s ", err.Error())
	}
}

// ReadBody decodes request.Body into the destination.
//  Warning: des should be a pointer
func ReadBody(request *http.Request, des interface{}, msg string) error {
	if kind := reflect.TypeOf(des).Kind(); kind != reflect.Ptr {
		return errgo.Errorf("des should be a pointer: but got %#v", kind)
	}
	if err := json.NewDecoder(request.Body).Decode(des); err != nil {
		return errgo.Errorf("Failed to decode request body,Error: %s", err.Error())
	}
	return nil
}
