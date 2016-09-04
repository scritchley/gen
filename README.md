# gen

Gen is a tool for implementing generic code in Go.

### Getting started

```
go install github.com/scritchley/gen
```

### Examples

See the `examples` directory. Full explanation can be found in my blog post: [Generic code generation in Go](http://simon-critchley.co.uk/generic-code-generation-in-go/).

### TODO

- [ ] If a stub type of `T__` has properties defined, gen should check that the target type includes those properties otherwise throw an error.
- [ ] The tool should be able to target specific generic types within a generic package. Something like @importpath:T__Foo
- [ ] Support for transformations, i.e. a method accepting T__ and returning U__. 
