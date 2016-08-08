package main

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kyokomi/bouillabaisse/firebase"
	"gopkg.in/yaml.v2"
)

// TODO: SQLiteとかにしたほうがよいかも?

type AuthStore struct {
	firebase.Auth
	UpdateAt  time.Time
	CreatedAt time.Time
}

func (a *AuthStore) ExpiresInText() string {
	expiredTime, isExpired := a.ExpiredTime()
	if isExpired || expiredTime.Before(time.Now()) {
		return "期限切れ" // TODO: text
	}
	return expiredTime.String()
}

func (a *AuthStore) ExpiredTime() (time.Time, bool) {
	expiresIn, err := strconv.Atoi(a.ExpiresIn)
	if err != nil {
		return time.Time{}, true
	}
	return a.UpdateAt.Add(time.Duration(expiresIn) * time.Second), false
}

func (a AuthStore) Save(dirPath string, fileName string) error {
	data, err := yaml.Marshal(&a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dirPath, fileName), data, 0666)
}

func (a *AuthStore) Load(dirPath string, fileName string) error {
	data, err := ioutil.ReadFile(filepath.Join(dirPath, fileName))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &a)
}
