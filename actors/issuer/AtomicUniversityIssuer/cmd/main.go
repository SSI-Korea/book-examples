package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"ssi-book/actors/issuer/AtomicUniversityIssuer"
	"ssi-book/protos"
)

func main() {
	argsWithoutProg := os.Args[1:]

	// New Issuer
	issr := new(AtomicUniversityIssuer.Issuer)
	issr.GenerateDID()

	if len(argsWithoutProg) > 0 {
		issr.CredentialSubjectJsonFilePath = argsWithoutProg[0]
		//loadJson(vcCustomFilePath)
	}

	lis, err := net.Listen("tcp", "localhost:1124")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	issuerServer := AtomicUniversityIssuer.Server{}
	issuerServer.Issuer = issr

	s := grpc.NewServer()
	protos.RegisterAtomicIssuerServer(s, &issuerServer)

	log.Printf("AtomicUniversityIssuer Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
