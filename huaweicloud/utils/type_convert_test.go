package utils_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestTypeConvertFunc_StringToJson(t *testing.T) {
	var (
		emptyInput     = "{}"
		correctInput   = "{\"foo\":\"bar\"}"
		incorrectInput = `func() {
			fmt.Println("Hello, this is a function!")
		}`
		emptyInputExpected   = make(map[string]interface{})
		correctInputExpected = map[string]interface{}{
			"foo": "bar",
		}
	)

	testOutput := utils.StringToJson(emptyInput)
	if !reflect.DeepEqual(testOutput, emptyInputExpected) {
		t.Fatalf("The processing result of the StringToJson method is not as expected, want %s, but got %s",
			utils.Green(emptyInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.StringToJson(correctInput)
	if !reflect.DeepEqual(testOutput, correctInputExpected) {
		t.Fatalf("The processing result of the StringToJson method is not as expected, want %s, but got %s",
			utils.Green(correctInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.StringToJson(incorrectInput)
	if !reflect.DeepEqual(testOutput, make(map[string]interface{})) {
		t.Fatalf("The processing result of the StringToJson method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the JsonToString method meets expectation")
}

func TestTypeConvertFunc_JsonToString(t *testing.T) {
	type Test struct {
		Foo string `json:"foo,omitempty"`
	}

	var (
		emptyInput   = Test{}
		correctInput = Test{
			Foo: "bar",
		}
		emptyInputExpected   = "{}"
		correctInputExpected = "{\"foo\":\"bar\"}"
		// Function is an unsupported type for JsonToString() function input and an error will be returned.
		functionInput = func() {
			fmt.Println("Hello, this is a function!")
		}

		inputWithNestedKeys = map[string]interface{}{
			"owner": "utils/function",
			"parameters": map[string]interface{}{
				"id":   "4c8374e0-8632-4151-80e8-374b81f80f3e",
				"name": "example",
				"age":  18,
			},
			"usage": "value conversion",
		}
		// Test whether both nested keys and normal keys can be removed correctly.
		nestedKeys                = []string{"parameters.id", "usage"}
		expectedWithoutNestedKeys = "{\"owner\":\"utils/function\",\"parameters\":{\"age\":18,\"name\":\"example\"}}"
	)

	testOutput := utils.JsonToString(emptyInput)
	if !reflect.DeepEqual(testOutput, emptyInputExpected) {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want %s, but got %s",
			utils.Green(emptyInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(correctInput)
	if !reflect.DeepEqual(testOutput, correctInputExpected) {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want %s, but got %s",
			utils.Green(correctInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(functionInput)
	if !reflect.DeepEqual(testOutput, "") {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(nil)
	if !reflect.DeepEqual(testOutput, "") {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(inputWithNestedKeys, nestedKeys...)
	if !reflect.DeepEqual(testOutput, expectedWithoutNestedKeys) {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want %s, but got %s",
			utils.Green(expectedWithoutNestedKeys), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the JsonToString method meets expectation")
}
