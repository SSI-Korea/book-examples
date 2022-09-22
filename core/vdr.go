package core

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"ssi-book/protos"
	"time"
)

func RegisterDid(did string, didDocument string) error {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Registrar not connect: %v\n", err)
		return errors.New(fmt.Sprintf("Registrar not connect: %v", err))
	}
	defer conn.Close()

	client := protos.NewRegistrarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RegisterDid(ctx, &protos.RegistrarRequest{Did: did, DidDocument: didDocument})
	if err != nil {
		log.Println("Failed to register DID.")
		return errors.New("Failed to register DID.")
	}

	log.Printf("Registrar Response: %s\n", res)

	return nil
}

func ResolveDid(did string) (string, error) {
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Resolver not connect: %v\n", err)
		return "", errors.New(fmt.Sprintf("Resolver not connect: %v", err))
	}
	defer conn.Close()

	client := protos.NewResolverClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.ResolveDid(ctx, &protos.ResolverRequest{Did: did})
	if err != nil {
		log.Fatalf("Failed to resolve DID.")
	}

	log.Printf("Resolver Response: %s\n", res)

	return res.DidDocument, nil
}
