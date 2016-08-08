package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"time"
	"os"

	"github.com/dustin/gojson"
	"github.com/kyokomi/bouillabaisse/firebase"
)

// TODO: SQLiteとかにしたほうがよいかも?

type AuthStore struct {
	firebase.Auth
	UpdateAt  time.Time
	CreatedAt time.Time
}

func (a AuthStore) Save(dirPath string, fileName string) error {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(a); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dirPath, fileName), buf.Bytes(), 0666)
}

func (a *AuthStore) Load(dirPath string, fileName string) error {
	f, err := os.Open(filepath.Join(dirPath, fileName))
	if err != nil {
		return err
	}
	return json.NewDecoder(f).Decode(&a)
}
