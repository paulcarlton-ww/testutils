package testutils

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type (
	// PrepTestI defines function to be called before running a test.
	PrepTestI func(t *testing.T, test *DefTest, setterFunc SetFieldFunc)
	// CheckTestI definesfunction to be called after test to check result.
	CheckTestI func(t *testing.T, actual []interface{}, test *DefTest, getterFunc GetFieldFunc) bool
	// ReportTestI defines function to be called to report test results.
	ReportTestI func(t *testing.T, actual []interface{}, test *DefTest, getterFunc GetFieldFunc)
	// GetFieldFunc is the function to call to get the value of a field of an object
	GetFieldFunc func(t *testing.T, obj interface{}, fieldName string)
	// SetFieldFunc is the function to call to set the value of a field of an object
	SetFieldFunc func(t *testing.T, obj interface{}, fieldName string, value interface{})
	// DefTest generic tests data structure used by tests.
	DefTest struct {
		Number      int         // Test number
		Description string      // Test description
		Config      interface{} // Test configuration, to be used by custom preTest Function
		Inputs      []interface{} // Test inputs
		Expected    []interface{} // Test Expected results
		Status      map[string]interface{} // map of object under test field names and expected values, used by CheckFunc to verify values
		PrepFunc    PrepTestI   // Function to be called before a test,
		// leave unset to call default - which prints the test number and name
		CheckFunc CheckTestI // Function to be called to check a test results,
		// leave unset to call default - which compares actual and expected as strings
		ReportFunc ReportTestI // Function to be called to report test results,
		// leave unset to call default - which reports input, actual and expected as strings
	}

	// ResultsErr is used to hold one or more return value and an error.
	ResultsErr struct {
		Items   uint8         // Number of Results
		Results []interface{} // Items returned
		Err     error         // Error returned
	}
)

var (
	// Prep is the default pre test function.
	Prep = DefaultPrep // nolint:gochecknoglobals // ok
	// Check is the default  post test result check.
	Check = DefaultCheck // nolint:gochecknoglobals // ok
	// Report is the default post test results reporter.
	Report = DefaultReport // nolint:gochecknoglobals // ok
	// NilValue the text used in place of a nil value in test report.
	// The user can change this value if needed.
	NilValue = "testutils.ToString returned nil value" // nolint:gochecknoglobals // ok
)

// RestoreDefaultTestFuncs is used to restore the default functions after a series of tests against a function
// Call defer RestoreDefaultTestFuncs() at the start of a test function and then set Prep, Check and Report
// to the functions to be used for testing the function being tested.
func RestoreDefaultTestFuncs() {
	// Default pre test function.
	Prep = DefaultPrep
	// Default post test result check.
	Check = DefaultCheck
	// Default post test results reporter.
	Report = DefaultReport
}

// GetPrepTestFunc calls the pre test function.
func GetPrepTestFunc(test *DefTest) PrepTestI {
	if test.PrepFunc == nil {
		return Prep
	}

	return test.PrepFunc
}

// GetCheckTestsFunc calls the check test function.
func GetCheckTestsFunc(test *DefTest) CheckTestI {
	if test.CheckFunc == nil {
		return Check
	}

	return test.CheckFunc
}

// GetReportTestsFunc calls the report test function.
func GetReportTestsFunc(test *DefTest) ReportTestI {
	if test.ReportFunc == nil {
		return Report
	}

	return test.ReportFunc
}

// DefaultPrep is the default pre test function that prints the test number and name.
func DefaultPrep(t *testing.T, test *DefTest, setterFunc SetFieldFunc) {
	t.Logf("Test: %d, %s\n", test.Number, test.Description)
}

// DefaultCheck is the default check test function that compares actual and expected as strings.
func DefaultCheck(t *testing.T, actual []interface{}, test *DefTest, getterFunc GetFieldFunc) bool {

	return reflect.DeepEqual(actual, test.Expected) && !FailTests
}

// DefaultReport is the default report test results function reports input, actual and expected as strings.
func DefaultReport(t *testing.T, actual []interface{}, test *DefTest, getterFunc GetFieldFunc) {
	t.Errorf("\nTest: %d, %s\nInput...: %s\nGot.....: %s\nExpected: %s",
		test.Number, test.Description, spew.Sdump(test.Inputs), spew.Sdump(actual), spew.Sdump(test.Expected))
}

// PostTestActions call after test to call check function and report function if check fails.
func PostTestActions(t *testing.T, result []interface{}, test *DefTest, getterFunc GetFieldFunc) {
	if !GetCheckTestsFunc(test)(t, result, test, getterFunc) {
		t.Logf("Test failed")
		GetReportTestsFunc(test)(t, result, test, getterFunc)
	}
}
