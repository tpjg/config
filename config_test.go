package config

import (
	"os"
	"testing"
	"time"
)

type TestStruct struct {
	S          string        `flag:"TEST_STRING"`
	I          int           `flag:"TEST_INT"`
	I64        int64         `flag:"TEST_INT64"`
	F64        float64       `flag:"TEST_FLOAT64"`
	D          time.Duration `flag:"TEST_DURATION"`
	UI         uint          `flag:"TEST_UINT"`
	UI64       uint64        `flag:"TEST_UINT64"`
	unexported string        `flag:"CANNOT_BE_SET"`
}

func TestParse(t *testing.T) {
	os.Setenv("TEST_STRING", "teststringenv")
	os.Setenv("TEST_INT", "421")
	os.Setenv("TEST_INT64", "42421")
	os.Setenv("TEST_FLOAT64", "3.14161")
	os.Setenv("TEST_DURATION", "7h2m3s")
	os.Setenv("TEST_UINT", "0xf")
	os.Setenv("TEST_UINT64", "0xfff")
	os.Setenv("CANNOT_BE_SET", "value")

	ts := TestStruct{}
	ReadStructFromEnv(&ts)

	if ts.S != "teststringenv" {
		t.Errorf("Unexpected value: %v", ts.S)
	}
	if ts.I != 421 {
		t.Errorf("Unexpected value: %v", ts.I)
	}
	if ts.I64 != 42421 {
		t.Errorf("Unexpected value: %v", ts.I64)
	}
	if ts.F64 != 3.14161 {
		t.Errorf("Unexpected value: %v", ts.F64)
	}
	if ts.D != 7*time.Hour+2*time.Minute+3*time.Second {
		t.Errorf("Unexpected value: %v", ts.D)
	}
	if ts.UI != 0xf {
		t.Errorf("Unexpected value: %v", ts.UI)
	}
	if ts.UI64 != 0xfff {
		t.Errorf("Unexpected value: %v", ts.UI64)
	}
	if ts.unexported != "" {
		t.Errorf("Unexpected value: %v", ts.unexported)
	}
}

type TestStructWithDefault struct {
	S          string        `flag:"DTEST_STRING" default:"teststring"`
	I          int           `flag:"DTEST_INT" default:"42"`
	I64        int64         `flag:"DTEST_INT64" default:"4242"`
	F64        float64       `flag:"DTEST_FLOAT64" default:"3.1416"`
	D          time.Duration `flag:"DTEST_DURATION" default:"1h2m3s"`
	UI         uint          `flag:"DTEST_UINT" default:"0xff"`
	UI64       uint64        `flag:"DTEST_UINT64" default:"0xffff"`
	unexported string        `flag:"CANNOT_BE_SET" default:"should not set"`
}

func TestParseWithDefault(t *testing.T) {
	ts := TestStructWithDefault{}
	ReadStructFromEnv(&ts)

	if ts.S != "teststring" {
		t.Errorf("Unexpected value: %v", ts.S)
	}
	if ts.I != 42 {
		t.Errorf("Unexpected value: %v", ts.I)
	}
	if ts.I64 != 4242 {
		t.Errorf("Unexpected value: %v", ts.I64)
	}
	if ts.F64 != 3.1416 {
		t.Errorf("Unexpected value: %v", ts.F64)
	}
	if ts.D != 1*time.Hour+2*time.Minute+3*time.Second {
		t.Errorf("Unexpected value: %v", ts.D)
	}
	if ts.UI != 0xff {
		t.Errorf("Unexpected value: %v", ts.UI)
	}
	if ts.UI64 != 0xffff {
		t.Errorf("Unexpected value: %v", ts.UI64)
	}
	if ts.unexported != "" {
		t.Errorf("Unexpected value: %v", ts.unexported)
	}
}

func TestParseWithOverrideOfDefault(t *testing.T) {
	os.Setenv("DTEST_STRING", "something else")
	os.Setenv("DTEST_INT", "-42")
	os.Setenv("DTEST_INT64", "-4242")
	os.Setenv("DTEST_FLOAT64", "-3.1416")
	os.Setenv("DTEST_DURATION", "3h2m1s")
	os.Setenv("DTEST_UINT", "0xaa")
	os.Setenv("DTEST_UINT64", "0xaaaa")

	ts := TestStructWithDefault{}
	ReadStructFromEnv(&ts)

	if ts.S != "something else" {
		t.Errorf("Unexpected value: %v", ts.S)
	}
	if ts.I != -42 {
		t.Errorf("Unexpected value: %v", ts.I)
	}
	if ts.I64 != -4242 {
		t.Errorf("Unexpected value: %v", ts.I64)
	}
	if ts.F64 != -3.1416 {
		t.Errorf("Unexpected value: %v", ts.F64)
	}
	if ts.D != 3*time.Hour+2*time.Minute+1*time.Second {
		t.Errorf("Unexpected value: %v", ts.D)
	}
	if ts.UI != 0xaa {
		t.Errorf("Unexpected value: %v", ts.UI)
	}
	if ts.UI64 != 0xaaaa {
		t.Errorf("Unexpected value: %v", ts.UI64)
	}
	if ts.unexported != "" {
		t.Errorf("Unexpected value: %v", ts.unexported)
	}
}

func TestOverrideWithArgs(t *testing.T) {
	tmp := os.Args
	defer func() { os.Args = tmp }()

	os.Args = []string{"TESTPROGNAME", "-DTEST_INT=55"}
	os.Setenv("DTEST_INT", "-555")
	os.Args = append(os.Args, "-DTEST_DURATION=5h4m3s")
	os.Setenv("DTEST_DURATION", "3h2m1s")

	ts := TestStructWithDefault{}
	ReadStructFromEnvOverrideWithArgs(&ts)
	if ts.I != 55 {
		t.Errorf("Unexpected value: %v", ts.I)
	}
	if ts.D != 5*time.Hour+4*time.Minute+3*time.Second {
		t.Errorf("Unexpected value: %v", ts.D)
	}
}
