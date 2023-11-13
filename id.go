package yap

import "github.com/gofrs/uuid"

func UUID() uuid.UUID {
	return uuid.Must(uuid.NewGen().NewV7())
}
