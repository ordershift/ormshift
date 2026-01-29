package testutils

import (
	"database/sql"
	"strings"
	"testing"
)

func AssertErrorMessage(t *testing.T, expectedErrorMessage string, err error, functionName string) {
	if err == nil {
		t.Errorf("%s should return not nil error", functionName)
		return
	}
	if err.Error() != expectedErrorMessage {
		t.Errorf(
			"%s error message expected [%s], but returned [%s]",
			functionName,
			expectedErrorMessage,
			err.Error(),
		)
	}
}

func AssertNotNilResultAndNilError[R any](t *testing.T, result *R, err error, functionName string) bool {
	res := true
	if result == nil {
		t.Errorf("%s cannot return nil Result", functionName)
		res = false
	}
	if !AssertNilError(t, err, functionName) {
		res = false
	}
	return res
}

func AssertNilResultAndNotNilError[R any](t *testing.T, result *R, err error, functionName string) bool {
	res := true
	if result != nil {
		t.Errorf("%s should return nil Result", functionName)
		res = false
	}
	if !AssertNotNilError(t, err, functionName) {
		res = false
	}
	return res
}

func AssertNilError(t *testing.T, err error, functionName string) bool {
	if err == nil {
		return true
	}
	t.Errorf("%s cannot return error %q", functionName, err.Error())
	return false
}

func AssertNotNilError(t *testing.T, err error, functionName string) bool {
	if err != nil {
		return true
	}
	t.Errorf("%s should return not nil error", functionName)
	return false
}

func AssertEqualWithLabel[T comparable](t *testing.T, expected, returned T, label string) bool {
	if expected == returned {
		return true
	}
	if label != "" && !strings.HasSuffix(label, ": ") {
		label += ": "
	}
	t.Errorf("%sexpected [%v], but returned [%v]", label, expected, returned)
	return false
}

func AssertNamedArgEqualWithLabel(t *testing.T, expected any, returned sql.NamedArg, label string) bool {
	if expected == returned {
		return true
	}
	expectedNamedArg := expected.(sql.NamedArg)
	if expectedNamedArg.Name == returned.Name && expectedNamedArg.Value == returned.Value {
		return true
	}

	if label != "" && !strings.HasSuffix(label, ": ") {
		label += ": "
	}
	t.Errorf("%sexpected [%v], but returned [%v]", label, expected, returned)
	return false
}
