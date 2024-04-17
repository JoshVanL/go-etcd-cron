/*
Copyright (c) 2024 Diagrid Inc.
Licensed under the MIT License.
*/

package etcdcron

import (
	"sync"
	"time"

	"github.com/dapr/kit/cron"
	"github.com/dapr/kit/ptr"
	"github.com/diagridio/go-etcd-cron/counting"
)

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	// The schedule on which this job should be run.
	schedule cron.Schedule

	// The Job o run.
	job *Job

	next *time.Time

	// Counter if has limit on number of triggers
	counter counting.Counter

	lock sync.Mutex
}

func (e *Entry) ScheduledTime() time.Time {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.next == nil {
		e.next = ptr.Of(e.schedule.Next(time.Now()))
	}

	return *e.next
}

func (e *Entry) Key() string {
	return e.job.Name
}
