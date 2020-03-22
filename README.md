# MGA: Modern Go Application tool

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sagikazarmark/mga/CI?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/sagikazarmark.dev/mga?style=flat-square)](https://goreportcard.com/report/sagikazarmark.dev/mga)

**Go application development tool applying modern practices.**

The `mga` tool is an experimental collection of tools and practices related to developing Go applications.
Its goal is to make the maintenance of applications, thus the life of developers easier.

**⚠️ This tool is still under heavy development! Things may change. ⚠️**

Currently it includes code generators and scaffolding tools:

- [Go kit](https://github.com/go-kit/kit/) [endpoint](http://gokit.io/faq/#endpoints-mdash-what-are-go-kit-endpoints) generator (based on a service interface)
- Testify mock generator (similar to [mockery](https://github.com/vektra/mockery))
- Event dispatcher generator (based on event interface) (compatible with [Watermill](https://github.com/ThreeDotsLabs/watermill))
- Event handler generator (based on event structs) (compatible with [Watermill](https://github.com/ThreeDotsLabs/watermill))

**Roadmap:**

- Application scaffolding (based on [Modern Go Application](https://github.com/sagikazarmark/modern-go-application))
- [Go kit](https://github.com/go-kit/kit/) [service](http://gokit.io/faq/#services-mdash-what-is-a-go-kit-service) scaffolding
- [Go kit](https://github.com/go-kit/kit/) [transport](http://gokit.io/faq/#transports-mdash-what-are-go-kit-transports) scaffolding/code generator
- and more

See the [Modern Go Application](https://github.com/sagikazarmark/modern-go-application) for a detailed usage example.


## Installation

Download a prebuilt binary from the [Releases](https://github.com/sagikazarmark/mga/releases) page,
or install the tool from source:

```bash
go get sagikazarmark.dev/mga
```


## Usage

### Endpoint generator

An endpoint can be generated based on a service interface:

```go
package my

import (
    "context"
)

// +kit:endpoint

// Service is a business service.
type Service interface{
    // DoSomething is a service call.
    //
    // Named parameters and results are optional, but they make the generated code nicer.
    DoSomething(ctx context.Context, myparam string) (id string, err error)
}
```

Then run the generator:

```bash
mga generate kit endpoint ./...
```

See [Modern Go Application](https://github.com/sagikazarmark/modern-go-application/blob/master/internal/app/mga/todo/tododriver/zz_generated.endpoint.go) for an example.


### Testify mock generator

```go
package my

// +testify:mock

// Service is a business service.
type Service interface{}

// +testify:mock:testOnly=true

// Service2 is a business service.
type Service2 interface{}
```

```bash
mga generate testify mock ./...
```


### Event dispatcher generator

```go
package my

import (
    "context"
)

// +mga:event:dispatcher

type Events interface{
    MyEvent(ctx context.Context, ev MyEvent) error
}

type MyEvent struct{}
```

```bash
mga generate event dispatcher ./...
```

See [Modern Go Application](https://github.com/sagikazarmark/modern-go-application/blob/master/internal/app/mga/todo/todogen/zz_generated.event_dispatcher.go) for an example.


### Event handler generator

```go
package my

// +mga:event:handler

type Event struct{
    ID string

    MyParam string
}
```

```bash
mga generate event handler ./...
```

See [Modern Go Application](https://github.com/sagikazarmark/modern-go-application/blob/master/internal/app/mga/todo/todogen/zz_generated.event_handler.go) for an example.


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
