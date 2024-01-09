package gapi

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/jithinlal/simplebank/pb"
	"github.com/jithinlal/simplebank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateGetAccountRequest(req)
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

	rsp := &pb.GetAccountResponse{Account: convertAccount(account)}

	return rsp, nil
}

func validateGetAccountRequest(req *pb.GetAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	id := strconv.Itoa(int(req.GetId()))
	if err := util.ValidateId(id); err != nil {
		violations = append(violations, fieldViolation("id", errors.New("unsupported id")))
	}

	return violations
}
