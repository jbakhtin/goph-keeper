package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	auth2 "github.com/jbakhtin/goph-keeper/gen/go/v1/auth"
	"github.com/jbakhtin/goph-keeper/gen/go/v1/kv"

	"github.com/go-faster/errors"
	"github.com/jbakhtin/goph-keeper/internal/client/infrastructure/persistance/grpc/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TokensPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	loginCMD := flag.NewFlagSet("login", flag.ExitOnError)

	loginEmail := loginCMD.String("email", "", "Your email address.")
	loginPassword := loginCMD.String("password", "", "Your password.")

	registrationCMD := flag.NewFlagSet("registration", flag.ExitOnError)

	registerEmail := registrationCMD.String("email", "", "Email address.")
	registerPassword := registrationCMD.String("password", "", "Password.")
	registerPasswordConfirmation := registrationCMD.String("password_confirmation", "", "Password confirmation.")

	logoutCMD := flag.NewFlagSet("logout", flag.ExitOnError)
	logoutType := logoutCMD.Int("type", 0, "Logout type")

	refreshTokenCMD := flag.NewFlagSet("refreshtoken", flag.ExitOnError)

	newKeyValueCMD := flag.NewFlagSet("newkeyvalue", flag.ExitOnError)

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		err := Login(loginCMD, loginEmail, loginPassword)
		if err != nil {
			log.Fatal(err)
		}
	case "registration":
		Registration(registrationCMD, registerEmail, registerPassword, registerPasswordConfirmation)
	case "refreshtoken":
		err := RefreshToken(refreshTokenCMD)
		if err != nil {
			log.Fatal(err)
		}
	case "logout":
		err := Logout(logoutCMD, logoutType)
		if err != nil {
			log.Fatal(err)
		}
	case "newkeyvalue":
		err := NewKeyValue(newKeyValueCMD)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Need to pass command")
	}
}

func Logout(cmd *flag.FlagSet, logoutType *int) error {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	jwtCredentials, err := credentials.NewJWTCredentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(jwtCredentials),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := auth2.NewAuthServiceClient(conn)

	pbLogoutRequest := &auth2.LogoutRequest{
		Type: auth2.LogoutType(*logoutType),
	}

	response, err := client.Logout(context.TODO(), pbLogoutRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Type)
	fmt.Println(response.Sessions)

	return nil
}

func Login(cmd *flag.FlagSet, login, password *string) error {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := auth2.NewAuthServiceClient(conn)

	pbLoginRequest := &auth2.LoginRequest{
		Email:    *login,
		Password: *password,
	}

	response, err := client.Login(context.TODO(), pbLoginRequest)
	if err != nil {
		log.Fatal(err)
	}

	file, err := openFile("./config.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "open file")
	}

	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	writer := bufio.NewWriter(file)
	if _, err = writer.Write(data); err != nil {
		return errors.Wrap(err, "write to buffer")
	}

	if _, err = file.Seek(0, 0); err != nil {
		return errors.Wrap(err, "seek file")
	}

	if err = writer.WriteByte('\n'); err != nil {
		return errors.Wrap(err, "write bytes to file with \\\n")
	}

	return writer.Flush()
}

func RefreshToken(cmd *flag.FlagSet) error {
	conn, err := grpc.Dial(":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(credentials.NewRefreshTokenCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := auth2.NewAuthServiceClient(conn)
	response, err := client.RefreshAccessToken(context.TODO(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)

	err = Write(response)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Registration(cmd *flag.FlagSet, email, password, passwordConfirmation *string) {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := auth2.NewAuthServiceClient(conn)

	pbRegisterRequest := &auth2.RegisterRequest{
		Email:                *email,
		Password:             *password,
		PasswordConfirmation: *passwordConfirmation,
	}

	response, err := client.Register(context.TODO(), pbRegisterRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response, err)
}

func NewKeyValue(cmd *flag.FlagSet) error {
	jwtCredentials, err := credentials.NewJWTCredentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(jwtCredentials),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := kv.NewKeyValueServiceClient(conn)

	pbRegisterRequest := &kv.CrateRequest{
		Key:      "password",
		Value:    "test",
		Metadata: "my new test password",
	}

	response, err := client.Create(context.TODO(), pbRegisterRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response, err)

	return nil
}

func openFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	if os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			return nil, err
		}

		file, err = os.OpenFile(path, flag, perm)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func Write(data any) (err error) {
	file, err := openFile("./config.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer func() {
		err = file.Close()
	}()

	marshaled, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err = file.Write(marshaled); err != nil {
		return errors.Wrap(err, "write to buffer")
	}

	return nil
}
