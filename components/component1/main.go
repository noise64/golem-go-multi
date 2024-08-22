package main

import (
	"github.com/golemcloud/golem-go/std"
	"golem-go-project/components/component1/binding"
)

func init() {
	binding.SetExportsGolemComponent1Component1Api(&Impl{})
}

type Impl struct {
	counter uint64
}

func (i *Impl) Add(value uint64) {
	std.Init(std.Packages{Os: true, NetHttp: true})

	component2 := binding.NewComponent2Api(binding.GolemRpc0_1_0_TypesUri{Value: "uri"})
	defer component2.Drop()

	component3 := binding.NewComponent3Api(binding.GolemRpc0_1_0_TypesUri{Value: "uri"})
	defer component3.Drop()

	component2.BlockingAdd(value)
	component3.BlockingAdd(value)

	i.counter += value
}

func (i *Impl) Get() uint64 {
	std.Init(std.Packages{Os: true, NetHttp: true})

	return i.counter
}

func main() {}
