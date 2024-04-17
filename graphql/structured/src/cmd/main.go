package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/handler"

	"github.com/tkhrk1010/go-samples/graphql/structured/src/gql"
)

func main() {
	schema, err := gql.BuildSchema()
	if err != nil {
		log.Fatal(err)
	}

	// handlerは、定義されたスキーマを元にGraphQLのリクエストをresolverに渡す
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	log.Fatal(http.ListenAndServe(":8080", nil))
}