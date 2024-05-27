package models_da

import "github.com/google/uuid"

type DocumentData struct {
	ID            uuid.UUID
	DocumentBytes []byte //pdf file -- the whole file
}
