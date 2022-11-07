package rds

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

func TestExpandParameters(t *testing.T) {
	expanded := []interface{}{
		map[string]interface{}{
			"name":         "character_set_client",
			"value":        "utf8",
			"apply_method": "immediate",
		},
	}
	parameters := ExpandParameters(expanded)

	expected := &rds.Parameter{
		ParameterName:  aws.String("character_set_client"),
		ParameterValue: aws.String("utf8"),
		ApplyMethod:    aws.String("immediate"),
	}

	if !reflect.DeepEqual(parameters[0], expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			parameters[0],
			expected)
	}
}

func TestFlattenParameters(t *testing.T) {
	cases := []struct {
		Input  []*rds.Parameter
		Output []map[string]interface{}
	}{
		{
			Input: []*rds.Parameter{
				{
					ParameterName:  aws.String("character_set_client"),
					ParameterValue: aws.String("utf8"),
				},
			},
			Output: []map[string]interface{}{
				{
					"name":  "character_set_client",
					"value": "utf8",
				},
			},
		},
		{
			Input: []*rds.Parameter{
				{
					ParameterName:  aws.String("character_set_client"),
					ParameterValue: aws.String("utf8"),
					ApplyMethod:    aws.String("immediate"),
				},
			},
			Output: []map[string]interface{}{
				{
					"name":         "character_set_client",
					"value":        "utf8",
					"apply_method": "immediate",
				},
			},
		},
	}

	for _, tc := range cases {
		output := FlattenParameters(tc.Input)
		if !reflect.DeepEqual(output, tc.Output) {
			t.Fatalf("Got:\n\n%#v\n\nExpected:\n\n%#v", output, tc.Output)
		}
	}
}
