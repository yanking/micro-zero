package main

import (
	"github.com/yanking/micro-zero/internal/apiserver"
	"math/rand"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	apiserver.NewApp("apiserver").Run()
}
