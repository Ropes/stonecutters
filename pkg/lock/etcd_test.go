package lock

import (
	"context"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/embed"
)

var (
	client *clientv3.Client
	e      *embed.Etcd
	err    error

	key string
	val string
)

func init() {
	key = "Denali"
	val = "wyeast"
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	cfg.ForceNewCluster = true
	e, err = embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	select {
	case <-e.Server.ReadyNotify():
		log.Infof("Server is ready!")
	case <-time.After(6 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Infof("Server took too long to start!")
	}
	go func() {
		log.Fatal(<-e.Err())
	}()

	ccfg := &clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	client, err = clientv3.New(*ccfg)
	if err != nil {
		log.Fatalf("error creating etcd client: %v", err)
	}

	time.Sleep(1 * time.Second)
}

func TestEtcd(t *testing.T) {
	t.Run("etcd tests", func(t *testing.T) {
		t.Run("deleteKeys", deleteKey)
		t.Run("txnStatic", txnStaticKey)
	})
	//client.Close()
	//e.Close()
}

func deleteKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dresp, err := client.Delete(ctx, key)
	if err != nil {
		t.Errorf("error deleting all keys: %v", err)
	}
	t.Logf("Deleted keys: %#v", dresp.Deleted)
}

func txnStaticKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kvc := clientv3.NewKV(client)

	// Get the key; test its value
	got, err := client.Get(ctx, key)
	if err != nil {
		t.Error(err)
	}
	if len(got.Kvs) > 0 {
		t.Errorf("data already logged from last run %s: %s", key, string(got.Kvs[0].Value))
	}

	resp, err := kvPutOrGet(kvc, ctx, key, val)
	if err != nil || resp == nil {
		t.Fatalf("error executing txn: %v", err)
	}
	t.Logf("first response: %#v", resp)

	var verified bool
	/*
		verified = verifyTxnResponse(resp, key, val) // This panics?
		if !verified {
			t.Errorf("kv verification failed")
		}
	*/

	// Get the key; test its value
	got, err = client.Get(ctx, key)
	if err != nil {
		t.Error(err)
	}
	rv := string(got.Kvs[0].Value)
	if rv != val {
		t.Errorf("value not set to %q", val)
	}
	t.Logf("%s: %s", key, rv)

	// Test comparision to overwrite; does not change the value
	resp, err = kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), ">", 0)).
		Then(clientv3.OpGet(key)).
		Else(clientv3.OpPut(key, "yyvest")).
		Commit()
	if err != nil {
		t.Errorf("error executing txn: %v", err)
	}
	verified = verifyTxnResponse(resp, key, val)
	if verified != true {
		t.Errorf("verification post if-already-exists failed")
	}

	// Get the key; test its value
	got, err = client.Get(ctx, key)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s: %s", key, string(got.Kvs[0].Value))
	if string(got.Kvs[0].Value) != val {
		t.Errorf("%s was overwritten to: %s", key, string(got.Kvs[0].Value))
	}
	t.Logf("%#v\n", string(got.Kvs[0].Value))
}
