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

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

// The PSI / pressure interface is described at
//   https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/Documentation/accounting/psi.txt
// Each resource (cpu, io, memory, ...) is exposed as a single file.
// Each file may contain up to two lines, one for "some" pressure and one for "full" pressure.
// Each line contains several averages (over n seconds) and a total in µs.

// example io pressure file:
// > some avg10=0.06 avg60=0.21 avg300=0.99 total=8537362
// > full avg10=0.00 avg60=0.13 avg300=0.96 total=8183134
var pressureLineRE = regexp.MustCompile(`(?P<type>some|full) avg10=(?P<avg10>[0-9.]+) avg60=(?P<avg60>[0-9.]+) avg300=(?P<avg300>[0-9.]+) total=(?P<total>\d+)`)

// ResourcePressureMeasurement represents a pressure measurement, full or partial, for a single resource.
type ResourcePressureMeasurement struct {
	CongestedPercent10Seconds  float64
	CongestedPercent60Seconds  float64
	CongestedPercent300Seconds float64
	TotalMicroseconds          uint64
}

// ResourcePressure is the parsed representation of a resource pressure file.
type ResourcePressure struct {
	Resource string
	Full     *ResourcePressureMeasurement
	Some     *ResourcePressureMeasurement
}

// NewResourcePressure loads the pressure/psi data for a given resource ("io", "memory" or "cpu").
func (fs FS) NewResourcePressure(resource string) (ResourcePressure, error) {
	result := ResourcePressure{Resource: resource}
	data, err := ioutil.ReadFile(fs.Path("pressure", resource))
	if err != nil {
		return result, err
	}

	for _, matches := range pressureLineRE.FindAllStringSubmatch(string(data), -1) {
		elements := make(map[string]string)
		for i, name := range pressureLineRE.SubexpNames() {
			elements[name] = matches[i]
		}
		measurement := &ResourcePressureMeasurement{}
		measurement.TotalMicroseconds, err = strconv.ParseUint(elements["total"], 10, 64)
		if err != nil {
			return result, fmt.Errorf("could not parse total %s: %v", elements["total"], err)
		}
		measurement.CongestedPercent10Seconds, err = strconv.ParseFloat(elements["avg10"], 64)
		if err != nil {
			return result, fmt.Errorf("could not parse avg10 %s: %v", elements["avg10"], err)
		}
		measurement.CongestedPercent60Seconds, err = strconv.ParseFloat(elements["avg60"], 64)
		if err != nil {
			return result, fmt.Errorf("could not parse avg60 %s: %v", elements["avg60"], err)
		}
		measurement.CongestedPercent300Seconds, err = strconv.ParseFloat(elements["avg300"], 64)
		if err != nil {
			return result, fmt.Errorf("could not parse avg300 %s: %v", elements["avg300"], err)
		}

		switch elements["type"] {
		case "some":
			result.Some = measurement
		case "full":
			result.Full = measurement
		default:
			return result, fmt.Errorf("unknown pressure measurement type %s", elements["type"])
		}
	}

	if result.Some == nil && result.Full == nil {
		return result, fmt.Errorf("could not parse pressure file for %s", resource)
	}

	return result, nil
}
