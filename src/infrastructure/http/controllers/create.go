package controllers

import (
	"email-masking-svc/src/application/services"
	"email-masking-svc/src/infrastructure/configuration"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type CreateEmailMaskController struct {
	logger     configuration.Logger
	createMask *services.CreateMask
}

func NewCreateEmailMaskController(createMask *services.CreateMask, logger configuration.Logger) *CreateEmailMaskController {
	return &CreateEmailMaskController{createMask: createMask, logger: logger}
}

var InvalidRequestBody = errors.New("error while decoding request body")
var CreateMaskError = errors.New("error creating mask")

func (c CreateEmailMaskController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	type bodyType struct {
		Alias  string `validate:"required" json:"alias"`
		Target string `validate:"required" json:"target"`
	}

	var body bodyType
	var err error

	err = json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		msg := fmt.Errorf("%w. %v", InvalidRequestBody, err).Error()
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		msg := fmt.Errorf("%w. %v", InvalidRequestBody, err).Error()
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	err = c.createMask.Execute(services.CreateMaskCommand{
		Alias:  body.Alias,
		Target: body.Target,
	})

	if err != nil {
		c.logger.Error("Error while creating mask %v", err)
		http.Error(writer, CreateMaskError.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
