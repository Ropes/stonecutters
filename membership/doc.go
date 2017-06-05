// Package membership handles atomic requests for associative identifiers in etcd.

// This package looks deep within a goroutine's soul and assigns it a name based on the order in which it joined.
//
// Designed to address the container namespace pollution of metric collection systems. Carbon https://github.com/graphite-project/carbon in particular, where namespace pollution kills Graphite's performance.)
//
// Distributed processes which need to share names but for identification and maintain uniqueness. This works but having a shared static set of names which are claimed using etcd(v3) as the distributed lock. Each process iterates over the ordered list, and claim the first name which isn't regestered/claimed in etcd.
//
// Provided static names are the top 100 highest mountains in North America, ordered by decending peak elevation. However any list of identifiers can be passed into membership.Join(...)
//
//
// Request an atomic ID from etcdv3:
//
//    IDs := []string{"coffee", "tea", "bikes"}
//    ctx, cancl := context.WithCancel()
//
//    // Responsibility of keeping lease alive is up to caller
//    lease, err :=  etcdclient.Grant(ctx, int64(5))
//    ...
//
//    // Join the stonecutters IDs list
//    member, err := membership.Join(etcdclient, ctx, leaseID.ID, "homer", IDs)
//    ...
//
//    // List all members
//    members, err := membership.Members(etcdclient, IDs)
//    ...
package membership
