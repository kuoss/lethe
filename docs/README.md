## letheql

Examples
```
pod{namespace="ingress-nginx"}
pod{namespace="ingress-nginx",container="controller"}
pod{namespace=~"(kube-system|ingress-nginx)",container="controller"}

pod{namespace="kube-system"}
pod{namespace="kube-system"} |= "hello"
pod{namespace="kube-system"} |= "hello" != "world"
pod{namespace="kube-system"} |~ "err|ERR" != "Liveness"
```


Label filter operators:
* `=`: exactly equal
* `!=`: not equal
* `=~`: regex matches
* `!~`: regex does not match

Line filter operators:
* `|=`: Line contains string
* `!=`: Line does not contain string
* `|~`: Line contains a match to the regex
* `!~`: Line does not contain a match to the regex


## HTTP API

query with curl
```shell
$ curl -G -s http://localhost:8080/api/v1/query \
-d 'query=pod{namespace="cert-manager",pod=~"cert-manager-.*"}' | jq
```

query_range with curl
```shell
$ curl -G -s http://localhost:8080/api/v1/query_range \
-d 'query=pod{namespace="cert-manager",pod=~"cert-manager-.*"}' \
-d 'start=1665549365.661' -d 'end=1665549665.661' | jq
```
