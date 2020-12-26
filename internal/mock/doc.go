// Package mock is simpley here to store all needed mocks
// To run tests you have to run go generate on this file.
package mock

//go:generate go run github.com/golang/mock/mockgen -destination ./mock_jsonball/jsonball.go github.com/tfaller/jsonball Registry
//go:generate go run github.com/golang/mock/mockgen -destination ./mock_propchange/propchange.go github.com/tfaller/propchange Detector,DocumentOps
