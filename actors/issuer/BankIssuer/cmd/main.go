package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"ssi-book/actors/issuer/BankIssuer"
	"ssi-book/protos"
)

func main() {
	argsWithoutProg := os.Args[1:]

	// New Issuer
	issr := new(BankIssuer.Issuer)
	issr.GenerateDID()

	if len(argsWithoutProg) > 0 {
		issr.CredentialSubjectJsonFilePath = argsWithoutProg[0]
		//loadJson(vcCustomFilePath)
	}

	lis, err := net.Listen("tcp", "localhost:1123")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	issuerServer := BankIssuer.Server{}
	issuerServer.Issuer = issr

	s := grpc.NewServer()
	protos.RegisterMultipleIssuerServer(s, &issuerServer)

	log.Printf("BankIssuer Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
