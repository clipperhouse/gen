package main

var (
	zero, first, second, third, anotherThird, fourth Thing
	fifth, sixth, seventh, eighth                    Thing
	things, noThings, lotsOfThings                   ThingSlice
	others                                           OtherSlice
)

func init() {
	zero = Thing{}
	first = Thing{"First", 60}
	second = Thing{"Second", 40}
	third = Thing{"Third", 100}
	anotherThird = Thing{"Third", 100}
	fourth = Thing{"Fourth", 40}
	fifth = Thing{"Fifth", 70}
	sixth = Thing{"Sixth", 10}
	seventh = Thing{"Seventh", 50}
	eighth = Thing{"Eighth", 110}

	things = ThingSlice{
		first,
		second,
		third,
		anotherThird,
		fourth,
	}

	noThings = ThingSlice{}

	lotsOfThings = ThingSlice{
		first,
		second,
		third,
		fourth,
		fifth,
		sixth,
		seventh,
		eighth,
	}

	others = OtherSlice{50, 100, 9, 7, 100, 99}
}
