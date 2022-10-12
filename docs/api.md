API docs

query with curl

```shell
$ curl -G -s http://localhost:8080/api/v1/query -d 'query=pod{namespace="cert-manager",pod="cert-manager-.*"}' | jq | head -4
{
  "data": {
    "result": [
      "2022-10-12T00:30:39Z[cert-manager|cert-manager-cainjector-6995cf7d4-582tw|cert-manager] I1012 00:30:39.116889...",
      ...
```


query_range with curl

```shell
$ curl -G -s http://localhost:8080/api/v1/query_range -d 'query=pod{namespace="cert-manager",pod="cert-manager-.*"}' -d 'start=1665549365.661' -d 'end=1665549665.661' | jq | head -4
{
  "data": {
    "result": [
      "2022-10-12T04:36:06Z[cert-manager|cert-manager-cainjector-6995cf7d4-582tw|cert-manager] I1012 04:36:06.709410...",
      ...
```
