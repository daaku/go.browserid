// Package browserid provides a way to have a shared identifier for an
// incoming request and allows for it to persist via cookies.
package browserid

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"github.com/nshah/go.domain"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const failId = "deadbeef0000000000000000deadbeef"

var (
	cookieName = flag.String(
		"browserid.cookie", "z", "Name of the cookie to store the ID.")
	maxAge = flag.Duration(
		"browserid.max-age", time.Hour*24*365*10, "Max age of the cookie.")
	idLen = flag.Uint(
		"browserid.len", 16, "Number of bytes to use for ID.")
)

// Check if a ID has been set.
func Has(r *http.Request) bool {
	cookie, err := r.Cookie(*cookieName)
	return err == nil && cookie != nil && isGood(cookie.Value)
}

// Get the ID, creating one if necessary.
func Get(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(*cookieName)
	if err != nil && err != http.ErrNoCookie {
		log.Printf("Error reading browserid cookie: %s", err)
	}
	if cookie != nil && isGood(cookie.Value) {
		return cookie.Value
	}
	id, err := genID()
	if err != nil {
		log.Printf("Error generating browserid: %s", err)
		return failId
	}
	http.SetCookie(w, &http.Cookie{
		Name:    *cookieName,
		Value:   id,
		Path:    "/",
		Expires: time.Now().Add(*maxAge),
		Domain:  cookieDomain(r.Host),
	})
	return id
}

// Returns an empty string on failure to skip explicit domain.
func cookieDomain(host string) string {
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err != nil {
			log.Printf("Error parsing host: %s", host)
			return ""
		}
		host = h
	}
	if host == "localhost" {
		return ""
	}
	if net.ParseIP(host) != nil {
		return ""
	}
	registered, err := domain.Registered(host)
	if err != nil {
		log.Printf("Error extracting base domain: %s", err)
		return ""
	}
	return "." + registered
}

func genID() (string, error) {
	i := make([]byte, *idLen)
	_, err := rand.Read(i)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(i), nil
}

func isGood(value string) bool {
	switch value {
	case "":
		return false
	case failId:
		return false
	}
	return uint(len(value)) == *idLen
}
