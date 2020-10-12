package mrtg

type Metric struct {
	CurrentIn  uint64
	CurrentOut uint64
	MaxIn      uint64
	MaxOut     uint64
	AverageIn  uint64
	AverageOut uint64
	// Average max 5 min values
	// https://github.com/oetiker/mrtg/blob/b51d650821287ab3f08e36ed465d23112d99c013/src/bin/mrtg#L1869-L1875
	AverageMaxIn  uint64
	AverageMaxOut uint64
}

type Metrics struct {
	Daily   Metric
	Weekly  Metric
	Monthly Metric
	Yearly  Metric
}
