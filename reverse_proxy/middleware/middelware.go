package middleware

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type HttpMiddewareAdapter func(http.Handler) http.Handler

func LogRequest(logger *log.Logger) HttpMiddewareAdapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s was called\n", r.URL.String())
			startTime := time.Now()
			h.ServeHTTP(w, r)
			endTime := time.Now()
			diff := endTime.Sub(startTime)
			logger.Printf("%s answered in %s\n", r.URL.String(), diff.String())
		})
	}
}

func DumpRequest(logger *log.Logger) HttpMiddewareAdapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dump, _ := httputil.DumpRequest(r, true)
			logger.Printf("****REQUEST****\n%s\n", dump)
			h.ServeHTTP(w, r)
		})
	}
}

func ChainHttpMiddlewareAdapters(
	h http.Handler,
	httpMiddewareAdapters []HttpMiddewareAdapter,
) http.Handler {
	numberOfAdapters := len(httpMiddewareAdapters)
	for index := numberOfAdapters - 1; index >= 0; index-- {
		adapter := httpMiddewareAdapters[index]
		h = adapter(h)

	}

	return h
}

type HttpResponseMiddewareAdapter func(*http.Response) error

func Cors() HttpResponseMiddewareAdapter {
	return func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "*")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		return nil
	}
}

func ChainHttpResponsetMiddlewareAdapters(
	httpResponseMiddewareAdapters []HttpResponseMiddewareAdapter,
) HttpResponseMiddewareAdapter {
	return func(r *http.Response) error {
		for _, adapter := range httpResponseMiddewareAdapters {
			err := adapter(r)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
