package main

import (
	"fmt"
	"huseynovvusal/gojudge/internal/judge"
)

func main() {
	// Python code to run a loop 10^5 times
	// pythonCode := "n = int(input())\nfor i in range(n):\n\tprint(i)"
	// result, err := judge.RunCode("python", pythonCode, "1000", 2, 256, 1)
	// C code to run a loop 10^5 times
	cCode := "#include <stdio.h>\nint main() {\n\tint n;\n\tscanf(\"%d\", &n);\n\tfor (int i = 0; i < n; i++) {\n\t\tprintf(\"%d\\n\", i);\n\t}\n\treturn 0;\n}"
	result, err := judge.RunCode("c", cCode, "1000", 2, 256, 1)
	// C++ code to run a loop 10^5 times
	// cppCode := "#include <iostream>\nusing namespace std;\nint main() {\n\tint n;\n\tcin >> n;\n\tfor (int i = 0; i < n; i++) {\n\t\tcout << i << endl;\n\t}\n\treturn 0;\n}"
	// result, err := judge.RunCode("cpp", cppCode, "1000",
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	println("Output:\n" + result.Output)
	println("Execution Time (ms):", result.ExecutionMs)
}
