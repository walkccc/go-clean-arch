package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/walkccc/go-clean-arch/internal/db/mock"
	"github.com/walkccc/go-clean-arch/internal/repository"
	"github.com/walkccc/go-clean-arch/internal/util"
	pb "github.com/walkccc/go-clean-arch/pkg"
)

type eqCreateUserParamsMatcher struct {
	arg      repository.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(repository.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("Matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserParams(arg repository.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
	password := util.RandomPassword()
	user := randomUser()

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, gotUser *repository.User, gotError error)
	}{
		{
			name: "Success",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mock.MockStore) {
				arg := repository.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, gotUser *repository.User, gotError error) {
				require.NoError(t, gotError)
				require.Equal(t, user.Username, gotUser.Username)
				require.Equal(t, user.HashedPassword, gotUser.HashedPassword)
				require.Equal(t, user.FullName, gotUser.FullName)
				require.Equal(t, user.Email, gotUser.Email)
			},
		},
		{
			name: "AlreadyExists",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					// 23505 := unique_violation
					Return(repository.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, gotUser *repository.User, gotError error) {
				require.Equal(t, &pq.Error{Code: "23505"}, gotError)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			testCase.buildStubs(store)

			userUsecase := NewCreateUserUsecase(store, time.Second*2)
			gotUser, err := userUsecase.CreateUser(context.Background(), testCase.req)
			testCase.checkResponse(t, gotUser, err)
		})
	}
}

func randomUser() repository.User {
	return repository.User{
		Username: util.RandomUsername(),
		FullName: util.RandomFullName(),
		Email:    util.RandomEmail(),
	}
}
