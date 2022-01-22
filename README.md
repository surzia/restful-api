# README

A RESTFul API implemented by Golang using standard library

##Overview
```
POST   /page/              :  create a page, returns ID
GET    /page/<pageid>      :  returns a single page by ID
GET    /page/              :  returns all pages
DELETE /page/<pageid>      :  delete a page by ID
PUT    /page/<pageid>      :  update a page by ID
GET    /tag/<tagname>      :  returns list of pages with this tag
GET    /due/<yy>/<mm>/<dd> :  returns list of pages due by this date
```
## Get Start

```shell
go run main.go
```
server runs in [localhost:8880](http://localhost:8880)
## Reference
- [Juejin](https://juejin.cn/post/7052931619962748958)
