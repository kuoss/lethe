## letheql

```
pod{namespace="kube-system"}
pod{namespace="kube-system"} |= "hello"
pod{namespace="kube-system"} |= "hello" != "world"
pod{namespace="kube-system"} |~ "err|ERR" != "Liveness"

pod{namespace="ingress-nginx",container="controller",image=~".*ingress-nginx.*"}
pod{namespace=~"(kube-system|ingress-nginx)",container="controller",image=~".*ingress-nginx.*"}
```


Label 필터 연산자
* `=` 같음
* `!=` 같지 않음
* `=~` 정규표현식에 맞음
* `!~` 정규표현식에 맞지 않음

Line 필터 연산자
* `|=` 같음
* `!=` 같지 않음
* `|~` 정규표현식에 맞음
* `!~` 정규표현식에 맞지 않음


## HTTP API

query with curl
```shell
$ curl -G -s http://localhost:8080/api/v1/query \
-d 'query=pod{namespace="cert-manager",pod="cert-manager-.*"}' | jq
{
  "data": {
    "result": [
      "2022-10-12T00:30:39Z[cert-manager|cert-manager-cainjector-6995cf7d4-582tw|cert-manager] I1012 00:30:39.116889...",
      ...
    ]
  }
}
```

query_range with curl
```shell
$ curl -G -s http://localhost:8080/api/v1/query_range \
-d 'query=pod{namespace="cert-manager",pod="cert-manager-.*"}' \
-d 'start=1665549365.661' -d 'end=1665549665.661' | jq
{
  "data": {
    "result": [
      "2022-10-12T04:36:06Z[cert-manager|cert-manager-cainjector-6995cf7d4-582tw|cert-manager] I1012 04:36:06.709410...",
      ...
    ]
  }
}
```
