package assert

import (
    "fmt"
    "os"
    "runtime/debug"
)

type AssertedArgs interface{}

var assertData map[string]AssertedArgs = map[string]AssertedArgs{}

func runAssert(msg string, args ...interface{}) {
	assertions := []interface{}{
		"msg",
		msg,
		"area",
		"Assert",
	}

	assertions = append(assertions, args...)

    // ya idk w.e
    for i := 0; i < 4; i++ {
        fmt.Println()
    }

	fmt.Fprintf(os.Stderr, "            [ASSERTION FAILURES]\n")
	fmt.Fprintf(os.Stderr, "[ARGS]: %+v\n", args)

	for k, v := range assertData {
		assertions = append(assertions, k, v)
	}

	fmt.Fprintf(os.Stderr, "[ASSERTION]\n")
	for i := 0; i < len(assertions); i += 2 {
		fmt.Fprintf(os.Stderr, "    %s=%v\n", assertions[i], assertions[i+1])
	}
	fmt.Fprint(os.Stderr, string(debug.Stack()))
	os.Exit(1)
}

func Assert(cond bool, msg string, data ...any) {
	if !cond {
		runAssert(msg, data...)
	}
}

func NoError(err error, msg string, data ...any) {
	if err != nil {
		data = append(data, "error", msg)
		runAssert(msg, data...)
	}
}

func NotNil(val any, msg string, data ...any) {
	if val == nil {
		runAssert(msg, data...)
	}
}
