[Stonecutters](https://youtu.be/HmEtR17A6ck?t=2m55s)
------------

**This package is still in active development; APIs may change so lock to a semantic version**

![](https://vignette1.wikia.nocookie.net/simpsons/images/1/16/Hammer_symbol.png/revision/latest?cb=20101006090032)

This package looks deep within a goroutine's soul and assigns it a name based on the order in which it joined.

Designed to address the container namespace pollution of metric collection systems.([Carbon](https://github.com/graphite-project/carbon) in particular, where namespace pollution kills Graphite's performance.)

For distributed systems; assigning UIDs to running processes is a common way to identify processes and their metrics. However older metric designes like Carbon/Graphite don't handle unique naming which polute and cause large wildcard search paths. Thus recycling names as prefixes can be used to reduce namespace pollution by uids, but still keep processes effectively unique. 

Distributed processes which need to share names but for identification and maintain uniqueness. This works but having a shared static set of names which are claimed using etcd(v3) as the distributed lock. Each process iterates over the ordered list, and claim the first name which isn't regestered/claimed in etcd.

Provided static names are the top 100 highest mountains in North America, ordered by decending peak elevation. However any list of identifiers can be passed into stonecutters.Join(...)  

eg owners/hosts named:

```
foo-ns31 -> Denali
foo-3le9 -> MtLogan
...
foo-wedc -> MtRainier
```

### etcd stonecutters Transaction to assign an Identifier 
1. if key does **not** exist 
2. PUT the key:value{Identifier: Owner} pair using a lease
  * If key is used, iterate to next identifier and retry claim transaction

If all identifiers are claimed return error.

## API

```
IDs := []string{"coffee", "tea", "bikes"}
ctx, cancl := context.WithCancel()

// Responsibility of keeping lease alive is up to caller
lease, err :=  etcdclient.Grant(ctx, int64(5))
...

// Join the stonecutters IDs list
member, err := stonecutters.Join(etcdclient, ctx, leaseID.ID, "homer", IDs)
...

// List all members
members, err := stonecutters.Members(etcdclient, IDs)
...
```

## Testing

Since etcd is critical to the stonecutters, tests are all effectively integration tests.

Recomended strategy is to run `etcd` in standalone along with the tests.Downloading [etcd](https://github.com/coreos/etcd/releases/tag/v3.2.10) and then run with eg:`./etcd-v3.2.10-linux-amd64/etcd`


Embeded `etcd` can be configured to start and run with the tests by setting `ETCDEMBED=1` in the test environment. This make starting tests take a while though.

All tests attempt to clean up and revoke all keys after finishing as to not polute etcd between runs.


## Unvendored Dependendies

List created with `dep`, but not Gopkg.toml not included to avoid unresolvable dependency trees.

```
## Manually curated list to describe minimal dependency versions
[[constraint]]
  name = "github.com/coreos/etcd"
  version = "3.1.7"

[[constraint]]
  name = "golang.org/x/text"
  version = "ccbd3f7"

[[constraint]]
  name = "github.com/Sirupsen/logrus"
  version = "cdd90c38c6e3718c731b555b9c3ed1becebec3ba"
```
