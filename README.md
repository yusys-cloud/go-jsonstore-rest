# Go-Jsonstore-Rest

A quick and easy way to setup a RESTful JSON DATASTORE server for small projects.

## Quick Start

Binary Download

https://github.com/yusys-cloud/go-jsonstore-rest/releases 

Startup

``` 
./go-jsonstore-rest
```

Replace ./json-db with the path to the drive or directory in which you want to store data.

```
./go-jsonstore-rest -path=./json-db -port=9999 
```

### RESTful API

- CRUD

``` 
// Create
curl localhost:9999/api/kv/meta/node -X POST -d '{"ip": "192.168.49.69","name":"redis-n1","dc":"default","lable":"Redis"}' --header "Content-Type: application/json"
// Read
curl localhost:9999/api/kv/meta/node
// Update
curl localhost:9999/api/kv/meta/node/node:1429991523109310464 -X PUT -d '{"ip": "192.168.49.69","name":"redis-n2","dc":"default","lable":"Redis"}' --header "Content-Type: application/json"
// Delete
curl localhost:9999/api/kv/meta/node/node:1429991523109310464 -X DELETE 
// Search
curl http://localhost:9999/api/search?b=snippets&k=code&key=v.name&value=linux&shortBy=weight,desc&offset=10&limit=2
```

``` 
[GIN-debug] POST   /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).create-fm (3 handlers)
[GIN-debug] GET    /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).readAll-fm (3 handlers)
[GIN-debug] GET    /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).read-fm (3 handlers)
[GIN-debug] PUT    /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).update-fm (3 handlers)
[GIN-debug] PUT    /api/kv/:b/:k/:kid/weight --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).updateWeight-fm (3 handlers)
[GIN-debug] DELETE /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).delete-fm (3 handlers)
[GIN-debug] DELETE /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).deleteAll-fm (3 handlers)
[GIN-debug] GET    /api/search               --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).search-fm (3 handlers)

```
<img width="550px" src="./docs/static/crud-time-ms.jpg">

## Benchmarks
- BenchmarkQuery100-12    	       1	2,102,479,566 ns/op
- BenchmarkCreate100-12    	       1	1082093859 ns/op
