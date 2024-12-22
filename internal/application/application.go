package application

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/H9ekoN/YandexExam_go/pkg/calculation"
)

type Config struct {
	Port string
}

type Server struct {
	port   string
	srv    *http.Server
	router *http.ServeMux
}

func NewServer(cfg *Config) *Server {
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	s := &Server{
		port:   port,
		router: http.NewServeMux(),
	}

	s.routes()

	s.srv = &http.Server{
		Addr:         port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/api/v1/calculate", s.handleCalculate())
}

func (s *Server) handleCalculate() http.HandlerFunc {
	type request struct {
		Expression string `json:"expression"`
	}

	type response struct {
		Result float64 `json:"result,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondWithError(w, http.StatusUnprocessableEntity, "Invalid request body")
			return
		}

		result, err := calculation.Calc(req.Expression)
		if err != nil {
			status := http.StatusUnprocessableEntity
			message := "Expression is not valid"

			switch err.Error() {
			case "на ноль делить нельзя":
				message = "Division by zero is not allowed"
			case "ошибочка в количестве скобок":
				message = "Parentheses are mismatched"
			case "переполнение":
				status = http.StatusInternalServerError
				message = "Internal server error"
			}

			s.respondWithError(w, status, message)
			return
		}

		s.respondWithJSON(w, http.StatusOK, response{Result: result})
	}
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{
		Error: message,
	})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *Server) Port() string {
	return s.port[1:]
}
