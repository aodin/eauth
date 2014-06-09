package eauth

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	config   Config
	users    UserManager
	sessions SessionManager
}

func (s *Server) Secrets(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("secrets!"))
}

func (s *Server) LoginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session key from cookie
		cookie, err := r.Cookie(s.config.Cookie.Name)
		// Check if the session is valid
		if err != nil || !IsValidSession(s.sessions, cookie.Value) {
			w.Write([]byte("redirect!"))
			return
		}
		h(w, r)
	}
}

// This is a massive security hole if left enabled
func (s *Server) Generate(w http.ResponseWriter, r *http.Request) {
	// Generate a new link for user 1
	link := CreateLink(s.config, s.users.Get(1))
	w.Write([]byte(link))
}

func (s *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Split the request url apart
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 5 {
		w.Write([]byte("Invalid link"))
		return
	}

	// Parse the user id and timestamp
	uid, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		w.Write([]byte("Invalid id"))
		return
	}

	timestamp, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		w.Write([]byte("Invalid timestamp"))
		return
	}

	// Get the user
	user := s.users.Get(uid)
	if user.Id == 0 {
		w.Write([]byte("Invalid id"))
		return
	}

	// Determine if the link is valid
	if !IsValid(s.config, user, time.Unix(timestamp, 0), parts[4]) {
		w.Write([]byte("Failed check"))
		return
	}

	// Create a new session
	session, err := NewSession(s.sessions, user.Id, s.config.Cookie)
	if err != nil {
		w.Write([]byte("Could not create session"))
		return
	}

	// TODO Update the user token
	if err := s.users.UpdateToken(user, session.Key); err != nil {
		w.Write([]byte("Failed to update user token"))
		return
	}

	SetCookie(w, s.config.Cookie, session)
	w.Write([]byte("valid!"))
}

func (s *Server) ListenAndServe() error {
	// Build the server address
	return http.ListenAndServe(s.config.Domain, nil)
}

func NewServer(c Config, u UserManager, s SessionManager) *Server {
	srv := &Server{
		config:   c,
		users:    u,
		sessions: s,
	}

	// Attach the handlers
	http.HandleFunc("/", srv.LoginRequired(srv.Secrets))
	http.HandleFunc("/gen/", srv.Generate)
	http.HandleFunc("/auth/", srv.Authenticate)
	return srv
}
