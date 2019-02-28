package cmd

import (
	"crypto/tls"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	v1 "minibox.ai/minibox/pkg/api/v1"
	log "minibox.ai/minibox/pkg/logger"
	"minibox.ai/minibox/pkg/oauth2/provider"
)

type ClientOption struct {
	Addr string
}

func saveAccesToken(data []byte) error {
	home, err := homedir.Dir()
	filename := path.Join(home, ".minibox/token")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	return err
}

func loadAccessToken() ([]byte, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	filename := path.Join(home, ".minibox/token")
	return ioutil.ReadFile(filename)
}

type AuthHandler func(token string) error
type GRPCConnHandler func(*grpc.ClientConn) error
type ClientHandler func(*v1.Clients) error

func LoadAuth(handle AuthHandler) error {
	buf, err := loadAccessToken()
	if err != nil {
		log.S().Errorf("invalid token %s", err)
		loginCmd.Execute()
		return err
	}
	tokenGen, err := provider.PublicKeyJWT()
	if err != nil {
		log.S().Errorf("invalid jwt token %s", err)
		return err
	}

	auth, err := tokenGen.DecryptoMap(string(buf))
	if err != nil {
		log.S().Errorf("auth error %s", err)
		return err
	}
	if token, ok := auth["tok"].(string); ok {
		return handle(token)
	}
	return err
}

func GRPCConn(addr string, fn GRPCConnHandler) error {
	return LoadAuth(func(t string) error {
		token := oauth2.Token{AccessToken: t}
		perRPC := oauth.NewOauthAccess(&token)
		opts := []grpc.DialOption{
			// In addition to the following grpc.DialOption, callers may also use
			// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
			// itself.
			// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
			grpc.WithPerRPCCredentials(perRPC),
			// oauth.NewOauthAccess requires the configuration of transport
			// credentials.
			grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
			),
		}

		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			// log.Fatalf("did not connect: %v", err)
			return err
		}
		return fn(conn)
	})
}

func AuthClient(opt *ClientOption, fn ClientHandler) error {

	return GRPCConn(opt.Addr, func(conn *grpc.ClientConn) error {
		client := v1.NewClients(conn)
		return fn(client)
	})
}

func LoadClientConfig() *ClientOption {
	var opt = ClientOption{
		Addr: viper.GetString("apiserver"),
	}
	return &opt
}

func setLoggerLevel(lvl zapcore.Level) {
	logger := log.NewLogger(log.Level(lvl))
	sugar := logger.SugarLogger()
	log.RegistryLogger(logger)
	log.RegistrySugar(sugar)
}

func enableDebug() {
	setLoggerLevel(zapcore.DebugLevel)
}
