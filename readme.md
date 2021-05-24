# fugologic

Fugologic is a naive implementation of a fuzzy logic system.

## Getting started

For more example, see [/fugologic/fuzzy/engine_test.go](https://github.com/sbiemont/fugologic/blob/master/fuzzy/engine_test.go)

## Define the system

### Crisp values definition

Defuzzification requires a crisp interval of discrete values.

It is defined as `crisp.Set` (x min, x max, dx)

```go
// Each values from 0.0 to 0.3 every 0.1 => [0.0, 0.1, 0.2, 0.3]
set, err := crisp.NewSet(0.0, 0.3, 0.1)
if err != nil{
  // Error if the crisp set is badly defined
}
```

### Fuzzy values definition

Fuzzy values and fuzzy sets are defined as :

* `fuzzy.IDVal`: a fuzzy value that contains,
  * an identifier
  * a list of fuzzy sets (only required for system and/or engine checks)
  * a crisp interval of values (only required for defuzzification)
* `fuzzy.IDSet`: fuzzy set that contains,
  * an identifier
  * a membership method
  * its parent `fuzzy.IDVal`

*Notes* :

* every identifier shall be unique
* `fuzzy.IDVal` and `fuzzy.IDSet` can be defined using a random generate ID or a custom one

### Define fuzzy inputs / outputs

First, create a fuzzy value and link it to a list of fuzzy sets.

Ensure that the crisp interval of the fuzzy value covers the fuzzy sets intervals.

```go
// Fuzzy value "a"
setA, _ := crisp.NewSet(-3, 3, 0.1)
fvA := fuzzy.NewIDValCustom("a", setA)

// Fuzzy sets "a1", "a2"
fsA1 := fuzzy.NewIDSetCustom("a1", fuzzy.NewSetTriangular(-3, -1, 1), &fvA)
fsA2 := fuzzy.NewIDSetCustom("a2", fuzzy.NewSetTriangular(-1, 1, 3), &fvA)
```

Create other inputs and outputs the same way.

### Define the rules

A rule is defined with 3 components:

```raw
rule = <expression> <implication> <consequence>
rule = A1 and B1    then          C1, D1
```

* `expression` : connect several fuzzy sets together
* `implication` : define an implication method
* `consequence` : defines several fuzzy sets as the outputs

#### Describe an input expression

Choose the input fuzzy sets and link them using a connector.

2 ways are possible to describe a rule:

* **compact** : uses by default the Zadeh connectors (And, Or)
* **explicit** : choose explicitely connectors for fuzzy sets

Simplest case : the expression has only one premise (directly use the fuzzy set)

```go
// A1
exp := fsA1
```

An expression can be a flat list of several fuzzy sets linked with the same connector.

For example : `A1 and B1 and C1`.

```go
// Using default connectors
// A1 and B1 and C1
exp := fsA1.And(fsB1).And(fsC1)
```

Or in a more explicit way

```go
// A1 and B1 and C1
exp := fuzzy.NewExpression([]fuzzy.Premise{fsA1, fsB1, fsC1}, fuzzy.ConnectorZadehAnd)
```

At last, an expression can be more complex like `(A1 and B1 and C1) or (D1 and E1)`.

```go
// Using default connectors
// (A1 and B1 and C1) or (D1 and E1)
exp := (fsA1.And(fsB1).And(fsC1)).Or(fsD1.And(fsE1))
```

Or in a more explicit way

```go
// A1 and B1 and C1
expABC := fuzzy.NewExpression([]fuzzy.Premise{fsA1, fsB1, fsC1}, fuzzy.ConnectorZadehAnd)

// D1 and E1
expDE := fuzzy.NewExpression([]fuzzy.Premise{fsD1, fsE1}, fuzzy.ConnectorZadehAnd)

// (A1 and B1 and C1) or (D1 and E1)
exp := fuzzy.NewExpression([]fuzzy.Premise{expABC, expDE}, fuzzy.ConnectorZadehOr)
```

#### Describe an implication

An implication links the input expression and the ouput consequence.
Several methods can be chosen like:

* `ImplicationMin` : Mamdani implication minimum
* `ImplicationProd` : Sugeno implication product
* ...

#### Describe an output consequence

A consequence is just a list of fuzzy sets.

#### Write a rule

Combine the several items previously seen to describe the rules.

The first method is useful when describing rules directly int the code (but it uses default connectors)

```go
rules := []fuzzy.Rule{
  // A1 and B1 => C1
  fuzzy.If(fsA1.And(fsB1))     // expression
    .Use(fuzzy.ImplicationMin) // implication
    .Then([]fuzzy.IDSet{fsC1}) // consequence
  // Describe other rules the same way
  // ...
}
```

The second method can be used to easily generate rules from an external builder.
Connectors can be explicitely choosen, unlike for the first method.

```go
rules := []fuzzy.Rule{
  // A1 and B1 => C1
  fuzzy.NewRule(
    fuzzy.NewExpression([]fuzzy.Premise{fsA1, fsB1}, fuzzy.ConnectorZadehAnd), // expression
    fuzzy.ImplicationMin,                                                      // implication
    []fuzzy.IDSet{fsC1},                                                       // consequence
  ),
  // Describe other rules the same way, for example:
  //  * A1 and B2 => C2
  //  * A2 and B1 => C1
  //  * A2 and B2 => C2
}
```

### Choose a defuzzing method

Create a defuzzer with a specific method (like the centroïd method)

```go
defuzzer := fuzzy.NewDefuzzer(fuzzy.DefuzzificationCentroid)
```

### Create an engine

A fuzzy engine combines some rules and a defuzzer.

#### Engine new instance

If the rules contains an error, the engine builder will fail.

```go
engine, err := fuzzy.NewEngine(rules, defuzzer)
if err != nil {
  // An error occurred, check the rules
  return err
}
```

#### Engine evaluation

Then, launch the evaluation process by setting a new input value for each `IDVal` of the engine.

The result contains a crisp value for each fuzzy output value defined.

```go
// Evaluate all the rules of the engine
result, err := engine.Evaluate(fuzzy.DataInput{
  "a": 1,
  "b": 0.05,
})
if err != nil {
  return err
}

// Result
// fuzzy.DataOutput{
//   "c": <crisp result>,
// }
```

### Create a system

A system is an ordered list of engines.
An output of an engine can be linked to the input of another engine.

#### System new instance

When creating a system, some contraints are checked, like:

* all identifiers shall be unique
* an output shall only be produced once
* loops are forbidden : an output cannot be linked to an input of a previous engine

```go
// Create engines
engine1, _ := fuzzy.NewEngine(rules1, defuzzer1)
engine2, _ := fuzzy.NewEngine(rules2, defuzzer2)

// Create and evaluate the system
system, err := NewSystem([]Engine{engine1, engine2})
if err != nil {
  // An error occurred, check the rules
  return err
}
```

#### System evaluation

Then, launch the evaluation process by setting a new input value for each `IDVal` of the system.

The result contains a crisp value for each `IDVal` output value defined.

```go
// Evaluation of the rules of each engines
result, err := system.Evaluate(fuzzy.DataInput{
  "a": 1,
  "b": 0.05,
})
if err != nil {
  return err
}

// Result
// fuzzy.DataOutput{
//   "c": <crisp result>,
// }
```

## Class diagram

Classes used to describe and evaluate a simple fuzzy system

```mermaid
classDiagram
  class Premise {
    <<interface>>
    + Evaluate()
  }
  class Expression {
    + Evaluate()
  }
  class Engine {
    + Evaluate()
  }
  class IDSet {
    + Evaluate()
  }
  class Defuzzer {
    + Defuzz()
  }
  class DataInput {
    - find()
  }
  class System {
    + Evaluate()
  }
  class Rule {
    - evaluate()
  }

  System --> "*" Engine
  DataInput --> "1" IDSet
  Defuzzer --> "1" IDSet : results

  IDSet "1" <--> "*" IDVal

  Premise "*" <-- Expression : premises
  Expression --|> Premise
  IDSet --|> Premise
  DataInput <.. Premise

  Engine --> "*" Rule : rules
  Engine --> "1" Defuzzer : defuzzer
  Engine ..> DataInput
  Engine ..> DataOutput

  Rule --> "*" Premise : inputs
  Rule --> "*" IDSet : outputs
```
