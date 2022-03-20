package main


import (
	"github.com/0qq/hello-world-go-http-example/pkg/hw"
)


func main() {
	s := hw.NewServer()
	go s.EmulateActivity()
	s.Start()
}
