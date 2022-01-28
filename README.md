# README

A RESTFul API implemented by Golang using graphql

## Overview
Here are the implementations of four different versions of the RESTFul API, corresponding to the four branches:

- [graphql](https://github.com/surzia/restful-api)
- [gin](https://github.com/surzia/restful-api/tree/gin)
- [gorilla/mux](https://github.com/surzia/restful-api/tree/gorilla/mux)
- [std-lib](https://github.com/surzia/restful-api/tree/std-lib)

## Get Start
It uses [gqlgen](https://github.com/99designs/gqlgen) to generate the GraphQL model and bindings.

The GraphQL schema is in `graph/schema.graphqls`; when the schema changes, rerun the tool to regenerate the code:

```shell
go run github.com/99designs/gqlgen generate
```

To run the server:
```shell
go run main.go
```

Then visit the printed link in a [browser](http://localhost:8880) for the GraphQL playground.

## Reference
- [Juejin](https://juejin.cn/post/7056018992976101389)
