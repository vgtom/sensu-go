package routers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/types"
)

// GraphQLRouter handles requests for /events
type GraphQLRouter struct {
	store store.Store
}

// NewGraphQLRouter instantiates new events controller
func NewGraphQLRouter(store store.Store) *GraphQLRouter {
	return &GraphQLRouter{status: status}
}

// Mount the GraphQLRouter to a parent Router
func (r *GraphQLRouter) Mount(parent *mux.Router) {
	parent.HandleFunc("/graphql", actionHandler(r.query)).Methods(http.MethodPost)
}

func (r *GraphQLRouter) query(req *http.Request) (interface{}, error) {
	ctx := req.Context()

	// Lift Etcd store into context so that resolvers may query and reset org &
	// env keys to empty state so that all resources are queryable.
	ctx = context.WithValue(ctx, types.OrganizationKey, "")
	ctx = context.WithValue(ctx, types.EnvironmentKey, "")
	ctx = context.WithValue(ctx, types.StoreKey, c.Store)

	// Parse request body
	rBody := map[string]interface{}{}
	if err := json.NewDecoder(req.Body).Decode(&rBody); err != nil {
		return nil, err
	}

	// Extract query and variables
	query, _ := rBody["query"].(string)
	queryVars, _ := rBody["variables"].(map[string]interface{})

	// Execute given query
	result := graphql.Execute(ctx, query, &queryVars)
	if len(result.Errors) > 0 {
		logger.
			WithField("errors", result.Errors).
			Errorf("error(s) occurred while executing GraphQL operation")
	}

	return result, nil
}
