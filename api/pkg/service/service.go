package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/manufacturer/api/pkg/v1"
)

type Service interface {
	PostDoLogicWithManufacturer(ctx context.Context, input []*v1.ManufacturerInput) (err error)
}

type service struct {
}

func (s *service) PostDoLogicWithManufacturer(ctx context.Context, input []*v1.ManufacturerInput) (err error) {
	if input == nil {
		err = v1.ErrInternal("bad input from transport")

		return
	}

	for _, in := range input {
		if in.Details == nil || !in.Details.NeedUpdate {
			if err == nil {
				err = fmt.Errorf("no details in manufacturer %s or needUpdate is false", in.ID)
			} else {
				err = errors.Wrapf(err, "no details in manufacturer %s or needUpdate is false", in.ID)
			}
		}
	}

	if err != nil {
		err = v1.ErrBadRequest(err.Error())
	}

	return
}

func NewService() Service {
	return &service{}
}
