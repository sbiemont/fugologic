package fuzzy

import (
	"fmt"

	"github.com/sbiemont/fugologic/graph"
	"github.com/sbiemont/fugologic/id"
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

		newOutput = newOutput.merge(output)
		newInput = newInput.merge(output)
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
	edges := graph.NewDirectedEdges()
	addEdge := func(i, j int) {
		edges.Add(nodes[i], nodes[j])
	}

	// Returns true if a common IDVal is found
	hasCommon := func(a, b map[*IDVal]struct{}) bool {
		for b1 := range b {
			if _, exists := a[b1]; exists {
				return true
			}
		}
		return false
	}

	// Compute inputs / outputs of all engines
	type inouts struct {
		inputs  map[*IDVal]struct{}
		outputs map[*IDVal]struct{}
	}
	savedIO := make(map[id.ID]inouts)
	for _, eng := range sys {
		in, out := eng.IO()
		savedIO[eng.uuid] = inouts{
			inputs:  IDSets(in).IDVals(),
			outputs: IDSets(out).IDVals(),
		}
	}

	// Add edges
	for i, iEng := range sys {
		// Edge at the current engine
		iIO := savedIO[iEng.uuid]
		if hasCommon(iIO.outputs, iIO.inputs) {
			addEdge(i, i)
		}

		// Edges with the other engines
		for j := i + 1; j < len(sys); j++ {
			jEng := sys[j]
			jIO := savedIO[jEng.uuid]
			if hasCommon(iIO.outputs, jIO.inputs) {
				addEdge(i, j)
			}
			if hasCommon(jIO.outputs, iIO.inputs) {
				addEdge(j, i)
			}
		}
	}

	// Check an sort nodes
	dg, err := graph.NewDirectedGraph(nodes, edges)
	if err != nil {
		return nil, err
	}

	// To engines
	engines := make([]Engine, len(sys))
	for i, node := range dg.Flatten() {
		engines[i] = node.Data().(Engine)
	}
	return engines, nil
}

// outputs flatten all outputs of the system
func (sys System) outputs() []IDSet {
	var result []IDSet
	for _, eng := range sys {
		_, outputs := eng.IO()
		result = append(result, outputs...)
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
