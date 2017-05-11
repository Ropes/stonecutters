package lock

import (
	"context"
	"fmt"

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
func acquireLease(c *clientv3.Client, ctx context.Context, timeout int64) (clientv3.LeaseID, error) {
	resp, err := c.Grant(ctx, timeout)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func revokeLease(client *clientv3.Client, ctx context.Context, lease clientv3.LeaseID) error {
	_, err := client.Revoke(ctx, lease)
	return err
}

var PutError = fmt.Errorf("error putting key:value pair")

// kvPutOrGet writes a key-val pair given that the key is not already in use.
// If the key exists it is returned, if it does not exist they key-val is Put.
// TxnResponse, error is returned.
func kvPutOrGet(kvc clientv3.KV, ctx context.Context, key, val string) (*clientv3.TxnResponse, error) {
	resp, err := kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), ">", 0)).
		Then(clientv3.OpGet(key)).
		Else(clientv3.OpPut(key, val)).
		Commit()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// kvPutOrGet writes a key-val pair with a lease given that the key is not already in use.
// If the key exists it is returned, if it does not exist they key-val is Put.
// TxnResponse, error is returned.
func kvPutLeaseOrGet(kvc clientv3.KV, ctx context.Context, lease clientv3.LeaseID, key, val string) (*clientv3.TxnResponse, error) {
	resp, err := kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), ">", 0)).
		Then(clientv3.OpGet(key)).
		Else(clientv3.OpPut(key, val, clientv3.WithLease(lease))).
		Commit()
	if err != nil {
		return nil, err
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

// verifyTxnResponse recieves a TxnResponse and validates that ek,ev key-value
// pair are written in etcd.
func verifyTxnResponse(resp clientv3.TxnResponse, ek, ev string) bool {
	responses := resp.Responses
	if len(responses) == 1 {
		resp := responses[0].GetResponse()
		fmt.Printf("---- %T %#v\n", resp, resp)
		switch T := resp.(type) {
		case *etcdserverpb.ResponseOp_ResponsePut:
			kv := T.ResponsePut
			fmt.Printf("---- %T %#v\n", kv, kv)

			/*
				if string(kv.Key) == ek && string(kv.Value) == ev {
					return true
				}
			*/
		case *etcdserverpb.ResponseOp_ResponseRange:
			rr := T.ResponseRange
			fmt.Printf("---- %T %#v\n", rr, rr)

		default:
			return false
		}
	}
	return false
}

// client, context, kv, lease
func claimName(c clientv3.Client, ctx context.Context, id string) (bool, error) {
	return true, nil
}
