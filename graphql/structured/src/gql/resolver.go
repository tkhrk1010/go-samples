package gql

import "github.com/graphql-go/graphql"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: "1", Name: "Alice"},
	{ID: "2", Name: "Bob"},
}

func resolveUsers(p graphql.ResolveParams) (interface{}, error) {
	return users, nil
}