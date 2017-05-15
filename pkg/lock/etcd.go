package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
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

// Lease Functionality
func acquireLease(c *clientv3.Client, ctx context.Context, timeout int64) (*clientv3.LeaseGrantResponse, error) {
	resp, err := c.Grant(ctx, timeout)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func revokeLease(client *clientv3.Client, ctx context.Context, leaseID clientv3.LeaseID) error {
	_, err := client.Revoke(ctx, leaseID)
	return err
}

var PutError = fmt.Errorf("lock: error putting key-value pair")

// kvPutLease writes a key-val pair with a lease given that the key is not already in use.
// If the key exists the Txn fails, if it does not exist they key-val is Put.
func kvPutLease(kvc clientv3.KV, ctx context.Context, leaseID clientv3.LeaseID, key, val string) (*clientv3.TxnResponse, error) {
	resp, err := kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, val, clientv3.WithLease(leaseID))).
		Commit()
	if err != nil {
		return nil, err
	}
	if resp.Succeeded == false {
		return nil, errors.New(fmt.Sprintf("key %q already registered", key))
	}
	return resp, nil
}

// respSingleKv takes a TxnResponse and returns the key-value pair
// assuming there is only one returned by the ResponseRange.
func respSingleKv(tr *clientv3.TxnResponse) (string, string) {
	if len(tr.Responses) == 1 {
		rop := tr.Responses[0].GetResponseRange()
		if len(rop.Kvs) == 1 {
			Kvs := rop.Kvs[0]
			return string(Kvs.Key), string(Kvs.Value)
		}
	}
	return "", ""
}

// verifyKvPair returns true if expected key-value strings match their expected values
func verifyKvPair(client *clientv3.Client, ek, ev string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	got, err := client.Get(ctx, ek)
	if err != nil {
		return false
	}
	if string(got.Kvs[0].Value) == ev {
		return true
	}
	return false
}

// deleteme
// verifyTxnResponse recieves a TxnResponse and validates that ek,ev key-value
// pair are written in etcd.
func verifyTxnResponse(client *clientv3.Client, resp clientv3.TxnResponse, ek, ev string) bool {
	responses := resp.Responses
	if len(responses) == 1 {
		resp := responses[0].GetResponse()
		switch resp.(type) {
		case *etcdserverpb.ResponseOp_ResponsePut:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			got, err := client.Get(ctx, ek)
			if err != nil {
				return false
			}
			if string(got.Kvs[0].Value) == ev {
				return true
			}
			return true //Txn Succeded, and Value matches
		case *etcdserverpb.ResponseOp_ResponseRange:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			got, err := client.Get(ctx, ek)
			if err != nil {
				//fmt.Printf("failed to GET key %s\n%v", ek, got.Kvs)
				return false
			}
			if string(got.Kvs[0].Value) == ev {
				return true
			}
		default:
			fmt.Printf("default response hit; failure\n")
			return false
		}
	}
	return false
}

// client, context, kv, lease
func claimName(c clientv3.Client, ctx context.Context, id string) (bool, error) {
	return true, nil
}
