/*
Copyright (c) 2024 Diagrid Inc.
Licensed under the MIT License.
*/

package etcdcron

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dapr/kit/cron"
	"github.com/dapr/kit/events/queue"
	"github.com/dapr/kit/ptr"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.etcd.io/etcd/client/pkg/v3/logutil"
	etcdclient "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/diagridio/go-etcd-cron/counting"
	"github.com/diagridio/go-etcd-cron/partitioning"
	"github.com/diagridio/go-etcd-cron/storage"
)

const (
	defaultEtcdEndpoint = "127.0.0.1:2379"
	defaultNamespace    = "etcd_cron"
)

// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type Cron struct {
	namespace   string
	triggerFunc TriggerFunction
	running     atomic.Bool
	readyCh     chan struct{}
	log         logr.Logger
	wg          sync.WaitGroup

	client       *etcdclient.Client
	jobStore     storage.JobStore
	organizer    partitioning.Organizer
	partitioning partitioning.Partitioner

	queue *queue.Processor[string, *Entry]
}

type Options struct {
	Log          logr.Logger
	Namespace    *string
	Partitioning partitioning.Partitioner
	Client       *etcdclient.Client
	JobStore     storage.JobStore

	TriggerFunc TriggerFunction
}

// New returns a new Cron job runner.
func New(opts Options) (*Cron, error) {
	ns := defaultNamespace
	if opts.Namespace != nil {
		ns = *opts.Namespace
	}

	client := opts.Client
	if client == nil {
		var err error
		client, err = etcdclient.New(etcdclient.Config{
			Endpoints: []string{defaultEtcdEndpoint},
		})
		if err != nil {
			return nil, err
		}
	}

	part := opts.Partitioning
	if part == nil {
		part = partitioning.NoPartitioning()
	}

	organizer := partitioning.NewOrganizer(ns, part)

	c := new(Cron)
	jobStore := opts.JobStore
	if jobStore == nil {
		jobStore = storage.NewEtcdJobStore(storage.Options{
			Client:       client,
			Organizer:    organizer,
			Partitioning: part,
			PutCallback: func(ctx context.Context, jobName string, r *storage.JobRecord) error {
				return c.scheduleJob(ctx, jobFromJobRecord(jobName, r))
			},
			DeleteCallback: func(ctx context.Context, jobName string) error {
				// Best effort to delete the counter (if present)
				// Jobs deleted by expiration will delete their counter first.
				// Jobs deleted manually need this logic.
				partitionId := c.partitioning.CalculatePartitionId(jobName)
				counterKey := c.organizer.CounterPath(partitionId, jobName)
				counter := counting.NewEtcdCounter(c.client, counterKey, 0)
				err := counter.Delete(ctx)
				c.errorsHandler(err, &Job{Name: jobName})
				// Ignore error as it is a best effort.
				return nil
			},
		})
	}

	log := opts.Log
	if log.GetSink() == nil {
		sink, err := logutil.CreateDefaultZapLogger(zap.InfoLevel)
		if err != nil {
			return nil, err
		}
		log = zapr.NewLogger(sink)
		log = log.WithName("diagrid-cron")
	}

	c.namespace = ns
	c.triggerFunc = opts.TriggerFunc
	c.readyCh = make(chan struct{})
	c.log = log
	c.client = client
	c.jobStore = jobStore
	c.organizer = organizer
	c.partitioning = part
	return c, nil
}

// Start the cron scheduler in its own go-routine.
func (c *Cron) Start(ctx context.Context) error {
	if !c.running.CompareAndSwap(false, true) {
		return errors.New("already running")
	}
	defer c.wg.Wait()

	if err := c.jobStore.Start(ctx); err != nil {
		return err
	}

	handleFunc := func(e *Entry) {
		if e.job.Status == JobStatusActive && e.job.expired(time.Now()) {
			c.errorsHandler(c.jobStore.Delete(ctx, e.job.Name), e.job)
			return
		}

		if e.job.Status == JobStatusDeleted {
			c.errorsHandler(c.jobStore.Delete(ctx, e.job.Name), e.job)
			return
		}

		result, err := c.triggerFunc(ctx, TriggerRequest{
			JobName:  e.job.Name,
			Metadata: e.job.Metadata,
			Payload:  e.job.Payload,
		})
		if err != nil {
			c.errorsHandler(err, e.job)
			return
		}

		if result == Delete {
			c.errorsHandler(c.jobStore.Delete(ctx, e.job.Name), e.job)
			return
		}

		if result == OK && e.counter != nil {
			// TODO: @joshvanl: garbage collect delete counter
			value := e.counter.Value()
			if value <= 1 {
				c.errorsHandler(c.jobStore.Delete(ctx, e.job.Name), e.job)
				return
			}

			// Needs to check number of triggers
			_, _, err = e.counter.Increment(ctx, -1)
			c.errorsHandler(err, e.job)
			// No need to abort if updating the count failed.
			// The count solution is not transactional anyway.

			c.errorsHandler(c.queue.Enqueue(e), e.job)
		}
	}

	c.queue = queue.NewProcessor[string, *Entry](func(e *Entry) {
		c.wg.Add(1)
		go func() {
			handleFunc(e)
			c.wg.Done()
		}()
	})

	close(c.readyCh)
	<-ctx.Done()

	return c.queue.Close()
}

func (c *Cron) AddJob(ctx context.Context, job *Job) error {
	select {
	case <-c.readyCh:
	case <-ctx.Done():
		return ctx.Err()
	}
	return c.jobStore.Put(ctx, job.Name, job.toJobRecord())
}

func (c *Cron) DeleteJob(ctx context.Context, jobName string) error {
	select {
	case <-c.readyCh:
	case <-ctx.Done():
		return ctx.Err()
	}
	return c.jobStore.Delete(ctx, jobName)
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (c *Cron) scheduleJob(ctx context.Context, job *Job) error {
	select {
	case <-c.readyCh:
	case <-ctx.Done():
		return ctx.Err()
	}

	var sched cron.Schedule
	var next *time.Time
	if job.Rhythm != nil {
		var err error
		sched, err = cron.ParseStandard(*job.Rhythm)
		if err != nil {
			return fmt.Errorf("failed to parse schedule: %s", err)
		}
	} else {
		next = ptr.Of(*job.DueTime)
	}

	partitionId := c.partitioning.CalculatePartitionId(job.Name)
	if !c.partitioning.CheckPartitionLeader(partitionId) {
		// It means the partitioning changed and persisted jobs are in the wrong partition now.
		return fmt.Errorf("host does not own partition %d", partitionId)
	}

	var counter counting.Counter
	if job.Repeats > 0 {
		counterKey := c.organizer.CounterPath(partitionId, job.Name)
		// Needs to count the number of invocations.
		counter = counting.NewEtcdCounter(c.client, counterKey, int(job.Repeats))
	}

	c.queue.Enqueue(&Entry{
		schedule: sched,
		next:     next,
		job:      job,
		counter:  counter,
	})

	return nil
}

// GetJob retrieves a job by name.
func (c *Cron) GetJob(ctx context.Context, jobName string) (*Job, error) {
	select {
	case <-c.readyCh:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	// TODO: @joshvanl
	return nil, nil
}

func (c *Cron) errorsHandler(err error, job *Job) {
	if err == nil {
		return
	}
	c.log.Error(err, "error handling job", "job", job.Name)
}
