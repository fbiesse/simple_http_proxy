package reverse_proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/fbiesse/simple_reverse_proxy/reverse_proxy/middleware"
)

type ProxyHandler struct {
	proxy *httputil.ReverseProxy
}

func (p ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

type Proxy struct {
	url                           url.URL
	port                          uint32
	httpMiddewareAdapters         []middleware.HttpMiddewareAdapter
	httpResponseMiddewareAdapters []middleware.HttpResponseMiddewareAdapter
	logger                        *log.Logger
}

func (p *Proxy) AppendHttpMiddewareAdapter(m middleware.HttpMiddewareAdapter) {
	p.httpMiddewareAdapters = append(p.httpMiddewareAdapters, m)
}

func (p *Proxy) AppendHttpResponseMiddewareAdapter(m middleware.HttpResponseMiddewareAdapter) {
	p.httpResponseMiddewareAdapters = append(p.httpResponseMiddewareAdapters, m)
}

func (p *Proxy) Start() {
	director := func(r *http.Request) {
		r.URL.Scheme = p.url.Scheme
		r.URL.Host = p.url.Host
		r.Host = p.url.Host
	}
	reverseProxy := &httputil.ReverseProxy{Director: director}
	reverseProxyHandler := ProxyHandler{proxy: reverseProxy}
	handler := middleware.ChainHttpMiddlewareAdapters(
		reverseProxyHandler,
		p.httpMiddewareAdapters,
	)
	reverseProxy.ModifyResponse = middleware.ChainHttpResponsetMiddlewareAdapters(p.httpResponseMiddewareAdapters)

	http.Handle("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", p.port), nil)
	if err != nil {
		p.logger.Panicf("Listen: %v", err)
	}
}

func CreateProxy(
	forwardUrl string,
	port uint32,
	logger *log.Logger,
) Proxy {
	if forwardUrl == "" {
		logger.Fatal("forUrl not provided")
	}

	url, err := url.Parse(forwardUrl)
	if err != nil {
		log.Fatal(err)
	}
	if port == 0 {
		logger.Fatal("Listen port not provided")
	}

	return Proxy{
		url:    *url,
		port:   port,
		logger: logger,
	}
}
