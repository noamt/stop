package stop

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

//This method registers a handler on the /shutdown path, with the given HTTP server
func OnHttpCall(mux *http.ServeMux, f http.HandlerFunc) {
	mux.HandleFunc("/shutdown", f)
}

//This method starts a background HTTP server on the given host and port.
//The given handler is registered on the /shutdown path.
//Once handling the call, the background HTTP server is shutdown
func OnHttpServerCall(host string, port int, f http.HandlerFunc) {
	m := http.NewServeMux()
	s := http.Server{Addr: host + ":" + strconv.Itoa(port), Handler: m}
	OnHttpCall(m, func(writer http.ResponseWriter, request *http.Request) {
		f(writer, request)
		go func() {
			s.Shutdown(context.Background())
		}()
	})

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
}
