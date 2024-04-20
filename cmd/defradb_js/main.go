//go:build js

package main

import (
	"context"
	"syscall/js"

	defradbJS "github.com/sourcenetwork/defradb/js"
)

func main() {
	ctx := context.Background()

	js.Global().Set("defraDB", map[string]any{
		"open": defradbJS.Open(ctx),
	})

	<-ctx.Done()
}
