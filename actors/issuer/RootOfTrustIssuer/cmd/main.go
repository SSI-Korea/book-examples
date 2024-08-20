package main

import (
	"log"
	"net"
	"os"
	"ssi-book/actors/issuer/RootofTrustIssuer"
	"ssi-book/protos"

	"google.golang.org/grpc"
)

func main() {
	argsWithoutProg := os.Args[1:]

	// New Issuer
	issr := new(RootofTrustIssuer.Issuer)
	issr.GenerateDID()

	if len(argsWithoutProg) > 0 {
		issr.CredentialSubjectJsonFilePath = argsWithoutProg[0]
		//loadJson(vcCustomFilePath)
	}

	lis, err := net.Listen("tcp", "localhost:1120")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	issuerServer := RootofTrustIssuer.Server{}
	issuerServer.Issuer = issr

	s := grpc.NewServer()
	protos.RegisterSimpleIssuerServer(s, &issuerServer)

	log.Printf("RootOfTrust Issuer Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
