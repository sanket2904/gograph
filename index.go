package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Number struct {
	Number int `graphql:"number:int"`
}

type Address struct {
	Street string `graphql:"street:string"`
	City   string `graphql:"city:string"`
	Number Number `graphql:"number:Number"`
}

type Customer struct {
	Name    string  `graphql:"name:string"`
	Address Address `graphql:"address:Address"`
}

func returnGraphqlType(obj string) graphql.Type {
	switch obj {
	case "string":
		return graphql.String
	case "int":
		return graphql.Int
	case "float":
		return graphql.Float
	case "bool":
		return graphql.Boolean
	case "id":
		return graphql.ID
	default:
		return graphql.String
	}
}

func graphobj[T any](obj T) *graphql.Object {
	// print all the fields of the struct with their types
	//  looping through the fields of the struct and creating a graphql field for each
	fields := graphql.Fields{}
	for i := 0; i < reflect.TypeOf(obj).NumField(); i++ {
		field := reflect.TypeOf(obj).Field(i)
		tag := field.Tag.Get("graphql")
		fmt.Println(field.Name, field.Type, tag)
		// use recursion if the field is another struct
		if field.Type.Kind() == reflect.Struct {
			fields[field.Name] = &graphql.Field{
				Type: graphobj(reflect.New(field.Type).Elem().Interface()),
			}
		} else if tag != "" {
			name := strings.Split(tag, ":")[0]
			gtype := strings.Split(tag, ":")[1]
			fields[name] = &graphql.Field{
				Type: returnGraphqlType(gtype),
			}
		}
	}
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   reflect.TypeOf(obj).Name(),
		Fields: fields,
	})
}

func main() {

	customerType := graphobj(Customer{})

	var rootPrivateQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "query",
		Fields: graphql.Fields{
			"customer": &graphql.Field{
				Type: customerType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					_, ok := p.Args["id"].(string)

					if ok {
						return Customer{
							Name: "John",
							Address: Address{
								Street: "123",
								City:   "NY",
							},
						}, nil
					}
					return nil, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootPrivateQuery,
	})

	qraphQlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", qraphQlHandler)
	http.ListenAndServe(":1337", nil)

}
