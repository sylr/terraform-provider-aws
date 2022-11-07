package tfresource_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestRetryWhenAWSErrCodeEquals(t *testing.T) { // nosemgrep:ci.aws-in-func-name
	var retryCount int32

	testCases := []struct {
		Name        string
		F           func() (interface{}, error)
		ExpectError bool
	}{
		{
			Name: "no error",
			F: func() (interface{}, error) {
				return nil, nil
			},
		},
		{
			Name: "non-retryable other error",
			F: func() (interface{}, error) {
				return nil, errors.New("TestCode")
			},
			ExpectError: true,
		},
		{
			Name: "non-retryable AWS error",
			F: func() (interface{}, error) {
				return nil, awserr.New("Testing", "Testing", nil)
			},
			ExpectError: true,
		},
		{
			Name: "retryable AWS error timeout",
			F: func() (interface{}, error) {
				return nil, awserr.New("TestCode1", "TestMessage", nil)
			},
			ExpectError: true,
		},
		{
			Name: "retryable AWS error success",
			F: func() (interface{}, error) {
				if atomic.CompareAndSwapInt32(&retryCount, 0, 1) {
					return nil, awserr.New("TestCode2", "TestMessage", nil)
				}

				return nil, nil
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			retryCount = 0

			_, err := tfresource.RetryWhenAWSErrCodeEquals(5*time.Second, testCase.F, "TestCode1", "TestCode2")

			if testCase.ExpectError && err == nil {
				t.Fatal("expected error")
			} else if !testCase.ExpectError && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestRetryWhenAWSErrMessageContains(t *testing.T) { // nosemgrep:ci.aws-in-func-name
	var retryCount int32

	testCases := []struct {
		Name        string
		F           func() (interface{}, error)
		ExpectError bool
	}{
		{
			Name: "no error",
			F: func() (interface{}, error) {
				return nil, nil
			},
		},
		{
			Name: "non-retryable other error",
			F: func() (interface{}, error) {
				return nil, errors.New("TestCode")
			},
			ExpectError: true,
		},
		{
			Name: "non-retryable AWS error",
			F: func() (interface{}, error) {
				return nil, awserr.New("TestCode1", "Testing", nil)
			},
			ExpectError: true,
		},
		{
			Name: "retryable AWS error timeout",
			F: func() (interface{}, error) {
				return nil, awserr.New("TestCode1", "TestMessage1", nil)
			},
			ExpectError: true,
		},
		{
			Name: "retryable AWS error success",
			F: func() (interface{}, error) {
				if atomic.CompareAndSwapInt32(&retryCount, 0, 1) {
					return nil, awserr.New("TestCode1", "TestMessage1", nil)
				}

				return nil, nil
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			retryCount = 0

			_, err := tfresource.RetryWhenAWSErrMessageContains(5*time.Second, testCase.F, "TestCode1", "TestMessage1")

			if testCase.ExpectError && err == nil {
				t.Fatal("expected error")
			} else if !testCase.ExpectError && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestRetryWhenNewResourceNotFound(t *testing.T) {
	var retryCount int32

	testCases := []struct {
		Name        string
		F           func() (interface{}, error)
		NewResource bool
		ExpectError bool
	}{
		{
			Name: "no error",
			F: func() (interface{}, error) {
				return nil, nil
			},
		},
		{
			Name: "no error new resource",
			F: func() (interface{}, error) {
				return nil, nil
			},
			NewResource: true,
		},
		{
			Name: "non-retryable other error",
			F: func() (interface{}, error) {
				return nil, errors.New("TestCode")
			},
			ExpectError: true,
		},
		{
			Name: "non-retryable other error new resource",
			F: func() (interface{}, error) {
				return nil, errors.New("TestCode")
			},
			NewResource: true,
			ExpectError: true,
		},
		{
			Name: "non-retryable AWS error",
			F: func() (interface{}, error) {
				return nil, awserr.New("Testing", "Testing", nil)
			},
			ExpectError: true,
		},
		{
			Name: "retryable NotFoundError not new resource",
			F: func() (interface{}, error) {
				return nil, &resource.NotFoundError{}
			},
			ExpectError: true,
		},
		{
			Name: "retryable NotFoundError new resource timeout",
			F: func() (interface{}, error) {
				return nil, &resource.NotFoundError{}
			},
			NewResource: true,
			ExpectError: true,
		},
		{
			Name: "retryable NotFoundError success new resource",
			F: func() (interface{}, error) {
				if atomic.CompareAndSwapInt32(&retryCount, 0, 1) {
					return nil, &resource.NotFoundError{}
				}

				return nil, nil
			},
			NewResource: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			retryCount = 0

			_, err := tfresource.RetryWhenNotFound(5*time.Second, testCase.F)

			if testCase.ExpectError && err == nil {
				t.Fatal("expected error")
			} else if !testCase.ExpectError && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestRetryWhenNotFound(t *testing.T) {
	var retryCount int32

	testCases := []struct {
		Name        string
		F           func() (interface{}, error)
		ExpectError bool
	}{
		{
			Name: "no error",
			F: func() (interface{}, error) {
				return nil, nil
			},
		},
		{
			Name: "non-retryable other error",
			F: func() (interface{}, error) {
				return nil, errors.New("TestCode")
			},
			ExpectError: true,
		},
		{
			Name: "non-retryable AWS error",
			F: func() (interface{}, error) {
				return nil, awserr.New("Testing", "Testing", nil)
			},
			ExpectError: true,
		},
		{
			Name: "retryable NotFoundError timeout",
			F: func() (interface{}, error) {
				return nil, &resource.NotFoundError{}
			},
			ExpectError: true,
		},
		{
			Name: "retryable NotFoundError success",
			F: func() (interface{}, error) {
				if atomic.CompareAndSwapInt32(&retryCount, 0, 1) {
					return nil, &resource.NotFoundError{}
				}

				return nil, nil
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			retryCount = 0

			_, err := tfresource.RetryWhenNotFound(5*time.Second, testCase.F)

			if testCase.ExpectError && err == nil {
				t.Fatal("expected error")
			} else if !testCase.ExpectError && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestRetryUntilNotFound(t *testing.T) {
	var retryCount int32

	testCases := []struct {
		Name        string
		F           func() (interface{}, error)
		ExpectError bool
	}{
		{
			Name: "no error",
			F: func() (interface{}, error) {
				return nil, nil
			},
			ExpectError: true,
		},
		{
			Name: "other error",
			F: func() (interface{}, error) {
				return nil, errors.New("TestCode")
			},
			ExpectError: true,
		},
		{
			Name: "AWS error",
			F: func() (interface{}, error) {
				return nil, awserr.New("Testing", "Testing", nil)
			},
			ExpectError: true,
		},
		{
			Name: "NotFoundError",
			F: func() (interface{}, error) {
				return nil, &resource.NotFoundError{}
			},
		},
		{
			Name: "retryable NotFoundError",
			F: func() (interface{}, error) {
				if atomic.CompareAndSwapInt32(&retryCount, 0, 1) {
					return nil, nil
				}

				return nil, &resource.NotFoundError{}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			retryCount = 0

			_, err := tfresource.RetryUntilNotFound(5*time.Second, testCase.F)

			if testCase.ExpectError && err == nil {
				t.Fatal("expected error")
			} else if !testCase.ExpectError && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestRetryContext_error(t *testing.T) {
	t.Parallel()

	expected := fmt.Errorf("nope")
	f := func() *resource.RetryError {
		return resource.NonRetryableError(expected)
	}

	errCh := make(chan error)
	go func() {
		errCh <- tfresource.RetryContext(context.Background(), 1*time.Second, f)
	}()

	select {
	case err := <-errCh:
		if err != expected {
			t.Fatalf("bad: %#v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOptionsApply(t *testing.T) {
	testCases := map[string]struct {
		options  tfresource.Options
		expected resource.StateChangeConf
	}{
		"Nothing": {
			options:  tfresource.Options{},
			expected: resource.StateChangeConf{},
		},
		"Delay": {
			options: tfresource.Options{
				Delay: 1 * time.Minute,
			},
			expected: resource.StateChangeConf{
				Delay: 1 * time.Minute,
			},
		},
		"MinPollInterval": {
			options: tfresource.Options{
				MinPollInterval: 1 * time.Minute,
			},
			expected: resource.StateChangeConf{
				MinTimeout: 1 * time.Minute,
			},
		},
		"PollInterval": {
			options: tfresource.Options{
				PollInterval: 1 * time.Minute,
			},
			expected: resource.StateChangeConf{
				PollInterval: 1 * time.Minute,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			conf := resource.StateChangeConf{}

			testCase.options.Apply(&conf)

			if a, e := conf.Delay, testCase.expected.Delay; a != e {
				t.Errorf("Delay: expected %s, got %s", e, a)
			}
			if a, e := conf.MinTimeout, testCase.expected.MinTimeout; a != e {
				t.Errorf("MinTimeout: expected %s, got %s", e, a)
			}
			if a, e := conf.PollInterval, testCase.expected.PollInterval; a != e {
				t.Errorf("PollInterval: expected %s, got %s", e, a)
			}
		})
	}
}
