package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, description string) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: description,
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badReques := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalidArgument := status.New(codes.InvalidArgument, "invalid argument")

	// WithDetails returns a new status with the provided details.
	statusDetails, err := statusInvalidArgument.WithDetails(badReques)
	if err != nil {
		return statusInvalidArgument.Err()
	}

	return statusDetails.Err()
}
