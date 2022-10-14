// +build tools

package internal

// These tools required by ginkgo
import (
	_ "github.com/go-task/slim-sprig"
	_ "github.com/google/pprof/profile"
	_ "golang.org/x/tools/go/ast/inspector"
)
