package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type Authorization struct {}

func (xri Authorization) AccessTokenClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("AccessTokenClientInterceptor")
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	if err != nil {
		log.Fatal(errors.Wrap(err, "read file to buffer"))
	}

	data, err := reader.ReadBytes('\n')

	var tokens map[string]string
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal file to structure"))
	}

	md := metadata.Pairs("Authorization", "Bearer " + tokens["access_token"])
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
