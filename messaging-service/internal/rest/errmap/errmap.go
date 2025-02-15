package errmap

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
)

type Response struct {
	Code int
	Body restapi.ErrorResponse
}

func Map(err error) Response {
	var domErr domain.Error
	if errors.As(err, &domErr) {
		if resp, ok := domainErrMap[domErr]; ok {
			return resp
		}
		panic(fmt.Errorf("domain error is not mapped: %s", domErr))
	}

	var servicesErr services.Error
	if errors.As(err, &servicesErr) {
		if resp, ok := servicesErrMap[servicesErr]; ok {
			return resp
		}
		panic(fmt.Errorf("services error is not mapped: %s", servicesErr))
	}

	// Every non-domain error is counted as internal server error.
	return Response{
		Code: http.StatusInternalServerError,
		Body: restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeInternal,
			ErrorMessage: "Internal Server Error",
		},
	}
}

var servicesErrMap = map[services.Error]Response{
	// TODO: fill it
}

var domainErrMap = map[domain.Error]Response{
	// TODO: fill it
}
