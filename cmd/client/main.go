package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/goph-keeper/internal/client/infastructure/network/interceptors/unary/auth"
	pb "github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc/gen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
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

	logoutCMD := flag.NewFlagSet("registration", flag.ExitOnError)
	logoutType := logoutCMD.String("type", "current", "Logout type")

	refreshTokenCMD := flag.NewFlagSet("refreshtoken", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'login' or 'registration' command")
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
		RefreshToken(refreshTokenCMD)
	case "logout":
		Logout(logoutCMD, logoutType)
	default:
		fmt.Println("Need to pass command")
	}
}

func Logout(cmd *flag.FlagSet, logoutType *string) error {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	authorization := auth.Authorization{}

	conn, err := grpc.Dial(":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authorization.AccessTokenClientInterceptor))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAuthServiceClient(conn)

	pbLogoutRequest := &pb.LogoutRequest{
		Type: pb.LogoutType_TYPE_UNSPECIFIED,
	}

	_, err = client.Logout(context.TODO(), pbLogoutRequest)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}


func Login(cmd *flag.FlagSet, login, password *string) error {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAuthServiceClient(conn)

	pbLoginRequest := &pb.LoginRequest{
		Email: *login,
		Password: *password,
	}

	fmt.Println("client", pbLoginRequest)

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
	tokens, err := Read()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAuthServiceClient(conn)
	response, err := client.RefreshAccessToken(context.TODO(), &pb.RefreshTokenRequest{
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = Write(response)

	return nil
}

func Registration(cmd * flag.FlagSet, email, password, passwordConfirmation *string) {
	fmt.Println("Registration ...")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAuthServiceClient(conn)

	pbRegisterRequest := &pb.RegisterRequest{
		Email: *email,
		Password: *password,
		PasswordConfirmation: *passwordConfirmation,
	}

	fmt.Println("client", pbRegisterRequest)

	response, err := client.Register(context.TODO(), pbRegisterRequest)

	fmt.Println(response, err)
}

func openFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(path, flag, perm)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
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

func Read() (*TokensPair, error) {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	var tokens TokensPair
	err = json.Unmarshal(content, &tokens)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal file to structure")
	}

	return &tokens, nil
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