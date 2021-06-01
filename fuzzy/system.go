package fuzzy

import (
	"fmt"

	"fugologic.git/graph"
)

// System groups engines and evaluate them all
// All engines are evaluated sequentially
type System []Engine

// NewSystem checks for errors, reorder the engines and creates a new system
func NewSystem(engines []Engine) (System, error) {
	tmp := System(engines)
	if err := tmp.checkDuplicatedOutputs(); err != nil {
		return nil, err
	}

	return tmp.reorder()
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

// reorder builds a graph of engines, check the presence of cycles and flatten the created graph
func (sys System) reorder() (System, error) {
	// To nodes
	nodes := make([]*graph.Node, len(sys))
	for i, engine := range sys {
		nodes[i] = graph.NewNode(engine)
	}

	// Init graph
	dg := graph.NewDirectedGraph(nodes)
	addEdge := func(i, j int) {
		dg.AddEdge(nodes[i], nodes[j])
	}

	// Returns true if a common IDVal is found
	hasCommon := func(a, b []IDSet) bool {
		for _, a1 := range a {
			for _, b1 := range b {
				if a1.parent == b1.parent {
					return true
				}
			}
		}
		return false
	}

	// Add edges
	for i, iEng := range sys {
		// Edge at the current engine
		iInputs := iEng.Inputs()
		iOutputs := iEng.Outputs()
		if hasCommon(iOutputs, iInputs) {
			addEdge(i, i)
		}

		// Edges with the other engines
		for j := i + 1; j < len(sys); j++ {
			jEng := sys[j]
			jInputs := jEng.Inputs()
			jOutputs := jEng.Outputs()
			if hasCommon(iOutputs, jInputs) {
				addEdge(i, j)
			}
			if hasCommon(jOutputs, iInputs) {
				addEdge(j, i)
			}
		}
	}

	// Check an sort nodes
	flat, err := dg.Flatten()
	if err != nil {
		return nil, err
	}

	// To engines
	engines := make([]Engine, len(sys))
	for i, node := range flat {
		engines[i] = node.Data().(Engine)
	}
	return engines, nil
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
