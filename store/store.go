package store

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
)

// Stores local auth TODO: SQLiteとかにしたほうがよいかも?
var Stores AuthStores

// AuthStores multi authStore
type AuthStores struct {
	Data map[string]AuthStore
}

// AuthStore firebase authを保存する単位
type AuthStore struct {
	firebase.Auth
	UpdateAt  time.Time
	CreatedAt time.Time
}

// NewAuthStore create AuthStore
func NewAuthStore(a firebase.Auth) AuthStore {
	nowTime := time.Now()
	return AuthStore{
		Auth:      a,
		UpdateAt:  nowTime,
		CreatedAt: nowTime,
	}
}

// ExpiresInText return expiresTime string
func (a *AuthStore) ExpiresInText() string {
	expiredTime, isExpired := a.ExpiredTime()
	if isExpired || expiredTime.Before(time.Now()) {
		return "期限切れ" // TODO: text
	}
	return expiredTime.Format("2006-01-02 15:04:05")
}

// ExpiredTime return expiresTime
func (a *AuthStore) ExpiredTime() (time.Time, bool) {
	expiresIn, err := strconv.Atoi(a.ExpiresIn)
	if err != nil {
		return time.Time{}, true
	}
	return a.UpdateAt.Add(time.Duration(expiresIn) * time.Second), false
}

// Add added auth store
func (a *AuthStores) Add(store AuthStore) {
	if a.Data == nil {
		a.Data = map[string]AuthStore{}
	}
	a.Data[store.LocalID] = store
}

// Remove removed auth store
func (a *AuthStores) Remove(localID string) {
	delete(a.Data, localID)
}

// Save Saved AuthStores in local file.
func (a *AuthStores) Save(dirPath string, fileName string) error {
	data, err := yaml.Marshal(&a.Data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dirPath, fileName), data, 0666)
}

// Load Loading AuthStores at local file
func (a *AuthStores) Load(dirPath string, fileName string) error {
	data, err := ioutil.ReadFile(filepath.Join(dirPath, fileName))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &a.Data)
}
