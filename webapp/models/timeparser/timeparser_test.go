package timeparser

import (
	"reflect"
	"testing"
	"time"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %+v (type %v) - Got %+v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func strptime(str string) time.Time {
	act, _ := time.Parse("2006-01-02 15:04:05 -0700", str)
	return act
}

var nowbak func() time.Time

func startup() {
	nowbak = now
	n := strptime("2014-07-01 12:00:00 +0900")
	now = func() time.Time {
		return n
	}
}

func teardown() {
	now = nowbak
}

func TestParse(t *testing.T) {
	func() {
		startup()
		defer teardown()

		act, err := Parse("18:30")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:30:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("11:30")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-02 11:30:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("1:3")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-02 01:03:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("07-01 18:30")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:30:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("07/01 18:30")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:30:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("7-1 18:3")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:03:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		_, err := Parse("07/01 11:30")
		if err == nil {
			t.Fatal("should be error")
		}
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("2014-07-01 18:30")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:30:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("2014/07/01 18:30")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:30:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		act, err := Parse("2014-7-1 18:3")
		if err != nil {
			t.Fatal(err)
		}
		exp := strptime("2014-07-01 18:03:00 +0900")
		expect(t, act, exp)
	}()

	func() {
		startup()
		defer teardown()

		_, err := Parse("2014/07/01 11:30")
		if err == nil {
			t.Fatal("should be error")
		}
	}()

	func() {
		startup()
		defer teardown()

		_, err := Parse("hoge")
		if err == nil {
			t.Fatal("should be error")
		}
	}()
}
