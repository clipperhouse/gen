package benchmarks

import (
	"fmt"
	"testing"
)

var (
	dummyObjects      dummyObjectSlice
	globalBoolResult  bool
	globalSliceResult []*dummyDestinationSelectObject
	globalSliceResul2 []*dummyObject
)

//Any
func BenchmarkDummyObjAny_Generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		globalBoolResult = dummyObjects.Any(func(d *dummyObject) bool { return d.Num > 10000 })
	}
}

func BenchmarkDummyObjAny_NativeLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		any := false
		for _, d := range dummyObjects {
			if d.Num > 10000 {
				any = true
				break
			}
		}
		globalBoolResult = any
	}
}

//Select
func BenchmarkDummyObjSelect_Generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		globalSliceResult = dummyObjects.SelectDummyDestinationSelectObject(func(d *dummyObject) *dummyDestinationSelectObject { return &dummyDestinationSelectObject{d.Name} })
		globalBoolResult = len(globalSliceResult) == len(dummyObjects)
	}
}

func BenchmarkDummyObjSelect_NativeLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		globalSliceResult = []*dummyDestinationSelectObject{}
		for _, d := range dummyObjects {
			globalSliceResult = append(globalSliceResult, &dummyDestinationSelectObject{d.Name})
		}
		globalBoolResult = len(globalSliceResult) == len(dummyObjects)
	}
}

//SortBy
func BenchmarkDummyObjSortBy_Generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		globalSliceResul2 = dummyObjects.SortBy(func(a *dummyObject, b *dummyObject) (isLess bool) { return a.Num%3 == 0 })
	}
}

/*TODO:
func BenchmarkDummyObjSortBy_NativeLoop(b *testing.B) {

}*/

func init() {
	dummyObjects = dummyObjectSlice([]*dummyObject{})
	for i := 0; i < 10000; i++ {
		dummyObjects = append(dummyObjects, &dummyObject{fmt.Sprintf("Name %d", i), i})
	}
}
