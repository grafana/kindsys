package kindsys

import "fmt"

// TODO docs
type Maturity string

const (
	MaturityMerged       Maturity = "merged"
	MaturityExperimental Maturity = "experimental"
	MaturityStable       Maturity = "stable"
	MaturityMature       Maturity = "mature"
)

func maturityIdx(m Maturity) int {
	// icky to do this globally, this is effectively setting a default
	if string(m) == "" {
		m = MaturityMerged
	}

	for i, ms := range maturityOrder {
		if m == ms {
			return i
		}
	}
	panic(fmt.Sprintf("unknown maturity milestone %s", m))
}

var maturityOrder = []Maturity{
	MaturityMerged,
	MaturityExperimental,
	MaturityStable,
	MaturityMature,
}

func (m Maturity) Less(om Maturity) bool {
	return maturityIdx(m) < maturityIdx(om)
}

func (m Maturity) String() string {
	return string(m)
}
