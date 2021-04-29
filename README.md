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
[GIN-debug] POST   /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).create-fm (3 handlers)
[GIN-debug] GET    /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).readAll-fm (3 handlers)
[GIN-debug] GET    /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).read-fm (3 handlers)
[GIN-debug] PUT    /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).update-fm (3 handlers)
[GIN-debug] PUT    /api/kv/:b/:k/:kid/weight --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).updateWeight-fm (3 handlers)
[GIN-debug] DELETE /api/kv/:b/:k/:kid        --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).delete-fm (3 handlers)
[GIN-debug] DELETE /api/kv/:b/:k             --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).deleteAll-fm (3 handlers)
[GIN-debug] GET    /api/search               --> github.com/yusys-cloud/go-jsonstore-rest/internal.(*Storage).search-fm (3 handlers)


```

- Search

``` 
curl localhost:9999/api/search?b=snippets&k=code&key=v.name&value=linux&shortBy=weight,desc
```