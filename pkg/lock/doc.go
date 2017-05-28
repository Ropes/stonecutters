// Package lock handles atomic requests for associative identifiers in etcd.
//
// Request an atomic ID from etcdv3:
//
//    IDs := []string{"coffee", "tea", "bikes"}
//
//    ctx, cancl := context.WithCancel()
//    defer cancl()
//
//    // Responsibility of keeping lease alive is up to caller
//    lease, err :=  etcdclient.Grant(ctx, int64(5))
//
//    id, err := st.Join(etcdclient, ctx, leaseID.ID, "hostname", IDs)
//    ...
package lock
