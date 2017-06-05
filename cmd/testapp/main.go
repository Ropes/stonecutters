package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	"github.com/ropes/stonecutters"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	name := flag.String("name", "default", "stonecutter member name")
	etcdUrl := flag.String("etcdaddr", "localhost:2379", "etcd connection address")
	flag.Parse()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	// Create etcd client
	ccfg := &clientv3.Config{
		Endpoints:   []string{*etcdUrl},
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(*ccfg)
	if err != nil {
		log.Fatalf("error creating etcd client: %v", err)
	}
	// Create etcd Lease
	lease, err := client.Grant(ctx, int64(10))
	if err != nil {
		os.Exit(1)
	}

	// Create etcd lease keepalive
	_, kaerr := client.KeepAlive(ctx, lease.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}

	IDs := stonecutters.PrefixedNumerics("/metrics/testapp", 100)

	// Request an ID from stonecutters Join
	ID, err := stonecutters.Join(client, ctx, lease.ID, *name, IDs)
	if err != nil {
		log.Fatalf("error joining stonecutters: %v", err)
		os.Exit(1)
	}

	// for{ print ID granted, list other members }
	for {
		select {
		case <-s:
			client.Revoke(ctx, lease.ID)
			client.Close()
			cancel()
			os.Exit(0)
		default:
			log.WithFields(log.Fields{"name": *name, "ID": ID}).Info("Member")
			members, err := stonecutters.Members(client, IDs)
			if err != nil {
				log.Fatalf("error listing members: %v", err)
				os.Exit(1)
			}
			log.WithFields(log.Fields{"count": len(members)}).Infof("members: %#v", members)
			time.Sleep(time.Second * 3)
		}
	}

}
