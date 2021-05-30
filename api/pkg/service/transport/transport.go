package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/manufacturer/api/pkg/service"
	"github.com/manufacturer/api/pkg/v1"
)

// Post Do Logic With Manufacturer .............................................

type PostDoLogicWithManufacturerRequest struct {
	input []*v1.ManufacturerInput
}

type PostDoLogicWithManufacturerResponse struct{}

func EncodePostDoLogicWithManufacturerResponse(_ context.Context, _ http.ResponseWriter, _ interface{}) error {
	return nil
}

func DecodePostDoLogicWithManufacturerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := PostDoLogicWithManufacturerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req.input); err != nil {
		return nil, v1.ErrBadRequest("failed to decode JSON request: %v", err)
	}

	return req, nil
}

func MakePostDoLogicWithManufacturerEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostDoLogicWithManufacturerRequest)
		err := svc.PostDoLogicWithManufacturer(ctx, req.input)

		return PostDoLogicWithManufacturerResponse{}, err
	}
}
