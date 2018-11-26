// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package procfs

import "testing"

func TestIOPressure(t *testing.T) {
	fs := FS("fixtures")

	// load io pressure

	iopressure, err := fs.NewResourcePressure("io")
	if err != nil {
		t.Fatalf("could not load io pressure: %v", err)
	}
	if iopressure.Some == nil {
		t.Fatal("expected to get 'some' metrics in io pressure")
	}
	if iopressure.Full == nil {
		t.Fatal("expected to get 'full' metrics in io pressure")
	}

	// check some data

	if iopressure.Full.TotalMicroseconds != 5933015 {
		t.Fatalf("expected full io pauses to be 5933015µs, got %v", iopressure.Full.TotalMicroseconds)
	}
	if iopressure.Some.TotalMicroseconds != 6164237 {
		t.Fatalf("expected some io pauses to be 6164237µs, got %v", iopressure.Some.TotalMicroseconds)
	}

	if iopressure.Full.CongestedPercent10Seconds != 0.01 {
		t.Fatalf("expected full io congestion time to be 0.01%% over 10s, got %v", iopressure.Full.CongestedPercent10Seconds)
	}
	if iopressure.Some.CongestedPercent300Seconds != 10.3 {
		t.Fatalf("expected some io congestion time to be 10.3%% over the last 300s, got %v", iopressure.Some.CongestedPercent300Seconds)
	}
}

func TestMemoryPressure(t *testing.T) {
	fs := FS("fixtures")

	// load io pressure

	mempressure, err := fs.NewResourcePressure("memory")
	if err != nil {
		t.Fatalf("could not load memory pressure: %v", err)
	}
	if mempressure.Some == nil {
		t.Fatal("expected to get 'some' metrics in memory pressure")
	}
	if mempressure.Full == nil {
		t.Fatal("expected to get 'full' metrics in memory pressure")
	}

	// check some data

	if mempressure.Full.TotalMicroseconds != 1 {
		t.Fatalf("expected full memory pauses to be 1µs, got %v", mempressure.Full.TotalMicroseconds)
	}
	if mempressure.Some.TotalMicroseconds != 10 {
		t.Fatalf("expected some memory pauses to be 10µs, got %v", mempressure.Some.TotalMicroseconds)
	}

	if mempressure.Full.CongestedPercent300Seconds != 0.01 {
		t.Fatalf("expected full memory congestion time to be 0.01%% over 300s, got %v", mempressure.Full.CongestedPercent300Seconds)
	}
	if mempressure.Some.CongestedPercent300Seconds != 0.02 {
		t.Fatalf("expected some memory congestion time to be 0.02%% over the last 300s, got %v", mempressure.Some.CongestedPercent300Seconds)
	}
}

func TestCPUPressure(t *testing.T) {
	fs := FS("fixtures")

	// load io pressure

	cpupressure, err := fs.NewResourcePressure("cpu")
	if err != nil {
		t.Fatalf("could not load cpu pressure: %v", err)
	}
	if cpupressure.Some == nil {
		t.Fatal("expected to get 'some' metrics in cpu pressure")
	}
	if cpupressure.Full != nil {
		t.Fatalf("expected to get no 'full' metrics in cpu pressure, got %v", cpupressure.Full)
	}

	// check some data

	if cpupressure.Some.TotalMicroseconds != 29151915 {
		t.Fatalf("expected some io pauses to be 29151915µs, got %v", cpupressure.Some.TotalMicroseconds)
	}

	if cpupressure.Some.CongestedPercent300Seconds != 5.0 {
		t.Fatalf("expected some io congestion time to be 5.0%% over the last 300s, got %v", cpupressure.Some.CongestedPercent300Seconds)
	}
}
