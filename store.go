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
var stores authStores

type authStores struct {
	stores map[string]authStore
}

type authStore struct {
	firebase.Auth
	UpdateAt  time.Time
	CreatedAt time.Time
}

func (a *authStore) ExpiresInText() string {
	expiredTime, isExpired := a.ExpiredTime()
	if isExpired || expiredTime.Before(time.Now()) {
		return "期限切れ" // TODO: text
	}
	return expiredTime.Format("2006-01-02 15:04:05")
}

func (a *authStore) ExpiredTime() (time.Time, bool) {
	expiresIn, err := strconv.Atoi(a.ExpiresIn)
	if err != nil {
		return time.Time{}, true
	}
	return a.UpdateAt.Add(time.Duration(expiresIn) * time.Second), false
}

func (a *authStores) Add(store authStore) {
	a.stores[store.LocalID] = store
}

func (a *authStores) Remove(localID string) {
	delete(a.stores, localID)
}

func (a *authStores) Save(dirPath string, fileName string) error {
	data, err := yaml.Marshal(&a.stores)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dirPath, fileName), data, 0666)
}

func (a *authStores) Load(dirPath string, fileName string) error {
	data, err := ioutil.ReadFile(filepath.Join(dirPath, fileName))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &a.stores)
}
