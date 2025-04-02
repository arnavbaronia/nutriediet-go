package constants

var PackageDayMap = map[string]int{
	"1 Month":  30,
	"2 Months": 60,
	"3 Months": 90,
}

type DietType uint32

const (
	RegularDiet DietType = 1
	DetoxDiet   DietType = 2
	DetoxWater  DietType = 3
)

func (d DietType) Uint32() uint32 {
	return uint32(d)
}

const (
	Motivation = "MOTIVATION"
)
