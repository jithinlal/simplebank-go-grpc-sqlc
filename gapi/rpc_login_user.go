package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/jithinlal/simplebank/db/sqlc"
	"github.com/jithinlal/simplebank/pb"
	"github.com/jithinlal/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no user found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorised: %s", err)
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	metadata := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		User:                  convertUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}

	return rsp, nil
}
