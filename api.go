/*
Copyright (c) 2024 Diagrid Inc.
Licensed under the MIT License.
*/

package etcdcron

import (
	"context"
	"errors"
	"fmt"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/diagridio/go-etcd-cron/api"
)

// ListResponse is the response of listing jobs with a given prefix.
type ListResponse struct {
	// Jobs is the list of jobs in the given prefix.
	Jobs []*api.Job

	// More indicates there are more jobs in the given prefix which do not fit in
	// the response size.
	More bool
}

// Add adds a new cron job to the cron instance.
func (c *cron) Add(ctx context.Context, name string, job *api.Job) error {
	select {
	case <-c.readyCh:
	case <-c.closeCh:
		return errors.New("cron is closed")
	case <-ctx.Done():
		return ctx.Err()
	}

	if err := c.validateName(name); err != nil {
		return err
	}

	if job == nil {
		return errors.New("job cannot be nil")
	}

	storedJob, err := c.schedBuilder.Parse(job)
	if err != nil {
		return fmt.Errorf("job failed validation: %w", err)
	}

	b, err := proto.Marshal(storedJob)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	_, err = c.client.Put(ctx, c.key.JobKey(name), string(b))
	return err
}

// Get gets a cron job from the cron instance.
func (c *cron) Get(ctx context.Context, name string) (*api.Job, error) {
	select {
	case <-c.readyCh:
	case <-c.closeCh:
		return nil, errors.New("cron is closed")
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := c.validateName(name); err != nil {
		return nil, err
	}

	resp, err := c.client.Get(ctx, c.key.JobKey(name))
	if err != nil {
		return nil, err
	}

	// No entry is not an error, but a nil object.
	if resp.Count == 0 {
		return nil, nil
	}

	var stored api.JobStored
	if err := proto.Unmarshal(resp.Kvs[0].Value, &stored); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return stored.GetJob(), nil
}

// Delete deletes a cron job from the cron instance.
func (c *cron) Delete(ctx context.Context, name string) error {
	select {
	case <-c.readyCh:
	case <-c.closeCh:
		return errors.New("cron is closed")
	case <-ctx.Done():
		return ctx.Err()
	}

	if err := c.validateName(name); err != nil {
		return err
	}

	if _, err := c.client.Delete(ctx, c.key.JobKey(name)); err != nil {
		return err
	}

	return c.queue.Dequeue(c.key.JobKey(name))
}

func (c *cron) List(ctx context.Context, prefix string) (*ListResponse, error) {
	select {
	case <-c.readyCh:
	case <-c.closeCh:
		return nil, errors.New("cron is closed")
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	resp, err := c.client.Get(ctx,
		c.key.JobKey(prefix),
		clientv3.WithPrefix(),
	)
	if err != nil {
		return nil, err
	}

	jobs := make([]*api.Job, 0, resp.Count)
	for _, kv := range resp.Kvs {
		var stored api.JobStored
		if err := proto.Unmarshal(kv.Value, &stored); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job from prefix %q: %w", prefix, err)
		}

		jobs = append(jobs, stored.GetJob())
	}

	return &ListResponse{
		Jobs: jobs,
		More: resp.More,
	}, nil
}

// validateName validates the name of a job.
func (c *cron) validateName(name string) error {
	if len(name) == 0 {
		return errors.New("job name cannot be empty")
	}

	for _, segment := range strings.Split(strings.ToLower(c.validateNameReplacer.Replace(name)), "||") {
		if errs := validation.IsDNS1123Subdomain(segment); len(errs) > 0 {
			return fmt.Errorf("job name is invalid %q: %s", name, strings.Join(errs, ", "))
		}
	}

	return nil
}
