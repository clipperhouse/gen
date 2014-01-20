package test

import (
"testing"
//"os"
"time"
)

/* Use short mode not to execute flaky time-based tests */

func TestBasicString(t *testing.T) {
  f := MakeStringFuture(func () String { return "foo"})
  if f.Get() != "foo" {
    t.Fail()
  }
}

func testBasicInt(t *testing.T) {
  f := MakeIntFuture(func () Int { return 1})
  if f.Get() != 1 {
    t.Fail()
  }
}

func testFlakyInt(t *testing.T) {
  t.Parallel()
  if testing.Short() {
    t.Skip()
  }
  f := MakeIntFuture(func () Int {
    time.Sleep(time.Second)
    return 1
  })
  if f.IsDone() {
    t.Fail()
  }
  time.Sleep(time.Second * time.Duration(2))
  if !f.IsDone() {
    t.Fail()
  }
}

