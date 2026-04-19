package wstypes

import (
	"context"
	"encoding/json"
)

type Request struct {
	Method  string          `json:"method"`
	Payload json.RawMessage `json:"payload"`
}

type Metadata struct {
	Ctx    context.Context
	UserID int64
}

type Method func(metadata *Metadata, cw *ConnWriter, req *Request)
