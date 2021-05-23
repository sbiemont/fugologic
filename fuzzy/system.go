package fuzzy

import (
	"fmt"
)

// System groups engines and evaluate them all
// For now, all engines have to be correctly ordered and will be evaluated sequentially
// TODO évaluation séquentielle pour l'instant => à rendre dynamique
type System []Engine

// NewSystem creates a new system and check for errors
func NewSystem(engines []Engine) (System, error) {
	system := System(engines)
	if err := system.check(); err != nil {
		return nil, err
	}
	return system, nil
}

// Evaluate all engines one by one
// The output of the first one is injected into the global input
// The global output is the result of merge of all outputs
func (sys System) Evaluate(input DataInput) (DataOutput, error) {
	newOutput := DataOutput{}
	newInput := input

	for _, eng := range sys {
		output, err := eng.Evaluate(newInput)
		if err != nil {
			return nil, err
		}

		newOutput = mergeData(newOutput, output)
		newInput = mergeData(newInput, output)
	}

	return newOutput, nil
}

// check all possible error of the system and return the first one
func (sys System) check() error {
	var err error
	err = sys.checkUnicity()
	if err != nil {
		return err
	}
	err = sys.checkDuplicatedOutputs()
	if err != nil {
		return err
	}
	err = sys.checkLoop()
	if err != nil {
		return err
	}
	return nil
}

// checkLoop builds an ordered list of inputs / outputs
// * allows an output to become an input
// * forbids an input to become an output
func (sys System) checkLoop() error {
	for i, engI := range sys {
		for _, in := range engI.Inputs() {
			for _, engJ := range sys[i+1:] {
				for _, out := range engJ.Outputs() {
					if in.parent == out.parent {
						return fmt.Errorf("input `%s` cannot become an output", in.parent.uuid)
					}
				}
			}
		}
	}
	return nil
}

// outputs flatten all outputs of the system
func (sys System) outputs() []IDSet {
	var result []IDSet
	for _, eng := range sys {
		result = append(result, eng.Outputs()...)
	}
	return result
}

// checkDuplicatedOutputs controls that an output is not produced twice
func (sys System) checkDuplicatedOutputs() error {
	outputs := sys.outputs()
	for i, outI := range outputs {
		for _, outJ := range outputs[i+1:] {
			if outI.parent == outJ.parent {
				return fmt.Errorf("output `%s` detected twice", outI.parent.uuid)
			}
		}
	}
	return nil
}

// checkUnicity controls the identifiers unicity
func (sys System) checkUnicity() error {
	var idsSets []IDSet
	for _, eng := range sys {
		idsSets = append(idsSets, eng.Inputs()...)
		idsSets = append(idsSets, eng.Outputs()...)
	}
	return newChecker(idsSets).check()
}
