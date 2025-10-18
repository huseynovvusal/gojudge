package main

import (
	"fmt"
	"huseynovvusal/gojudge/internal/judge"
)

func main() {
	// Python code to run a loop 10^5 times
	// pythonCode := "for i in range(10**7):\n\tif i % 10000 == 0: print(i)"
	// result, err := judge.RunCode("python", pythonCode, 1, 256, 1)
	// C code to print numbers from 0 to 10^5
	cCode := "#include <stdio.h>\nint main() {\n    for (int i = 0; i < 100000; i++) {\n        if (i % 10000 == 0) printf(\"%d\\n\", i);\n    }\n    return 0;\n}"
	result, err := judge.RunCode("c", cCode, 2, 256, 1)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	println("Output:\n" + result.Output)
	println("Execution Time (ms):", result.ExecutionMs)
}
