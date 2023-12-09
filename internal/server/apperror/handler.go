package apperror

import (
	"net/http"
)

type errorHandler func(w http.ResponseWriter, r *http.Request) error

func Handler(h errorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "core/json")
		//
		//var appErr *AppError
		//err := h(w, r)
		//
		//if err != nil {
		//	isAppError := errors.As(err, &appErr)
		//
		//	if isAppError {
		//		switch appErr {
		//		case NotFound:
		//			w.WriteHeader(http.StatusNotFound)
		//		case ErrUserAlreadyExists:
		//			w.WriteHeader(http.StatusConflict)
		//		default:
		//			w.WriteHeader(http.StatusBadRequest)
		//		}
		//	} else {
		//		appErr = systemError(err)
		//	}
		//
		//	w.Write(appErr.Marshal())
		//}
	}
}
