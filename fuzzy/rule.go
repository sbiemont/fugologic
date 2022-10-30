package fuzzy

// flattenIDSets extracts the IDSets from a list of premises
func flattenIDSets(init []IDSet, premises []Premise) []IDSet {
	for _, premise := range premises {
		switch p := premise.(type) {
		case IDSet:
			init = append(init, p)
		case Expression:
			init = flattenIDSets(init, p.premises)
		}
	}
	return init
}

// Implication links an expression and produces a single fuzzy Set
type Implication func(set Set, k float64) Set

var (
	// ImplicationProd returns the product of a Set with a constant factor
	ImplicationProd Implication = func(set Set, k float64) Set { return set.Multiply(k) }

	// ImplicationMin sets the max upper bound
	ImplicationMin Implication = func(set Set, k float64) Set { return set.Min(k) }
)

// Rule evaluates the input expression + implication + fuzzy output
type Rule struct {
	inputs      Premise
	implication Implication
	outputs     []IDSet
}

// NewRule builds a new Rule instance
// rule = <premise> <implication> <outputs>
// rule = A and B   then          C
func NewRule(inputs Premise, implication Implication, outputs []IDSet) Rule {
	return Rule{
		inputs:      inputs,
		implication: implication,
		outputs:     outputs,
	}
}

// evaluate and return the fuzzy output using crisp input
// Outputs
// * One fuzzy IDSet for each output
// * An error is returned if inputs are missing
func (rule Rule) evaluate(input DataInput) ([]IDSet, error) {
	// Evaluate inputs
	y, err := rule.inputs.Evaluate(input)
	if err != nil {
		return nil, err
	}

	// Evaluate outputs => create a NEW fuzzy Set with the same output ID
	result := make([]IDSet, len(rule.outputs))
	for i, out := range rule.outputs {
		result[i] = IDSet{
			uuid:   out.uuid,
			set:    rule.implication(out.set, y),
			parent: out.parent,
		}
	}

	return result, nil
}

// IO gather and flatten all IDSet from rules' expressions
func (rule Rule) IO() ([]IDSet, []IDSet) {
	return flattenIDSets(nil, []Premise{rule.inputs}), rule.outputs
}

type rules []Rule

// io extracts inputs and outputs IDSet from a list of rules
func (r rules) io() ([]IDSet, []IDSet) {
	var inputs, outputs []IDSet
	for _, rule := range r {
		in, out := rule.IO()
		inputs = append(inputs, in...)
		outputs = append(outputs, out...)
	}
	return inputs, outputs
}
