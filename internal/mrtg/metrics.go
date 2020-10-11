package mrtg

type Metric struct {
	CurrentIn     uint64
	CurrentOut    uint64
	MaxIn         uint64
	MaxOut        uint64
	AverageIn     uint64
	AverageOut    uint64
	AverageMaxIn  uint64
	AverageMaxOut uint64
}

type Metrics struct {
	Daily   Metric
	Weekly  Metric
	Monthly Metric
	Yearly  Metric
}
