package gapi

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	db "github.com/jithinlal/simplebank/db/sqlc"
	"github.com/jithinlal/simplebank/pb"
	"github.com/jithinlal/simplebank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateAccountParams{
		ID:      req.GetId(),
		Balance: req.GetBalance(),
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

	updatedAccount, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update account: %s", err)
	}

	rsp := &pb.UpdateAccountResponse{Account: convertAccount(updatedAccount)}

	return rsp, nil
}

func validateUpdateAccountRequest(req *pb.UpdateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	id := strconv.Itoa(int(req.GetId()))

	if err := util.ValidateId(id); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := util.ValidateBalance(req.GetBalance()); err != nil {
		violations = append(violations, fieldViolation("balance", err))
	}

	return violations
}
