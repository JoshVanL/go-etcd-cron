/*
Copyright (c) 2024 Diagrid Inc.
Licensed under the MIT License.
*/

package etcdcron

import (
	"context"
	"time"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

type TriggerRequest struct {
	JobName   string
	Timestamp time.Time
	Actual    time.Time
	Metadata  map[string]string
	Payload   *anypb.Any
}

type TriggerResult int

const (
	OK TriggerResult = iota
	Failure
	Delete
)

type TriggerFunction func(ctx context.Context, req TriggerRequest) (TriggerResult, error)
