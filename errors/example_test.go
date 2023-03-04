package errors_test

import (
	"fmt"

	"github.com/KEINOS/go-errors/errors"
)

func Example_wrapping() {
	innerFn := func() error {
		return errors.New("error of innerFn()")
	}

	outerFn := func() error {
		err := innerFn()

		// Wrap returns nil if err is nil.
		return errors.Wrap(err, "error of outerFn()")
	}

	if err := outerFn(); err != nil {
		// For verbose/errorstack output, use fmt.Printf("%+v", err)
		fmt.Println(err)
	}
	// Output: error of outerFn(): error of innerFn()
}

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")

	fmt.Printf("%+v", err)
	// Example output:
	// whoops
	// github.com/KEINOS/go-errors/errors_test.ExampleNew_printf
	// 	/Users/path/to/repo/GitHub/KEINOS/go-errors/errors/example_test.go:62
	// testing.runExample
	// 	/usr/local/Cellar/go/1.20.1/libexec/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/Cellar/go/1.20.1/libexec/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/Cellar/go/1.20.1/libexec/src/testing/testing.go:1908
	// main.main
	// 	_testmain.go:137
	// runtime.main
	// 	/usr/local/Cellar/go/1.20.1/libexec/src/runtime/proc.go:250
	// runtime.goexit
	// 	/usr/local/Cellar/go/1.20.1/libexec/src/runtime/asm_amd64.s:1598
}

func ExampleWithMessage() {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")

	fmt.Println(err)
	// Output: oh noes: whoops
}

func ExampleWithStack() {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)

	fmt.Println(err)
	// Output: whoops
}

func ExampleWithStack_printf() {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)

	fmt.Printf("%+v", err)
	// Example Output:
	// whoops
	// github.com/KEINOS/go-errors/errors_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/KEINOS/go-errors/errors/example_test.go:55
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/KEINOS/go-errors/errors/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
	// github.com/KEINOS/go-errors/errors_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/KEINOS/go-errors/errors/example_test.go:56
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/KEINOS/go-errors/errors/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")

	fmt.Println(err)
	// Output: oh noes: whoops
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")

	return errors.Wrap(e3, "outer")
}

func ExampleCause() {
	err := fn()

	fmt.Println(err)
	fmt.Println(errors.Cause(err))
	// Output:
	// outer: middle: inner: error
	// error
}

func ExampleWrap_extended() {
	err := fn()

	fmt.Printf("%+v\n", err)
	// Example output:
	// error
	// github.com/KEINOS/go-errors/errors_test.fn
	//         /home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:47
	// github.com/KEINOS/go-errors/errors_test.ExampleCause_printf
	//         /home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:63
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/KEINOS/go-errors/errors/_test/_testmain.go:104
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
	// github.com/KEINOS/go-errors/errors_test.fn
	// 	  /home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:48: inner
	// github.com/KEINOS/go-errors/errors_test.fn
	//        /home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:49: middle
	// github.com/KEINOS/go-errors/errors_test.fn
	//      /home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:50: outer
}

func ExampleWrapf() {
	cause := errors.New("whoops")
	err := errors.Wrapf(cause, "oh noes #%d", 2)

	fmt.Println(err)
	// Output: oh noes #2: whoops
}

func ExampleErrorf_extended() {
	err := errors.Errorf("whoops: %s", "foo")

	fmt.Printf("%+v", err)
	// Example output:
	// whoops: foo
	// github.com/KEINOS/go-errors/errors_test.ExampleErrorf
	//         /home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:101
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/KEINOS/go-errors/errors/_test/_testmain.go:102
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func Example_stackTrace() {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	err, ok := errors.Cause(fn()).(stackTracer)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	st := err.StackTrace()

	fmt.Printf("%+v", st[0:2]) // top two frames
	// Example output:
	// github.com/KEINOS/go-errors/errors_test.fn
	//	/home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:47
	// github.com/KEINOS/go-errors/errors_test.Example_stackTrace
	//	/home/dfc/src/github.com/KEINOS/go-errors/errors/example_test.go:127
}

func ExampleCause_printf() {
	err := errors.Wrap(func() error {
		return func() error {
			return errors.New("hello world")
		}()
	}(), "failed")

	fmt.Printf("%v", err)
	// Output: failed: hello world
}
