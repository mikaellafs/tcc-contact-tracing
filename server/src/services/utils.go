package services

import (
	"contacttracing/src/models/dto"
	"contacttracing/src/utils"
	"errors"
	"net/http"
)

func validateGrpcMessage(request any, pk string, signature []byte) (dto.GrpcResult, error) {
	var result dto.GrpcResult

	isValid, err := utils.ValidateMessage(request, pk, signature)
	if err != nil {
		result.Status = http.StatusBadRequest
		result.Message = "Failed to validate message: " + err.Error()
		return result, errors.New(result.Message)
	}

	if !isValid {
		result.Status = http.StatusForbidden
		result.Message = "Signature is not valid for this message"
		return result, errors.New(result.Message)
	}

	return result, nil
}
