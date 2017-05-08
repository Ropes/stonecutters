package lock

import (
	"context"

	"github.com/coreos/etcd/clientv3"
)

// TODO: Abstract to interfaces
// for now; implement etcd as locking/lease mechanism for claiming names

// GetID iterates over the passed 'ids' and attempts to claim one in
// etcd with a Lease which is persisted until the context is closed.
// If the list of ids are all claimed, the function will pause for 5
// seconds before iterating over all the ids again, if it fails to lock
// it returns an error.
func GetID(c clientv3.Client, ctx context.Context, ids []string) (string, error) {
	return "", nil
}

func acquireLease(c clientv3.Client, ctx context.Context, timeout int64) (clientv3.LeaseID, error) {
	resp, err := c.Grant(ctx, timeout)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

// client, context, kv, lease
func claimName(c clientv3.Client, ctx context.Context, id string) (bool, error) {
	return true, nil
}
