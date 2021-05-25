package builder

import "fugologic.git/fuzzy"

// NewMamdaniBuilder sets the default configuration for Mamdani inference system
func NewMamdaniBuilder() Builder {
	return Builder{
		and:  fuzzy.ConnectorZadehAnd,
		or:   fuzzy.ConnectorZadehOr,
		impl: fuzzy.ImplicationMin,
	}
}

// NewSugenoBuilder sets the default configuration for Takagi-Sugeno inference system
func NewSugenoBuilder() Builder {
	return Builder{
		and:  fuzzy.ConnectorZadehAnd,
		or:   fuzzy.ConnectorZadehOr,
		impl: fuzzy.ImplicationProd,
	}
}
