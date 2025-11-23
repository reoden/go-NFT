package pipeline

import (
	"context"

	"github.com/reoden/go-echo-template/pkg/logger"
	"github.com/reoden/go-echo-template/pkg/validation"

	"github.com/mehdihadeli/go-mediatr"
)

type mediatorValidationPipeline struct {
	logger logger.Logger
}

func NewMediatorValidationPipeline(l logger.Logger) mediatr.PipelineBehavior {
	return &mediatorValidationPipeline{logger: l}
}

func (m mediatorValidationPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	v, ok := request.(validation.Validator)
	if ok {
		err := v.Validate()
		if err != nil {
			return nil, err
		}
	}

	return next(ctx)
}
