package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/ibks-bank/profile/cmd/swagger"
	"github.com/ibks-bank/profile/config"
	"github.com/ibks-bank/profile/internal/app/profile"
	"github.com/ibks-bank/profile/internal/pkg/auth"
	"github.com/ibks-bank/profile/internal/pkg/email"
	"github.com/ibks-bank/profile/internal/pkg/store"
	gw "github.com/ibks-bank/profile/pkg/profile"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	conf := config.GetConfig()
	grpcPort := "3002"
	tcpPort := "3001"

	pgConnString := fmt.Sprintf(
		"port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Port,
		conf.Database.Address,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Name)

	postgres, err := sql.Open("postgres", pgConnString)
	if err != nil {
		log.Fatal("can't open database")
	}
	st := store.New(postgres)

	auther := auth.NewAuthorizer(
		conf.Auth.SigningKey,
		conf.Auth.HashSalt,
		time.Duration(conf.Auth.TokenTTL)*time.Second,
		st,
	)

	emailer := email.NewSender(
		conf.Auth.Email2FA,
		conf.Auth.Password2FA,
		conf.Auth.SmtpHost,
		conf.Auth.SmtpPort,
	)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(auth.Interceptor))
	gw.RegisterProfileServer(
		s,
		profile.NewServer(st, auther, emailer),
	)
	log.Println("Serving gRPC on 0.0.0.0:" + grpcPort)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:"+grpcPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
		if s == auth.TokenKey {
			return s, true
		}

		return s, false
	}))
	err = gw.RegisterProfileHandler(ctx, gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	gwServer := &http.Server{
		Addr:    ":" + tcpPort,
		Handler: mux,
	}

	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("/docs/", staticServer)
	mux.Handle("/docs/", sh)

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:" + tcpPort)
	log.Fatalln(gwServer.ListenAndServe())
}
