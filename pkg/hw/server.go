package hw

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

var HTTP_PORT = "8000"

type Server struct {
	TotalRequests       prometheus.CounterVec
	RequestResponceTime prometheus.SummaryVec
	addr                string
	ready               bool
}

func NewServer() *Server {
	p, present := os.LookupEnv("HTTP_PORT")

	if !present {
		p = HTTP_PORT
	}

	tr := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP Requests.",
		},
		[]string{"path"},
	)

	rrt := promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_responce_time",
			Help: "Response latency in seconds.",
		},
		[]string{"path"},
	)

	s := Server{
		TotalRequests:       *tr,
		RequestResponceTime: *rrt,
		addr:                fmt.Sprintf(":%v", p),
	}

	return &s
}

func (s *Server) Start() {
	s.routes()
	log.Print("The service is ready to listen and serve.")
	http.ListenAndServe(s.addr, nil)
}

func (s *Server) routes() {
	http.Handle("/hello", s.useMetrics(s.printHelloWorld))
	http.HandleFunc("/liveness", s.CheckLiviness)
	http.HandleFunc("/readiness", s.CheckReadiness)
	http.Handle("/metrics", promhttp.Handler())
}

func (s *Server) useMetrics(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		f(w, r)
		deltaTime := time.Since(startTime)
		s.RequestResponceTime.WithLabelValues(r.URL.Path).Observe(deltaTime.Seconds())
		s.TotalRequests.WithLabelValues(r.URL.Path).Inc()
	})
}

func (s *Server) printHelloWorld(w http.ResponseWriter, r *http.Request) {
	time.Sleep(90 * time.Millisecond)
	fmt.Fprint(w, "Hello world!")
}

func (s *Server) CheckLiviness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) CheckReadiness(w http.ResponseWriter, _ *http.Request) {
	sts := http.StatusOK
	if !s.ready {
		sts = http.StatusInternalServerError
	}
	w.WriteHeader(sts)
}

// Emulate activity where server can't be serving requests
func (s *Server) EmulateActivity() {
	upSt := 50 * time.Second
	downSt := 10 * time.Second

	for true {
		s.ready = false
		time.Sleep(downSt)
		s.ready = true
		time.Sleep(upSt)
	}

	s.ready = true
}
