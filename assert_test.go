package ormshift_test

import (
	"strings"
	"testing"
)

func assertErrorMessage(t *testing.T, pExpectedErrorMessage string, pError error, pFunctionName string) {
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

func assertNotNilResultAndNilError[R any](t *testing.T, pResult *R, pError error, pFunctionName string) bool {
	lResult := true
	if pResult == nil {
		t.Errorf("%s cannot return nil Result", pFunctionName)
		lResult = false
	}
	if !assertNilError(t, pError, pFunctionName) {
		lResult = false
	}
	return lResult
}

func assertNilResultAndNotNilError[R any](t *testing.T, pResult *R, pError error, pFunctionName string) bool {
	lResult := true
	if pResult != nil {
		t.Errorf("%s should return nil Result", pFunctionName)
		lResult = false
	}
	if !assertNotNilError(t, pError, pFunctionName) {
		lResult = false
	}
	return lResult
}

func assertNilError(t *testing.T, pError error, pFunctionName string) bool {
	if pError == nil {
		return true
	}
	t.Errorf("%s return an error: %s", pFunctionName, pError.Error())
	return false
}

func assertNotNilError(t *testing.T, pError error, pFunctionName string) bool {
	if pError != nil {
		return true
	}
	t.Errorf("%s should return not nil error", pFunctionName)
	return false
}

// func assertEqual[T comparable](t *testing.T, pExpected, pReturned T) bool {
// 	return assertEqualWithLabel(t, pExpected, pReturned, "")
// }

func assertEqualWithLabel[T comparable](t *testing.T, pExpected, pReturned T, pLabel string) bool {
	if pExpected == pReturned {
		return true
	}
	if pLabel != "" && !strings.HasSuffix(pLabel, ": ") {
		pLabel += ": "
	}
	t.Errorf("%sexpected [%v], but returned [%v]", pLabel, pExpected, pReturned)
	return false
}

// func assertNotEqual[T comparable](t *testing.T, pNotExpected, pReturned T) bool {
// 	return assertNotEqualWithLabel(t, pNotExpected, pReturned, "")
// }

// func assertNotEqualWithLabel[T comparable](t *testing.T, pNotExpected, pReturned T, pLabel string) bool {
// 	if pNotExpected != pReturned {
// 		return true
// 	}
// 	if pLabel != "" && !strings.HasSuffix(pLabel, ": ") {
// 		pLabel += ": "
// 	}
// 	t.Errorf("%s[%v] was not expected", pLabel, pReturned)
// 	return false
// }
