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
[GIN-debug] POST   /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).create-fm (3 handlers)
[GIN-debug] GET    /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).readAll-fm (3 handlers)
[GIN-debug] GET    /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).read-fm (3 handlers)
[GIN-debug] PUT    /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).update-fm (3 handlers)
[GIN-debug] PUT    /api/kv/:b/:k/:kid/weight --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).updateWeight-fm (3 handlers)
[GIN-debug] DELETE /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).delete-fm (3 handlers)
[GIN-debug] DELETE /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).deleteAll-fm (3 handlers)
[GIN-debug] GET    /api/search               --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).search-fm (3 handlers)

```

- Search

``` 
curl http://localhost:9999/api/search?b=snippets&k=code&key=v.name&value=linux&shortBy=weight,desc&offset=10&limit=2
```
- Example
``` 
curl --location --request POST 'http://localhost:9999/api/kv/snippets/code' \
--header 'Content-Type: application/json' \
--data-raw '{"name":"linux io","code":"iostat -x -t 2"}'
```
