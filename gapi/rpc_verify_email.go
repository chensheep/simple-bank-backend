package gapi

import (
	"context"

	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/chensheep/simple-bank-backend/pb"
	"github.com/chensheep/simple-bank-backend/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailReq(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	res, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify email: %v", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: res.User.IsEmailVerified,
	}

	return rsp, nil
}

func validateVerifyEmailReq(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailID(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err.Error()))
	}
	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err.Error()))
	}

	return violations
}
