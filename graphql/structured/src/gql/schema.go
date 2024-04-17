// src/pkg/graphql/schema.go
package gql

import (
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// 
// GraphQLのスキーマを構築する
func BuildSchema() (graphql.Schema, error) {
	// スキーマ定義を読み込む
	schemaBytes, err := os.ReadFile("src/gql/schema.graphql")
	if err != nil {
			return graphql.Schema{}, err
	}

	schemaString := string(schemaBytes)
	// GraphQLのスキーマ定義(.graphql)をASTに変換する
	// ASTとは、GraphQLのスキーマを表現するためのデータ構造
	astDoc, err := parser.Parse(parser.ParseParams{
			Source: &source.Source{Body: []byte(schemaString)},
	})
	if err != nil {
		return graphql.Schema{}, err
	}

	// getQueryTypeを使って、ASTからQueryの定義を取得
	// NewSchema関数を使って、Queryの定義からSchemaを生成
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    getQueryType(astDoc),
		Mutation: getMutationType(astDoc),
	})
	if err != nil {
		return graphql.Schema{}, err
	}

	return schema, nil
}

// AST (.graphqlをパースした結果) からQueryの定義を取得
func getQueryType(doc *ast.Document) *graphql.Object {
	var queryType *ast.ObjectDefinition
	for _, def := range doc.Definitions {
		if def, ok := def.(*ast.ObjectDefinition); ok && def.Name.Value == "Query" {
			queryType = def
			break
		}
	}

	fields := graphql.FieldsThunk(func() graphql.Fields {
    result := graphql.Fields{}
    if queryType != nil {
        for _, field := range queryType.Fields {
            result[field.Name.Value] = &graphql.Field{
                Type:    getGraphQLType(field.Type),
								// ここでは、getResolverが、usersを受け取ったらresolveUsersを返すようにしている
                Resolve: getResolver(field.Name.Value),
            }
        }
    }
    return result
	})

	return graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
	})
}

func getMutationType(doc *ast.Document) *graphql.Object {
	// Mutationの実装は省略
	return nil
}

func getGraphQLType(astType ast.Type) graphql.Type {
	switch t := astType.(type) {
	case *ast.Named:
			switch t.Name.Value {
			case "ID":
					return graphql.ID
			case "String":
					return graphql.String
			default:
					return graphql.NewObject(graphql.ObjectConfig{
							Name: t.Name.Value,
							Fields: graphql.Fields{
									"id": &graphql.Field{
											Type: graphql.NewNonNull(graphql.ID),
									},
									"name": &graphql.Field{
											Type: graphql.NewNonNull(graphql.String),
									},
							},
					})
			}
	case *ast.NonNull:
			return graphql.NewNonNull(getGraphQLType(t.Type))
	case *ast.List:
			return graphql.NewList(getGraphQLType(t.Type))
	default:
			panic("Unknown AST type")
	}

	// ここに到達することはないが、compilerがエラーを出さないようにnilを返す記述を追加
	return nil
}

func getResolver(fieldName string) graphql.FieldResolveFn {
	switch fieldName {
	case "users":
			return resolveUsers
	default:
			return nil
	}
}