package main

import (
	"github.com/golemcloud/golem-go/std"
	"golem-go-project/components/component2/binding"
)

func init() {
	binding.SetExportsGolemComponent2Component2Api(&Impl{})
}

type Impl struct {
	counter uint64
}

func (i *Impl) Add(value uint64) {
	std.Init(std.Packages{Os: true, NetHttp: true})

	component3 := binding.NewComponent3Api(binding.GolemRpc0_1_0_TypesUri{Value: "uri"})
	defer component3.Drop()
	component3.BlockingAdd(value)

	i.counter += value
}

func (i *Impl) Get() uint64 {
	std.Init(std.Packages{Os: true, NetHttp: true})

	return i.counter
}

func main() {}
