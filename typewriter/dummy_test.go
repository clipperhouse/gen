package typewriter

// +test foo:"bar"
type dummy int

type (
	// +test foo:"bar"
	dummy2 map[string]dummy

	// +test foo:"bar"
	dummy3 string
)
