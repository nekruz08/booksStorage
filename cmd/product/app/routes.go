package app

import (
	"github.com/nekruz08/BookStorage/pkg/core/token"
	"github.com/nekruz08/BookStorage/pkg/mux/middleware/authenticated"
	"github.com/nekruz08/BookStorage/pkg/mux/middleware/authorized"
	"github.com/nekruz08/BookStorage/pkg/mux/middleware/jwt"
	"github.com/nekruz08/BookStorage/pkg/mux/middleware/logger"
	"reflect"
)

func (s Server) InitRoutes() {


	s.router.GET(
		"/api/products/{id}",
		s.handleBooksByID(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get Books by id"),
	)
	s.router.GET(
		"/api/products",
		s.handleBooksList(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get list"),
	)



	s.router.DELETE(
		"/api/products/{id}",
		s.handleDeleteBook(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("delete Books"),
	)

	s.router.POST(
		"/api/products/{id}",
		s.handBook(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post Books"),
	)

}