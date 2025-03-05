package tests

import (
	"context"

	"github.com/um3ra/auth-microservice/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ApiTestSuite) TestIntegration_AuthApi_Register() {
	req := auth_v1.RegisterRequest{
		Name:            "michael",
		Email:           "michael@mail.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}
	res, err := s.authClient.Register(context.Background(), &req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotEmpty(res.Token)
}

func (s *ApiTestSuite) TestIntegration_AuthApi_Register_Validation_Failure() {
	testCases := []struct {
		name string
		req  auth_v1.RegisterRequest
	}{
		{
			name: "Password too short",
			req: auth_v1.RegisterRequest{
				Name:            "michael",
				Email:           "michael@mail.com",
				Password:        "1234",
				ConfirmPassword: "1234",
			},
		},
		{
			name: "Invalid email format",
			req: auth_v1.RegisterRequest{
				Name:            "michael",
				Email:           "michae",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
		},
	}

	for i := range testCases {
		tc := &testCases[i]
		s.Run(tc.name, func() {
			res, err := s.authClient.Register(context.Background(), &tc.req)

			if err != nil {
				grpcErr, _ := status.FromError(err)
				s.T().Logf("Validation failed as expected: %s (code: %v)", grpcErr.Message(), grpcErr.Code())
			}

			s.Require().Error(err, "Expected validation error but got none")
			s.Require().Nil(res, "Response should be nil on validation failure")
			grpcErr, _ := status.FromError(err)
			s.Require().Equal(codes.InvalidArgument, grpcErr.Code(), "Expected Invalid argument error")
		})
	}
}
