package http

import (
	"net/http"
	"time"

	kv "github.com/adoublef/yap/nats"
	"github.com/nats-io/nats.go"
)

type CookieJar struct {
	ttl time.Duration
	kv  nats.KeyValue
}

func (jar *CookieJar) Get(w http.ResponseWriter, r *http.Request, dst any) (err error) {
	var (
		name = jar.kv.Bucket()
	)
	// TODO, will use casting just like the context package
	c, err := r.Cookie(name)
	if err != nil {
		return err
	}
	_, err = kv.Get(jar.kv, c.Value, dst)
	return
}
func (jar *CookieJar) Has(w http.ResponseWriter, r *http.Request, dst any) (ok bool) {
	return jar.Get(w, r, dst) == nil
}
func (jar *CookieJar) Put(w http.ResponseWriter, r *http.Request, key string, src any) (err error) {
	var (
		name   = jar.kv.Bucket()
		maxAge = int(jar.ttl.Seconds())
	)
	_, err = kv.Put(jar.kv, key, src)
	if err != nil {
		return err
	}
	// edit the ttl
	c := &http.Cookie{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// TODO make this clearer
		Name:   name,
		Value:  key,
		MaxAge: maxAge,
	}
	http.SetCookie(w, c)
	return
}

func NewCookieJar(nc *kv.Conn, name string, ttl time.Duration) (jar *CookieJar, err error) {
	jar = &CookieJar{
		ttl: ttl,
	}
	jsc, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	jar.kv, err = kv.UpsertKV(jsc, &nats.KeyValueConfig{
		Bucket:  name,
		Storage: nats.MemoryStorage,
		TTL:     ttl,
	})
	return
}
