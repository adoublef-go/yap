package nats

import (
	"errors"

	"github.com/nats-io/nats.go"
)

func UpsertKV(js nats.JetStreamContext, c *nats.KeyValueConfig) (kv nats.KeyValue, err error) {
	if c == nil || c.Bucket == "" {
		return nil, errors.New("invalid config")
	}
	kv, err = js.KeyValue(c.Bucket)
	switch {
	case errors.Is(err, nats.ErrBucketNotFound):
		kv, err = js.CreateKeyValue(c)
	case err != nil:
		return nil, err
	}
	return
}

func Put(kv nats.KeyValue, name string, value any) (revision int, err error) {
	b, err := Encode(value)
	if err != nil {
		return 0, err
	}
	n, err := kv.Put(name, b)
	return int(n), err
}

func Get(kv nats.KeyValue, name string, dst any) (revision int, err error) {
	entry, err := kv.Get(name)
	if err != nil {
		return 0, err
	}
	err = Decode(entry.Value(), dst)
	revision = int(entry.Revision())
	return revision, err
}
