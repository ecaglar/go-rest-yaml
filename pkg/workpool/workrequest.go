package workpool

import (
	"../model"
	"github.com/google/uuid"
)

//WorkRequest defines the work that can be processed by workers.
type WorkRequest struct {
	ID      uuid.UUID
	Payload model.Metadata
}
