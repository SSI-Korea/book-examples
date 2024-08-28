package over19

import (
	"github.com/consensys/gnark/frontend"
)

// AgeCircuit 정의
type AgeCircuit struct {
	Age    frontend.Variable `gnark:",secret"`
	MinAge frontend.Variable `gnark:",public"`
}

func (circuit *AgeCircuit) Define(api frontend.API) error {
	api.AssertIsLessOrEqual(circuit.MinAge, circuit.Age)
	return nil
}

type CircuitConfig struct {
	Circuit string `json:"circuit"`
	Inputs  map[string]struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Value       int    `json:"value,omitempty"`
	} `json:"inputs"`
	Description string `json:"description"`
}
