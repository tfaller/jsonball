package startup

import (
	"github.com/tfaller/propchange"
	"github.com/tfaller/propchange/mysql"
)

const envDetectorConnectionString = "DETECTOR_CS"

// GetDetector gets a detector
func GetDetector() (propchange.Detector, error) {
	cs := MustGetEnvVar(envDetectorConnectionString)
	return mysql.NewDetector(cs)
}

// MustGetDetector gets an connector, otherwise panics
func MustGetDetector() propchange.Detector {
	detector, err := GetDetector()
	if err != nil {
		panic(err)
	}
	return detector
}
