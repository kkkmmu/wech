package db_test

import (
	. "db"
	"testing"
)

func TestNewDBr(t *testing.T) {
	db, err := NewFileDB("test")
	if err != nil {
		t.Error("Cannot create DB")
	}

	err = db.Save("Hello.txt", []byte("Hello world"))
	if err != nil {
		t.Error("Cannot write DB")
	}

	v, err := db.Get("Hello.txt")
	if err != nil {
		t.Error("Cannot Read content")
	}

	if string(v) != "Hello world" {
		t.Error("Value Changed after store into DB")
	}

	if err = db.Destroy(); err != nil {
		t.Error("Cannot Destroy DB")
	}
}
