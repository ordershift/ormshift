package testutils

import (
	"database/sql"
	"strings"
	"testing"
)

func AssertErrorMessage(t *testing.T, pExpectedErrorMessage string, pError error, pFunctionName string) {
	if pError == nil {
		t.Errorf("%s should return not nil error", pFunctionName)
		return
	}
	if pError.Error() != pExpectedErrorMessage {
		t.Errorf(
			"%s error message expected [%s], but returned [%s]",
			pFunctionName,
			pExpectedErrorMessage,
			pError.Error(),
		)
	}
}

func AssertNotNilResultAndNilError[R any](t *testing.T, pResult *R, pError error, pFunctionName string) bool {
	result := true
	if pResult == nil {
		t.Errorf("%s cannot return nil Result", pFunctionName)
		result = false
	}
	if !AssertNilError(t, pError, pFunctionName) {
		result = false
	}
	return result
}

func AssertNilResultAndNotNilError[R any](t *testing.T, pResult *R, pError error, pFunctionName string) bool {
	result := true
	if pResult != nil {
		t.Errorf("%s should return nil Result", pFunctionName)
		result = false
	}
	if !AssertNotNilError(t, pError, pFunctionName) {
		result = false
	}
	return result
}

func AssertNilError(t *testing.T, pError error, pFunctionName string) bool {
	if pError == nil {
		return true
	}
	t.Errorf("%s return an error: %s", pFunctionName, pError.Error())
	return false
}

func AssertNotNilError(t *testing.T, pError error, pFunctionName string) bool {
	if pError != nil {
		return true
	}
	t.Errorf("%s should return not nil error", pFunctionName)
	return false
}

func AssertEqualWithLabel[T comparable](t *testing.T, pExpected, pReturned T, pLabel string) bool {
	if pExpected == pReturned {
		return true
	}
	if pLabel != "" && !strings.HasSuffix(pLabel, ": ") {
		pLabel += ": "
	}
	t.Errorf("%sexpected [%v], but returned [%v]", pLabel, pExpected, pReturned)
	return false
}

func AssertNamedArgEqualWithLabel(t *testing.T, pExpected any, pReturned sql.NamedArg, pLabel string) bool {
	if pExpected == pReturned {
		return true
	}
	pExpectedNamedArg := pExpected.(sql.NamedArg)
	if pExpectedNamedArg.Name == pReturned.Name && pExpectedNamedArg.Value == pReturned.Value {
		return true
	}

	if pLabel != "" && !strings.HasSuffix(pLabel, ": ") {
		pLabel += ": "
	}
	t.Errorf("%sexpected [%v], but returned [%v]", pLabel, pExpected, pReturned)
	return false
}
