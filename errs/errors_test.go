package errs_test

import (
	"testing"

	"github.com/ordershift/ormshift/errs"
	"github.com/ordershift/ormshift/internal/testutils"
)

type errorTester struct {
	expectedMessageError string
	expectedTypeError    error
	testedError          error
}

func TestErrors(t *testing.T) {
	testers := []errorTester{
		{
			expectedMessageError: "invalid driver",
			expectedTypeError:    errs.ErrInvalid,
			testedError:          errs.Invalid("driver"),
		},
		{
			expectedMessageError: "databasedriver cannot be nil",
			expectedTypeError:    errs.ErrNil,
			testedError:          errs.Nil("databasedriver"),
		},
		{
			expectedMessageError: "failed to get db schema",
			expectedTypeError:    errs.ErrFailedTo,
			testedError:          errs.FailedTo("get db schema"),
		},
		{
			expectedMessageError: `on open database: failed to get db schema`,
			expectedTypeError:    errs.ErrFailedTo,
			testedError:          errs.WithContext("on open database", errs.FailedTo("get db schema")),
		},
	}
	for _, tester := range testers {
		testutils.AssertErrorType(t, tester.expectedTypeError, tester.testedError)
		testutils.AssertErrorMessage(t, tester.expectedMessageError, tester.testedError, "errs pkg")
	}
}
