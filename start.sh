#!/bin/bash

# Registrar port: 9000
nohup go run vdr/registrar/registrar.go > ./logs/registrar.log 2>&1 &
# Resover port: 9001
nohup go run vdr/resolver/resolver.go > ./logs/resolver.log 2>&1 &

# Issuer ROT port: 1120
nohup go run actors/issuer/RootOfTrustIssuer/cmd/main.go > ./logs/RootOfTrustIssuer.log 2>&1 &
# Issuer University port: 1121
nohup go run actors/issuer/UniversityIssuer/cmd/main.go > ./logs/UniversityIssuer.log 2>&1 &
# Issuer Company port: 1122
nohup go run actors/issuer/CompanyIssuer/cmd/main.go > ./logs/CompanyIssuer.log 2>&1 &
# Issuer Bank port: 1123
nohup go run actors/issuer/BankIssuer/cmd/main.go > ./logs/BankIssuer.log 2>&1 &
# Issuer Atomic University port: 1124
nohup go run actors/issuer/AtomicUniversityIssuer/cmd/main.go > ./logs/AtomicUniversityIssuer.log 2>&1 &
