package main

var (
	zero, first, second, third, anotherThird, fourth Thing
	fifth, sixth, seventh, eighth                    Thing
	things, noThings, lotsOfThings                   Things
	others                                           Others
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

	things = Things{
		first,
		second,
		third,
		anotherThird,
		fourth,
	}

	noThings = Things{}

	lotsOfThings = Things{
		first,
		second,
		third,
		fourth,
		fifth,
		sixth,
		seventh,
		eighth,
	}

	others = Others{50, 100, 9, 7, 100, 99}
}
