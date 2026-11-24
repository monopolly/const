package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite examples/*/*.go and *.swift")

// TestExamples generates every examples/*/*.const and compares the result with
// the committed .go / .swift files. Run `go test -update` after changing the generator.
// `go build ./...` then checks that the generated Go code compiles.
func TestExamples(t *testing.T) {
	inputs, err := filepath.Glob("examples/*/*.const")
	if err != nil || len(inputs) == 0 {
		t.Fatal("no examples found")
	}

	for _, input := range inputs {
		t.Run(filepath.Base(input), func(t *testing.T) {
			src, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}

			goCode, swiftCode := generate(string(src))
			base := strings.TrimSuffix(input, ".const")

			compare(t, base+".go", goCode)
			if swiftCode != nil {
				compare(t, base+".swift", swiftCode)
			}
		})
	}
}

func compare(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run go test -update)", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s is out of date, run go test -update", path)
	}
}

func TestGoName(t *testing.T) {
	cases := map[string]string{
		"new":                                "New",
		"user_id":                            "UserId",
		"sensor.alarm.*":                     "SensorAlarmAny",
		"sensor.>":                           "SensorAll",
		">":                                  "All",
		"myapp.us-east.users.profile.update": "MyappUsEastUsersProfileUpdate",
		"obs.gnss.BEMT":                      "ObsGnssBEMT",
		"game.room.42.action":                "GameRoom42Action",
		"42.answer":                          "N42Answer",
		"wasmbus.ctl.v1.{lattice}.{host-id}": "WasmbusCtlV1LatticeHostId",
		"$SYS.REQ.SERVER.PING":               "SysReqServerPing",
		"$JS.API.STREAM.INFO.ORDERS":         "JsApiStreamInfoOrders",
	}
	for in, want := range cases {
		if got := goName(in); got != want {
			t.Errorf("goName(%q) = %q, want %q", in, got, want)
		}
	}

	if got := goUpperName("sensor.alarm.*"); got != "SENSOR_ALARM_ANY" {
		t.Errorf("goUpperName = %q", got)
	}
	if got := swiftName("$SYS.REQ.SERVER.PING"); got != "sysReqServerPing" {
		t.Errorf("swiftName = %q", got)
	}
}
