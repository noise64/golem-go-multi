# Golem Go Example with Multiple Components and Worker to Worker RPC Communication

## Building
The project uses [magefile](https://magefile.org/) for building. Either install the tool binary `mage`,
or use the __zero install option__: `go run mage.go`. This readme will use the latter.

To see the available commands use:

```shell
go run mage.go
Targets:
  addStubDependency             adds generated and built stub dependency to componentGolemCliAddStubDependency
  build                         alias for BuildAllComponents
  buildAllComponents            builds all components
  buildComponent                builds component by name
  buildStubComponent            builds RPC stub for component
  clean                         cleans the projects
  generateBinding               generates go binding from WIT
  generateNewComponent          generates a new component based on the component-template
  stubCompose                   composes dependencies
  tinyGoBuildComponentBinary    build wasm component binary with tiny go
  updateRpcStubs                builds rpc stub components and adds them as dependency
  wasmToolsComponentEmbed       embeds type info into wasm component with wasm-tools
  wasmToolsComponentNew         create golem component with wasm-tools
```

For building the project for the first time (or after `clean`) use the following commands:

```shell
go run mage.go updateRpcStubs
go run mage.go build
```

After this, using the `build` command is enough, unless there are changes in the RPC dependencies,
in that case `updateRpcStubs` is needed again.

The final components that are usable by golem are placed in the `target/components` folder.

## Adding Components

Use the `generateNewComponent` command to add new components to the project:

```shell
go run mage.go generateNewComponent component-four
```

The above will create a new component in the `components/component-four` directory based on the template at [/component-template/component](/component-template/component).

After adding a new component the `build` command will also include it.

## Using Worker to Worker RPC calls

The dependencies between components are defined in  the [/magefiles/magefile.go](/magefiles/magefile.go) build script:

```go
// componentDeps defines the Worker to Worker RPC dependencies
var componentDeps = map[string][]string{
    "component-one": {"component-two", "component-three"},
    "component-two": {"component-three"},
}
```

After changing dependencies the `updateRpcStubs` command can be used to create the necessary stubs:

```shell
go run mage.go updateRpcStubs
```

The command will create stubs for the dependency projects in the ``/target/stub`` directory and will also place the required stub _WIT_ interfaces on the dependant component's `wit/deps` directory.

To actually use the dependencies in a project it also has to be manually imported in the component's world.

E.g. with the above definitions the following import has to be __manually__ added to `/components/component-one/wit/component-one.wit`:

```wit
import golem:component-two-stub;
import golem:component-three-stub;
```

So the component definition should like similar to this:

```wit
package golem:component-one;

// See https://component-model.bytecodealliance.org/design/wit.html for more details about the WIT syntax

interface component-one-api {
  add: func(value: u64);
  get: func() -> u64;
}

world component-one {
  // Golem dependencies
  import golem:api/host@0.2.0;
  import golem:rpc/types@0.1.0;

  // WASI dependencies
  import wasi:blobstore/blobstore;
  // .
  // .
  // .
  // other dependencies
  import wasi:sockets/instance-network@0.2.0;

  // Project Component dependencies
  import golem:component-two-stub;
  import golem:component-three-stub;

  export component-one-api;
}
```

After this `build` (or the `generateBinding`) command can be used to update bindings, which now should include the
required functions for calling other components.

Here's an example that delegates the `Add` call to another component and waits for the result:

```go
import (
	"github.com/golemcloud/golem-go/std"

	"golem-go-project/components/component-one/binding"
)


func (i *Impl) Add(value uint64) {
    std.Init(std.Packages{Os: true, NetHttp: true})
    
    componentTwo := binding.NewComponentTwoApi(binding.GolemRpc0_1_0_TypesUri{Value: "uri"})
    defer componentTwo.Drop()
    componentTwo.BlockingAdd(value)

    i.counter += value
}

```


