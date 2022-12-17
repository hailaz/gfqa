package main

import (
	_ "gfqa/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"gfqa/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
