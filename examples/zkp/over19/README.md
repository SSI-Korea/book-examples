# ZKP Example with gnark

- 검증자가 19세 이상인지 증명하는 예제.
- datas 폴더에 생성되는 파일을 적절히 공유
- 아래 순서로 실행한다.

1. 검증자가 회로(Circuit)을 설계하고, proving_key와 circuit_metadata.json 파일을 공유한다.
2. 증명자는 proving_key와 circuit_metadata.json 파일을 이용해서 proof와 public witness 파일을 만들어 공유한다.
3. 검증자는 다시 verification_key와 proof, public witness를 이용해 verify한다.

- 파일로 공유하거나 네트워크로 전송하여 사용한다.

```
#step 1. Verifier_ready
    cd Verifier_ready
    go run Verifier.go

#step 2. Prover
    cd Prover
    go run Prover.go

#step 1. Verifier_verify
    cd Verifier_verify
    go run Verify.go
```