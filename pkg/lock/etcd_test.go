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
)

func init() {
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
		t.Run("deleteKeys", deleteAll)
		t.Run("txnStatic", txnStaticKey)
	})
	//client.Close()
	//e.Close()
}

func deleteAll(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dresp, err := client.Delete(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		t.Errorf("error deleting all keys: %v", err)
	}
	t.Logf("Deleted keys: %#v", dresp.Deleted)
}

func txnStaticKey(t *testing.T) {
	key := "Denali"
	val := "wyeast"

	ctx, cancel := context.WithCancel(context.Background())
	kvc := clientv3.NewKV(client)

	// Create txn to write key if it does not exist
	resp, err := kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), ">", 0)).
		Then(clientv3.OpGet(key)).
		Else(clientv3.OpPut(key, val)).
		Commit()
	if err != nil {
		t.Errorf("error executing txn: %v", err)
	}
	t.Logf("write txn resp: %#v", resp)

	// Get the key; test its value
	got, err := client.Get(ctx, key)
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
	t.Logf("failed write resp: %#v", resp)

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
	cancel()
}
