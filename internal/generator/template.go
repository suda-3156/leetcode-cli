package generator

var templateMap = map[string]string{
	"golang":  golang,
	"python":  python,
	"python3": python,
	"default": defaultTmpl,
}

const golang = `{{ .Header }}
package main

{{ .CodeSnippet }}

func main() {
	testCases := []struct {
		// Define your test case structure here
		input struct {}
		want  string	
	}{
		// Add your test cases here
	}

	for _, tc := range testCases {
		result := {{ .FunctionName }}(tc.input)
		if result != tc.want {
			fmt.Printf("Test failed for input %v: got %v, want %v\n", tc.input, result, tc.want)
		} else {
			fmt.Printf("Test passed for input %v\n", tc.input)
		}
	}
}
`

const python = `{{ .Header }}
from typing import List, Any, TypedDict


{{ .CodeSnippet }}

TestCase = TypedDict("TestCase", {"input": Any, "want": Any})


def main():
    test_cases: List[TestCase] = [
		# Add your test cases here
        {"input": {}, "want": ""}
    ]

    s = Solution()

    for tc in test_cases:
        result = s.{{ .FunctionName }}(tc["input"])
        if result != tc["want"]:
            print(f"Test failed for input {tc["input"]}: got {result}, want {tc["want"]}")
        else:
            print(f"Test passed for input {tc["input"]}")


if __name__ == "__main__":
    main()
`

const defaultTmpl = `{{ .Header }}
{{ .CodeSnippet }}
`
