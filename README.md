# Golem Go Example with Multiple Components and Worker to Worker RPC Communication

## Building
The project uses [magefile](https://magefile.org/) for building. Either install the tool binary `mage`,
or use the __zero install option__: `go run mage.go`. This readme will use the latter.

To see the available commands use:

```shell
go run mage.go
Targets:
  build                         alias for BuildAllComponents
  buildAllComponents            builds all components
  buildComponent                builds component by name
  clean                         cleans the projects
  generateBinding               generates go binding from WIT
  golemCliAddStubDependency     adds generated and built stub dependency to componentGolemCliAddStubDependency
  golemCliBuildStubComponent    builds RPC stub for component
  golemCliStubCompose           composes dependencies
  tinyGoBuildComponentBinary    build wasm component binary with tiny go
  updateRpcStubs                builds rpc stub components and adds them as dependency
  wasmToolsComponentEmbed       embeds type info into wasm component with wasm-tools
  wasmToolsComponentNew         create golem component with wasm-tools and compose dependencies
```

For building the project the first time use the following commands:
```shell
go run mage.go updateRpcStubs
go run mage.go build`
```

After this, using the `build` command is enough, unless there are changes in the RPC dependencies,
in that case `updateRpcStubs` is needed again.

## Adding Components
TODO

## Defining and updating RPC dependencies
TODO

[/magefiles/magefile.go](/magefiles/magefile.go)
```go
// components and RPC dependencies
var components = map[string][]string{
	"component1": {"component2", "component3"},
	"component2": {"component3"},
	"component3": nil,
}
```
