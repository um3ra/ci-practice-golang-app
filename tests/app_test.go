

package tests

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/um3ra/auth-microservice/internal/app"
	authGrpc "github.com/um3ra/auth-microservice/pkg/auth_v1"
	"github.com/um3ra/auth-microservice/pkg/interceptor"
	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcTestAddr = os.Getenv("GRPC_TEST_ADDR")
)


type ApiTestSuite struct {
	suite.Suite
	grpcServer *grpc.Server
	userClient userGrpc.UserV1Client
	authClient authGrpc.AuthV1Client
	grpcConn   *grpc.ClientConn
}

// SetupTestSuite has a SetupTest method, which will run before each
// test in the suite.
func (s *ApiTestSuite) SetupSuite() {

	ctx := context.Background()
	provider := app.NewServiceProviderTest()

	// grpc server init
	s.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.UnaryInterceptor(interceptor.ValidateInterceptor))

	userHandler := provider.UserHandler(ctx)
	authHandler := provider.AuthHandler(ctx)

	userGrpc.RegisterUserV1Server(s.grpcServer, userHandler)
	authGrpc.RegisterAuthV1Server(s.grpcServer, authHandler)

	l, err := net.Listen("tcp", grpcTestAddr)
	if err != nil {
		log.Fatalf("Failed to start grpc server: %s", err.Error())
	}
	go func() {
		if err = s.grpcServer.Serve(l); err != nil {
			log.Fatalf("grpc server exited: %s", err.Error())
		}
	}()

	time.Sleep(time.Second)

	// grpc client init
	conn, err := grpc.NewClient(grpcTestAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create grpc client connection: %s", err.Error())
	}
	s.grpcConn = conn
	authClient := authGrpc.NewAuthV1Client(conn)
	userClient := userGrpc.NewUserV1Client(conn)
	s.authClient = authClient
	s.userClient = userClient
}

func (s *ApiTestSuite) TearDownSuite() {

	_ = s.grpcConn.Close()
	s.grpcServer.GracefulStop()
}

func Test_HttpTestSuite_Integration(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
