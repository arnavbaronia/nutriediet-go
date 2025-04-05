package constants

var PackageDayMap = map[string]int{
	"4 weeks":  28,
	"8 weeks":  56,
	"12 weeks": 84,
}

var SuperAdminEmail = "nutriedietplan@gmail.com"

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
