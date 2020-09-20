package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	ListenOn string `default:"0.0.0.0:8080"`
	TlsCert  string `default:"/etc/az-fx-proxy/certs/cert.pem"`
	TlsKey   string `default:"/etc/az-fx-proxy/certs/key.pem"`
	DacFxUrl string `default:"https://dac-demo-fx.azurewebsites.net/api/AdmissionControlFx"`
}

var config = &Config{}

func main() {
	if err := envconfig.Process("DAC_PROXY", config); err != nil {
		log.Panic("Failed to load configuration", err)
	}

	log.Infoln(config)
	server := GetAdmissionValidationServer(config.ListenOn)
	if err := server.ListenAndServeTLS(config.TlsCert, config.TlsKey); err != nil {
		log.Panic("Listener failed", err)
	}
}

func GetAdmissionValidationServer(listenOn string) *http.Server {
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/", reverseProxyHandler)
	server := &http.Server{
		Handler: mux,
		Addr:    listenOn,
	}

	return server
}

func reverseProxyHandler(res http.ResponseWriter, req *http.Request) {
	log.Infoln("Sending request to function")
	processRequest(config.DacFxUrl, res, req)
}

func processRequest(target string, res http.ResponseWriter, req *http.Request) {
	dacFxUrl, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(dacFxUrl)

	req.URL.Host = dacFxUrl.Host
	req.URL.Scheme = dacFxUrl.Scheme
	req.Host = dacFxUrl.Host

	proxy.ServeHTTP(res, req)
}
