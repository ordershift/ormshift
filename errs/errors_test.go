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
			expectedMessageError: "database driver cannot be nil",
			expectedTypeError:    errs.ErrNil,
			testedError:          errs.Nil("database driver"),
		},
		{
			expectedMessageError: "failed to get db schema",
			expectedTypeError:    errs.ErrFailedTo,
			testedError:          errs.FailedTo("get db schema", nil),
		},
		{
			expectedMessageError: "failed to get db schema: db cannot be nil",
			expectedTypeError:    errs.ErrFailedTo,
			testedError:          errs.FailedTo("get db schema", errs.Nil("db")),
		},
	}
	for _, tester := range testers {
		testutils.AssertErrorType(t, tester.expectedTypeError, tester.testedError)
		testutils.AssertErrorMessage(t, tester.expectedMessageError, tester.testedError, "errs pkg")
	}
}
