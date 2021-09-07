package proxy

import (
	"congress/logger"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ProxyServer struct {
	hosts map[string]*Host
}

var log = logger.Default.Child("proxy")

func New() *ProxyServer {
	return &ProxyServer{make(map[string]*Host)}
}

func (s *ProxyServer) AddHost(h *Host) {
	s.hosts[h.name] = h
}

func (s *ProxyServer) Print() {
	for _, host := range s.hosts {
		log.Debug("%s: %s", host.name, host.defaultBackend)
		for _, path := range host.paths {
			log.Debug("  %s: %s", path.pattern, path.backend)
		}
	}
}

func (s *ProxyServer) Listen(port uint) error {
	log.Info("Listening at port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}

func (s *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("---> %s http://%s%s", r.Method, r.Host, r.URL)

	if host, ok := s.hosts[hostname(r.Host)]; ok {
		for _, path := range host.paths {
			if path.Matches(r.URL) {
				url, err := path.Replace(r.URL)
				if err != nil {
					log.Error("Failed to rewrite url: %s", err)
				}
				rewriteRequest(w, r, url)
				return
			}
		}

		if host.defaultBackend != "" {
			url, err := host.ReplaceDefault(r.URL)
			if err != nil {
				log.Error("Failed to rewrite url: %s", err)
			}
			rewriteRequest(w, r, url)
			return
		}
	}

	log.Debug("<--- %v %s", 404, r.URL)
	http.NotFound(w, r)
}

func hostname(host string) string {
	return strings.Split(host, ":")[0]
}

func rewriteRequest(w http.ResponseWriter, r *http.Request, url *url.URL) {
	req, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		log.Error("Failed to create request: %s", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Failed to execute request: %s", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	} else {
		log.Debug("<--- %v %s", res.Status, res.Request.URL)
		for key, values := range res.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	}
}
