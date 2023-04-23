package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/walkccc/go-clean-arch/internal/db/mock"
	"github.com/walkccc/go-clean-arch/internal/util"
	"github.com/walkccc/go-clean-arch/internal/util/token"
	pb "github.com/walkccc/go-clean-arch/pkg"
)

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name       string
		req        *pb.LoginUserRequest
		buildStubs func(store *mock.MockStore)
	}{
		{
			name: "Success",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			testCase.buildStubs(store)

			tokenMaker, _ := token.NewPasetoMaker(util.RandomString(32))
			authUsecase := NewLoginUserUsecase(store, tokenMaker, time.Minute*15, time.Hour*24, time.Second*2)
			loggedInUser, _, err := authUsecase.LoginUser(context.Background(), testCase.req)
			require.NoError(t, err)
			require.Equal(t, user.Username, loggedInUser.Username)
			require.Equal(t, user.FullName, loggedInUser.FullName)
			require.Equal(t, user.Email, loggedInUser.Email)
			require.Equal(t, user.HashedPassword, loggedInUser.HashedPassword)
		})
	}
}
