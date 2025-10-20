package main

import (
	"context"
	"fmt"
	"huseynovvusal/gojudge/internal/executor"
	pb "huseynovvusal/gojudge/internal/proto"
	"net"

	"google.golang.org/grpc"
)

// func main() {
// 	// Python code to run a loop 10^5 times
// 	// pythonCode := "n = int(1e3)\nfor i in range(n):\n\tprint(i)"
// 	// result, err := judge.RunCode("python", pythonCode, "", 2, 256, 1)

// 	// C code to run a loop 10^5 times
// 	// cCode := "#include <stdio.h>\nint main() {\n\tint n;\n\tscanf(\"%d\", &n);\n\tfor (int i = 0; i < n; i++) {\n\t\tprintf(\"%d\\n\", i);\n\t}\n\treturn 0;\n}"
// 	// result, err := judge.RunCode("c", cCode, "1000", 2, 256, 1)

// 	// C++ code to run a loop 10^5 times
// 	// cppCode := "#include <iostream>\nusing namespace std;\nint main() {\n\tint n;\n\tcin >> n;\n\tfor (int i = 0; i < n; i++) {\n\t\tcout << i << endl;\n\t}\n\treturn 0;\n}"
// 	// result, err := judge.RunCode("cpp", cppCode, "1000",

// 	goCode := "package main\nimport \"fmt\"\nfunc main() {\n\tn := 1000\n\tfor i := 0; i < n; i++ {\n\t\tfmt.Println(i)\n\t}\n}"
// 	result, err := judge.RunCode("go", goCode, "", 2, 1024, 1)

// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}
// 	println("Output:\n" + result.Output)
// 	println("Execution Time (ms):", result.ExecutionMs)
// }

type server struct {
	pb.UnimplementedExecutorServiceServer
}

func (*server) Execute(ctx context.Context, req *pb.ExecutorRequest) (*pb.ExecutorResponse, error) {
	result, err := executor.RunCode(req.Language, req.Code, req.Input, int16(req.TimeLimit), int16(req.MemoryLimit), int16(req.CpuLimit))

	if err != nil {
		return &pb.ExecutorResponse{}, err
	}

	return &pb.ExecutorResponse{
		Output:        result.Output,
		ExecutionTime: result.ExecutionMs,
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterExecutorServiceServer(s, &server{})

	fmt.Println("Server is running on port :50051")

	if err := s.Serve(lis); err != nil {
		fmt.Println(err)
	}

}
