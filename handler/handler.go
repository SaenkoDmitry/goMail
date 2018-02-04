package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"goMail/utils"
	"goMail/dbs/mysql"
	"go.uber.org/zap"
	"time"
	"errors"
	"fmt"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.HandleToken(r.Header.Get("Authorization"))
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusUnauthorized)
			utils.Logger.Error("cannot handle token",
				zap.Error(err),
			)
			return
		}
		utils.Mu.Lock()
		a, exists := utils.Cookies[token]
		utils.Mu.Unlock()
		if !exists {
			utils.WriteResponse(w, errors.New("cannot find token in cookies"), nil, http.StatusUnauthorized)
			utils.Logger.Error("cannot find token in cookies",
				zap.Error(errors.New("cannot find token in cookies")),
			)
			return
		}
		_, err = mysql.GetUser(a)
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
			utils.Logger.Error("cannot get user because of mysql error",
				zap.Error(err),
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Logger.Info("within the method",
			zap.String("url", r.URL.Path),
			zap.String("method", r.Method),
			zap.Duration("backoff", time.Second),
		)
		next.ServeHTTP(w, r)
		utils.Logger.Info("out of the method",
			zap.String("url", r.URL.Path),
			zap.String("method", r.Method),
			zap.Duration("backoff", time.Second),
		)
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.Logger.Error("Panic middleware: recovered",
					zap.String("url", r.URL.Path),
					zap.String("method", r.Method),
				)
				utils.WriteResponse(w, errors.New(fmt.Sprintf("%v", err)), nil, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func InitHandlers(router *mux.Router) {
	router.Handle("/users",
		accessLogMiddleware(panicMiddleware(http.HandlerFunc(getAllUsers))),
	).Methods("GET")

	router.Handle("/spaces",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(getAllSpaces)))),
	).Methods("GET")

	router.Handle("/history",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(getUserHistory)))),
	).Methods("GET")

	router.Handle("/space/{name_space}/history",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(getSpaceHistory)))),
	).Methods("GET")

	router.Handle("/spaces/{name_space}/permissions",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(getAllSpacePermissions)))),
	).Methods("GET")

	router.Handle("/spaces/{name_space}/tuples",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(getAllTuples)))),
	).Methods("GET")

	router.Handle("/spaces/{name_space}/tuples/{id_tuple}",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(getTuple)))),
	).Methods("GET")

	router.Handle("/users",
		accessLogMiddleware(panicMiddleware(http.HandlerFunc(loginUser))),
	).Methods("POST")

	router.Handle("/spaces/{name_space}",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(addSpace)))),
	).Methods("POST")

	router.Handle("/spaces/{name_space}/tuples/{id_tuple}",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(addTuple)))),
	).Methods("POST")

	router.Handle("/users/{name}/spaces/{name_space}/permissions",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(addPermission)))),
	).Methods("POST")

	router.Handle("/spaces/{name_space}/tuples/{id_tuple}",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(updateTuple)))),
	).Methods("PUT")

	router.Handle("/spaces/{name_space}/tuples/{id_tuple}",
		authMiddleware(accessLogMiddleware(panicMiddleware(http.HandlerFunc(deleteTuple)))),
	).Methods("DELETE")

}
