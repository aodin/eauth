package eauth

import (
	"net/http"
)

type Server struct {
	config   Config
	users    UserManager
	sessions SessionManager
}

// This is a massive security hole if left enabled
func (s *Server) Generate(w http.ResponseWriter, r *http.Request) {
	// Generate a new link for user 1
	link := CreateLink(s.config, s.users.GetId(1))
	w.Write([]byte(link))
}

func (s *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Split the request url apart
	// Determine if the link is valid
	// If valid, Create a session and add the cookie
	w.Write([]byte("invalid"))
}

func (s *Server) ListenAndServe() error {
	// Build the server address
	return http.ListenAndServe(s.config.Domain, nil)
}

func NewServer(config Config) (s *Server) {
	s.config = config
	s.users = InMemoryUsers()
	s.sessions = InMemorySessions()

	// Attach a test user
	admin := User{}
	if err := s.users.Create(admin); err != nil {
		panic(err)
	}

	// Attach the handlers
	http.HandleFunc("/gen/", s.Generate)
	http.HandleFunc("/auth/", s.Authenticate)
	return
}
