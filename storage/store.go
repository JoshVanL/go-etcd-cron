/*
Copyright (c) 2024 Diagrid Inc.
Licensed under the MIT License.
*/

package storage

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/mirror"
	"google.golang.org/protobuf/proto"

	"github.com/diagridio/go-etcd-cron/partitioning"
)

// The JobStore persists and reads jobs from Etcd.
type JobStore interface {
	Start(ctx context.Context) error
	Put(ctx context.Context, jobName string, job *JobRecord) error
	Delete(ctx context.Context, jobName string) error
	Wait()
}

type etcdStore struct {
	runWaitingGroup sync.WaitGroup
	etcdClient      *etcdclient.Client
	kvStore         etcdclient.KV
	partitioning    partitioning.Partitioner
	organizer       partitioning.Organizer
	putCallback     func(context.Context, string, *JobRecord) error
	deleteCallback  func(context.Context, string) error
}

type Options struct {
	Client         *etcdclient.Client
	Organizer      partitioning.Organizer
	Partitioning   partitioning.Partitioner
	PutCallback    func(context.Context, string, *JobRecord) error
	DeleteCallback func(context.Context, string) error
}

func NewEtcdJobStore(opts Options) JobStore {
	return &etcdStore{
		etcdClient:     opts.Client,
		kvStore:        etcdclient.NewKV(opts.Client),
		partitioning:   opts.Partitioning,
		organizer:      opts.Organizer,
		putCallback:    opts.PutCallback,
		deleteCallback: opts.DeleteCallback,
	}
}

func (s *etcdStore) Start(ctx context.Context) error {
	for _, partitionId := range s.partitioning.ListPartitions() {
		// TODO(artursouza): parallelize this per partition.
		partitionPrefix := s.organizer.JobsPath(partitionId) + "/"
		partitionSyncer := mirror.NewSyncer(s.etcdClient, partitionPrefix, 0)
		rc, errc := partitionSyncer.SyncBase(ctx)

		for r := range rc {
			for _, kv := range r.Kvs {
				err := s.notifyPut(ctx, kv, s.putCallback)
				if err != nil {
					return err
				}
			}
		}

		err := <-errc
		if err != nil {
			return err
		}

		s.sync(ctx, partitionPrefix, partitionSyncer)
	}

	return nil
}

func (s *etcdStore) Put(ctx context.Context, jobName string, job *JobRecord) error {
	bytes, err := proto.Marshal(job)
	if err != nil {
		return err
	}
	_, err = s.kvStore.Put(
		ctx,
		s.organizer.JobPath(jobName),
		string(bytes),
	)
	return err
}

func (s *etcdStore) Delete(ctx context.Context, jobName string) error {
	_, err := s.kvStore.Delete(
		ctx,
		s.organizer.JobPath(jobName))
	return err
}

func (s *etcdStore) Wait() {
	s.runWaitingGroup.Wait()
}

func (s *etcdStore) notifyPut(ctx context.Context, kv *mvccpb.KeyValue, callback func(context.Context, string, *JobRecord) error) error {
	_, jobName := filepath.Split(string(kv.Key))
	record := JobRecord{}
	err := proto.Unmarshal(kv.Value, &record)
	if err != nil {
		return fmt.Errorf("could not unmarshal job for key %s: %v", string(kv.Key), err)
	}
	if jobName == "" || (record.Rhythm == nil && record.DueTimestamp == nil) {
		return fmt.Errorf("could not deserialize job for key %s", string(kv.Key))
	}

	return callback(ctx, jobName, &record)
}

func (s *etcdStore) notifyDelete(ctx context.Context, name string, callback func(context.Context, string) error) error {
	return callback(ctx, name)
}

func (s *etcdStore) sync(ctx context.Context, prefix string, syncer mirror.Syncer) {
	s.runWaitingGroup.Add(1)
	go func() {
		log.Printf("Started sync for path: %s\n", prefix)
		wc := syncer.SyncUpdates(ctx)
		for {
			select {
			case <-ctx.Done():
				s.runWaitingGroup.Done()
				return
			case wr := <-wc:
				for _, ev := range wr.Events {
					t := ev.Type
					switch t {
					case mvccpb.PUT:
						if err := s.notifyPut(ctx, ev.Kv, s.putCallback); err != nil {
							log.Printf("Error notifying put: %v", err)
						}
					case mvccpb.DELETE:
						_, name := filepath.Split(string(ev.Kv.Key))
						if err := s.notifyDelete(ctx, name, s.deleteCallback); err != nil {
							log.Printf("Error notifying delete: %v", err)
						}
					default:
						log.Printf("Unknown etcd event type: %v", t.String())
					}
				}
			}
		}
	}()
}
