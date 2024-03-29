# Go-Jsonstore-Rest

This is a plug-and-play JSON file storage REST API service. It offers a simple way to create, read, update, and delete 
JSON data without the need for an additional database. You can interact with this service by sending HTTP requests, and 
it will store the provided JSON data in files and respond with appropriate results based on the requests. This service is 
particularly suitable for small-scale applications or prototyping, enabling you to quickly set up a REST API environment with data storage support.

## Quick Start

Binary Download

https://github.com/yusys-cloud/go-jsonstore-rest/releases 

Startup

``` 
./go-jsonstore-rest
```

Replace ./json-db with the path to the drive or directory in which you want to store data.

```
./go-jsonstore-rest --path=./json-db --port=9999 
./go-jsonstore-rest --basicAuth=admin:admin 
```
## Installation
``` 
go get github.com/yusys-cloud/go-jsonstore-rest
```
## Usage

### Base Codes
```
import (
	"github.com/gin-gonic/gin"
	"github.com/yusys-cloud/go-jsonstore-rest/rest"
	"net/http"
)
func main() {
	r := gin.Default()
	//REST-APIs-json
	s := rest.NewJsonStoreRest("./api-data")
	
	basicAuth := make(map[string]string)
    basicAuth["admin"] = "admin"
	s.BasicAuth = basicAuth
	
	s.DisableCors = true
	s.ConfigHandles(r)
	r.Run(":9991")
}
```

## Features
RESTful JSON DATASTORE
- CRUD
``` 
curl localhost:9999/kv/meta/node -X POST -d '{"ip": "192.168.x.x","name":"redis-n1","idc":"default","lable":"Redis"}' --header "Content-Type: application/json"
curl localhost:9999/kv/meta/node
curl localhost:9999/kv/meta/node/node:1429991523109310464 -X PUT -d '{"ip": "192.168.49.69","name":"redis-n2","idc":"default","lable":"Redis"}' --header "Content-Type: application/json"
curl localhost:9999/kv/meta/node/node:1429991523109310464 -X DELETE
```
- 批量保存,有id字段对应的历史数据则执行更新 ```/kvAPI/batch?b=b&k=k -X POST -d '[{},{}]' ```
- 查询字段不为空 ```/kvAPI/search?b=datasets&k=train&key=v.declaration&relation=isNotEq&value=null```
- 删除全部Delete ```/kv/meta/node -X DELETE```
- First In, First Out (FIFO) ```curl 'localhost:9999/kvAPI/fifo?size=10' -X POST -d '{"k":"key1","v":"hellow world"}'```
- key:value store/cache ```curl localhost:9999/kvAPI/cache -X POST -d '{"k":"k1","v":"hellow world"}'```
- 搜索 
  - 分页 ```/kvAPI/search?b=meta&k=node&key=v.name&value=linux&shortBy=weight,desc&page=1&size=10```
  - 排序 ```/kvAPI/search?b=node&k=node&key=v.name&value=linux&shortBy=v.weight,desc&page=1&size=2```
  - 多字段 ```/kvAPI/search?b=meta&k=node&key=v.name,v.lable&value=zk|zk2,zk&relation=like```
  - 指定要返回的字段范围 &fields=name,dc

## 搜索参数
  /kv/:b/:k 
- 参数 b 为存储的json文件名，类似数据库名称
- 参数 k 为存储文件中json对象数组名，类似数据库中表名称

- 参数 key 为搜索json对象中字段名称
- 参数 value 为搜索json对象字段的搜索值
- 参数 key 为json对象的搜索字段名称,需多层json key的完整路径值，如：v.name
- 参数 relation 为查询关系，默认为equal精准查询，可选like模糊查询
- 参数 value 为json对象字段的搜索值

## RESTful API

``` 
[GIN-debug] POST   /kv/:b/:k                 --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).create-fm (4 handlers)
[GIN-debug] GET    /kv/:b/:k                 --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).readAll-fm (4 handlers)
[GIN-debug] GET    /kv/:b/:k/:kid            --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).read-fm (4 handlers)
[GIN-debug] PUT    /kv/:b/:k/:kid            --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).update-fm (4 handlers)
[GIN-debug] PUT    /kv/:b/:k/:kid/weight     --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).updateWeight-fm (4 handlers)
[GIN-debug] DELETE /kv/:b/:k/:kid            --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).delete-fm (4 handlers)
[GIN-debug] DELETE /kv/:b/:k                 --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).deleteAll-fm (4 handlers)
[GIN-debug] GET    /kvAPI/search               --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).search-fm (4 handlers)
[GIN-debug] DELETE /kvAPI/search               --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).deleteSearch-fm (4 handlers)
[GIN-debug] POST   /kvAPI/cache                --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).cache-fm (4 handlers)
[GIN-debug] GET    /kvAPI/cache/:key           --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).cacheGet-fm (4 handlers)
[GIN-debug] POST   /kvAPI/fifo                 --> github.com/yusys-cloud/go-jsonstore-rest/rest.(*Storage).fifo-fm (4 handlers)

```
<img width="550px" src="./docs/static/json-rest-crud-costtime.jpg">
<img width="550px" src="./docs/static/crud-time-ms.jpg">

## Benchmarks
- BenchmarkQuery100-12    	       1	2,102,479,566 ns/op
- BenchmarkCreate100-12    	       1	1082093859 ns/op
