package main

import (
	"testing"
	"strings"

	"github.com/rafaeljusto/redigomock"
)

func TestEncodeDecode(t *testing.T) {
	id := int64(42424242429)
	short := encodeURL(id)
	decodedID, err := decodeURL(short)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if id != decodedID {
		t.Errorf("encoding/Decoding does not work, id: %v, decodedID: %v", id,
			decodedID)
	}
}

func TestSingleUrl(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Clear()
	// data
	url := "https://www.google.com/search?q=golang+spec"
	parsed, err := ParseURL(url)
	if err != nil {
		t.Errorf("Error during parsing valid url: %s\n", err)
	}
	// saving & shortening
	id := int64(9910050000)
	conn.Command("GET", IdKey).Expect(id)
	conn.Command("INCR", IdKey).Expect(id + 1)
	conn.Command("SET", id, url).Expect("ok")

	shorten, err := parsed.SaveURL(conn)
	if err != nil {
		t.Errorf("error shortening url: %s\n", err)
	}
	// fetching
	conn.Clear()
	conn.Command("GET", id).Expect(url)
	decoded, err := DbGetURL(shorten, conn)
	if err != nil {
		t.Errorf("error fetching url: %s\n", err)
	}
	if strings.Compare(decoded, url) != 0 {
		t.Errorf("urls does not match: %s and %s\n", url, decoded)
	}
}

func TestMultipleUrls(t *testing.T) {

}

func TestDuplicateUrl(t *testing.T) {

}

func TestInvalidUrl(t *testing.T) {

}
