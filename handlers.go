package eauth

import (
	"fmt"
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

func (s *Server) LoginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session key from the cookie
		cookie, err := r.Cookie(s.config.Cookie.Name)

		// Check if the session is valid
		if err != nil || !IsValidSession(s.sessions, cookie.Value) {
			w.Write([]byte("Login Required"))
			return
		}

		// Execute the wrapped handler
		h(w, r)
	}
}

// func (s *Server) UserIs(test UserTest, h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Get the session key from the cookie
// 		cookie, err := r.Cookie(s.config.Cookie.Name)

// 		// TODO Get the user for the given session
// 	}
// }

var emailForm = []byte(`<!DOCTYPE html>
<html>
	<body>
		<form method="POST">
			<input type="email" name="email" placeholder="Enter your email">
			<input type="submit">
		</form>
	</body>
</html>`)

var emailBody = []byte(``)

// Send email will send a auth link to the given email address if the given
// POST data contains a valid user email.
func (s *Server) SendEmail(w http.ResponseWriter, r *http.Request) {
	// If
	if strings.ToUpper(r.Method) == "POST" {
		email := r.FormValue("email")
		w.Write([]byte(email))
		w.Write([]byte("\n"))

		// TODO Email normalization
		user := s.users.GetEmail(email)
		if user.Id == 0 {
			w.Write([]byte("Invalid user"))
			return
		}

		w.Write([]byte(fmt.Sprint(user)))
		w.Write([]byte("\n"))

		// TODO email this link!
		link := CreateLink(s.config, s.users.Get(1))
		w.Write([]byte(link))
		return
	}

	// TODO How to make the templates extensible?
	w.Write(emailForm)
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
	http.HandleFunc("/login/", srv.SendEmail)
	http.HandleFunc("/auth/", srv.Authenticate)
	return srv
}
