package user

import (
	"context"
	"database/sql"
	"net/http"
	"reflect"

	v "github.com/core-go/core/v10"
	"github.com/core-go/search/query"
	q "github.com/core-go/sql"

	"go-service/internal/user/handler"
	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

type UserTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(db *sql.DB, logError func(context.Context, string, ...map[string]interface{})) (UserTransport, error) {
	validator, err := v.NewValidator()
	if err != nil {
		return nil, err
	}

	userType := reflect.TypeOf(model.User{})
	queryBuilder := query.NewBuilder(db, "users", userType)
	userSearchBuilder, err := q.NewSearchBuilder(db, userType, queryBuilder.BuildQuery)
	if err != nil {
		return nil, err
	}
	// userRepository, err := adapter.NewUserAdapter(db, adapter.BuildQuery)
	userRepository, err := q.NewRepository(db, "users", userType)
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userSearchBuilder.Search, userService, logError, validator.Validate, nil)
	return userHandler, nil
}
