package socketio_test

import (
	"net/url"
	"testing"
)

func TestParse(t *testing.T) {
	u, err := url.Parse("http://username:password@host:8080/directory/file?query#ref")
	if err != nil || u.String() != "http://username:password@host:8080/directory/file?query#ref" {
		t.Error("Invalid url parse")
	}
}

func TestParseRelativePath(t *testing.T) {
	u, err := url.Parse("https://woot.com/test")

	if err != nil || u.Scheme != "https" || u.Host != "woot.com" || u.Path != "/test" {
		t.Error("Invalid url parse")
	}
}

// TODO - others port not work with go `net/url package`
