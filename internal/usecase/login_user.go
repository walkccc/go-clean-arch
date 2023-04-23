package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/util"
	"github.com/walkccc/go-clean-arch/internal/util/token"
	pb "github.com/walkccc/go-clean-arch/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginMetadata struct {
	SessionId             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type LoginUserUsecase interface {
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*repository.User, *LoginMetadata, error)
}

type loginUserUsecase struct {
	store                repository.Store
	tokenMaker           token.Maker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	timeout              time.Duration
}

func NewLoginUserUsecase(
	store repository.Store,
	tokenMaker token.Maker,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
	timeout time.Duration) LoginUserUsecase {
	return &loginUserUsecase{
		store:                store,
		tokenMaker:           tokenMaker,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		timeout:              timeout}
}

func (u *loginUserUsecase) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*repository.User, *LoginMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, status.Errorf(codes.NotFound, "User '%s' not found", req.GetUsername())
		}
		return nil, nil, status.Errorf(codes.Internal, "Failed to find user '%s'", req.GetUsername())
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, nil, status.Errorf(codes.NotFound, "User '%s' wrong password", req.GetUsername())
	}

	accessToken, accessPayload, err := u.tokenMaker.CreateToken(
		user.Username,
		u.accessTokenDuration)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := u.tokenMaker.CreateToken(
		user.Username,
		u.refreshTokenDuration)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Failed to create refresh token: %s", err)
	}

	mtdt := u.extractMetadata(ctx)
	session, err := u.store.CreateSession(ctx, repository.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Failed to create session: %s", err)
	}

	return &user, &LoginMetadata{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt}, nil
}
