/*
Copyright (c) 2024 Diagrid Inc.
Licensed under the MIT License.
*/

package elector

import (
	"context"
	"testing"

	"github.com/diagridio/go-etcd-cron/internal/key"
	"github.com/diagridio/go-etcd-cron/tests/framework/etcd"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/types/known/anypb"
)

func Test_attemptNewLeadership(t *testing.T) {
	t.Parallel()

	t.Run("non-existing key should write leadership key", func(t *testing.T) {
		t.Parallel()

		client := etcd.Embedded(t)
		key, err := key.New(key.Options{
			Namespace: "abc",
			ID:        "helloworld",
		})
		require.NoError(t, err)

		replicaData := &anypb.Any{Value: []byte("hello")}

		lease, err := client.Grant(context.Background(), 20)
		require.NoError(t, err)

		e := New(Options{
			Log:         logr.Discard(),
			Client:      client,
			Key:         key,
			ReplicaData: replicaData,
			LeaseID:     lease.ID,
		})

		resp, err := client.Get(context.Background(), "abc/leadership", clientv3.WithPrefix())
		require.NoError(t, err)

		ok, err := e.attemptNewLeadership(context.Background(), resp)
		require.NoError(t, err)
		assert.True(t, ok)

		resp, err = client.Get(context.Background(), "abc/leadership", clientv3.WithPrefix())
		require.NoError(t, err)
		assert.Len(t, resp.Kvs, 1)

		// TODO: @joshvanl: assert kv value to be leadership with correct replica
		// data and uid.
	})
}
