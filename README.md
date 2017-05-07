[Stonecutters](https://youtu.be/HmEtR17A6ck?t=2m55s)
------------

This package looks deep within a goroutine's soul and assigns it a name based on the order in which it joined.

For distributed systems; assigning UIDs to running processes is a common way to identify processes and their metrics. However older metric designes like Carbon/Graphite don't handle unique naming which polute and cause large wildcard search paths. Thus recycling names as prefixes can be used to reduce namespace polution by uids, but still keep processes effectively unique. 

Distributed processes which need to share names but for identification and maintain uniqueness. This works but having a shared static set of names which are claimed using etcd(v3) as the distributed lock. Each process iterates over the ordered list, and claim the first name which isn't regestered/claimed in etcd.

Current static names are the top 100 highest mountains in North America, ordered by decending peak elevation.
eg:
foo-ns31 -> foo-Denali
foo-3le9 -> foo-MtLogan
...
foo-wedc -> foo-MtRainier



