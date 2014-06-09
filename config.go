package eauth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

// TODO A mechanism should be in place to rotate secret keys
// TODO Multiple domains?
type Config struct {
	Domain    string                    `json:"domain"`
	Secret    string                    `json:"secret"`
	Https     bool                      `json:"https"`
	Cookie    CookieConfig              `json:"cookie"`
	SMTP      map[string]SMTPConfig     `json:"smtp"`
	Databases map[string]DatabaseConfig `json:"databases"`
}

var localhost = Config{
	Domain: "localhost:8008",
	Secret: "yabbadabbado", // TODO this should NEVER be used
	Cookie: defaultCookie,
}

// Cookie names are valid tokens as defined by RFC 2616 section 2.2:
// http://tools.ietf.org/html/rfc2616#section-2.2
// TL;DR: Any non-control or non-separator character.
type CookieConfig struct {
	Age      time.Duration `json:"age"`
	Domain   string        `json:"domain"`
	HttpOnly bool          `json:"http_only"`
	Name     string        `json:"name"`
	Path     string        `json:"path"`
	Secure   bool          `json:"secure"`
}

// The default cookie implementation - not very secure
var defaultCookie = CookieConfig{
	Age:      14 * 24 * time.Hour, // Two weeks
	Domain:   "",
	HttpOnly: true,
	Name:     "eauthid",
	Path:     "/",
	Secure:   false,
}

type SMTPConfig struct {
	Port     int64  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	From     string `json:"from"`
	Alias    string `json:"alias"`
}

func (c SMTPConfig) FromAddress() string {
	if c.Alias != "" {
		return fmt.Sprintf(`"%s" <%s>`, c.Alias, c.From)
	}
	return fmt.Sprintf("<%s>", c.From)
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Return a string of credentials approriate for Go's sql.Open() func
func (db DatabaseConfig) Credentials() string {
	// TODO Different credentials for different drivers
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s",
		db.Host,
		db.Port,
		db.Name,
		db.User,
		db.Password,
	)
}

// By default, the parser will look for a file called settings.json in
// current directory.
func Parse() (Config, error) {
	return ParseFile("./settings.json")
}

func ParseFile(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	return parse(f)
}

func parse(f io.Reader) (Config, error) {
	var c Config
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return c, err
	}
	if err = json.Unmarshal(contents, &c); err != nil {
		return c, err
	}
	// TODO Allow flag values to override configuration?

	// Fall back to the default cookie if none was set
	if c.Cookie.Name == "" {
		c.Cookie = defaultCookie
	}
	return c, nil
}
