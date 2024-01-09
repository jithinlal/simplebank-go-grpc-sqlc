package gapi

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/jithinlal/simplebank/pb"
	"github.com/jithinlal/simplebank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateDeleteAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	account, err := server.store.GetAccount(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "account does not exist")
		}

		return nil, status.Errorf(codes.Internal, "failed to get account: %s", err)
	}

	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to authenticated user")
		return nil, unauthenticatedError(err)
	}

	err = server.store.DeleteAccount(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "account does not exist")
		}

		return nil, status.Errorf(codes.Internal, "failed to get account: %s", err)
	}

	return new(emptypb.Empty), nil
}

func validateDeleteAccountRequest(req *pb.DeleteAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	id := strconv.Itoa(int(req.GetId()))
	if err := util.ValidateId(id); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}
