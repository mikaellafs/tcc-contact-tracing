package service

import (
	"contacttracing/pkg/grpc/pb"
	"contacttracing/pkg/models/dto"
	"contacttracing/pkg/utils"
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

func parseGrpcReport(report *pb.Report) dto.GrpcReport {
	return dto.GrpcReport{
		UserId:            report.GetUserId(),
		DateStartSymptoms: report.GetDateStartSymptoms().AsTime().UnixMilli(),
		DateDiagnostic:    report.GetDateDiagnostic().AsTime().UnixMilli(),
		DateReport:        report.GetDateReport().AsTime().UnixMilli(),
	}
}
