package app

import (
	"github.com/nekruz08/BookStorage/pkg/core/Books"
	"github.com/nekruz08/jwt/pkg/jwt"
	"github.com/nekruz08/mux/pkg/mux"
	"github.com/nekruz08/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	router   *mux.ExactMux
	BooksSvc *Books.Service
	secret   jwt.Secret
}

func NewServer(router *mux.ExactMux, productSvc *Books.Service, secret jwt.Secret) *Server {
	return &Server{router: router, BooksSvc: productSvc, secret: secret}
}

func (s Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s Server) Start() {
	s.InitRoutes()
}

func (s Server) handleBooksList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.BooksSvc.BooksList(request.Context())
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &list)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleDeleteBook() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = s.BooksSvc.RemoveByID(request.Context(), int64(id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func (s Server) handBook() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		context, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(context)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		prod := Books.Product{}
		err = rest.ReadJSONBody(request, &prod)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if id == 0 {
			err = s.BooksSvc.AddNewBook(request.Context(), prod)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		if id > 0 {
			err = s.BooksSvc.UpdateProduct(request.Context(), int64(id), prod)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

	}
}

func (s Server) handleBooksByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prod, err := s.BooksSvc.BookByID(request.Context(), int64(id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &prod)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}
