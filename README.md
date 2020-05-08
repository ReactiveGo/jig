# generics

    $ go get github.com/reactivego/generics/cmd/jig

[![](http://godoc.org/github.com/reactivego/generics?status.png)](http://godoc.org/github.com/reactivego/generics)


Package `generics` provides a generator to add support for generic programming to Go 1.

What is called generics here, is more accurately described as parametric polymorphism. In essence, generics are fragments of code with type place-holders. The code fragments are transformed into statically typed code by replacing these type place-holders with concrete types.

The generator tool (named [jig](cmd/jig)) can expand generic functions and datatypes from a library into statically typed code.
To test whether it was installed correctly, go to the command-line and run `jig -h`:

```bash
$ jig -h
Usage of jig [flags] [<dir>]:
  -c, --clean     Remove files generated by jig
  -r, --regen     Force regeneration of all code by jig (default)
  -m, --missing   Only generate code that is missing
  -v, --verbose   Print details of what jig is doing
```
## Stack example

Let's implement a generic stack.

We start with two templates that use `foo` as a type place-holder:

```go
package stack

type foo int

//jig:template <Foo>Stack

type FooStack []foo
var zeroFoo foo

//jig:template <Foo>Stack Push

func (s *FooStack) Push(v foo) {
  *s = append(*s, v)
}
```
Note that this library contains **valid** Go code that compiles normally.

This code is also available [online](example/stack/generic) and can be imported as follows:

```go
import _ "github.com/reactivego/generics/example/stack/generic"
```

Given some code that uses this Library in a file `main.go`.

```go
package main

import _ "github.com/reactivego/generics/example/stack/generic"

func main() {
  var stack StringStack
  stack.Push("Hello, World!")
}
```
Running `jig` will expand the two templates and write them to the file `stack.go`.

```go
// Code generated by jig; DO NOT EDIT.

//go:generate jig

package main

//jig:name StringStack

type StringStack []string
var zeroString string

//jig:name StringStackPush

func (s *StringStack) Push(v string) {
  *s = append(*s, v)
}
```
The generated code is a type-safe stack generated by expanding templates for the type `string`

## rx example

In this example you will learn how to write a minimal program that uses [Reactive Extensions for Go](https://github.com/reactivego/rx).

### Installation


To install `rx`, open a terminal and from the command-line run:

```bash
$ go get github.com/reactivego/rx
```
### Overview

- You write code that references templates from the `rx` library.
- You run the `jig` command in the directory where your code is located.
- **Now `jig` analyzes your code and determines what additional code is needed to make it build**.
- *Jig* takes templates from the `rx` library and specializes them on specific types.
- Specializations are generated into the file `rx.go` alongside your own code.
- If all went well, your code will now build.

### Prepare Folder

Let's create a new folder for our simple program and start editing the file `main.go`.

```bash
$ cd $(go env GOPATH)
$ mkdir -p ./src/helloworld
$ cd ./src/helloworld
$ subl main.go
```
> NOTE we use `subl` to open *Sublime Text*, but any text editor will do.

### Write Code
Now that you have your [`main.go`](../example/rx/main.go) file open in your editor of choice, type the following code:

![Hello World Program](doc/helloworld.png)

If you would now go to the command-line and run `jig` it would do the following:

1. Import `github.com/reactivego/rx/generic` to access the generics in that library.
2. Generate `FromString` by specializing the template `From<Foo>` from the library on type `string`. Generate dependencies of `FromString` like the type `ObservableString` that is returned by `FromString`.
3. Generate the `MapString` method of `ObservableString` by specializing the template `Observable<Foo> Map<Bar>` for `Foo` as type `string` and `Bar` (the type to map to) also as type `string`.
4. Map function from `string` to `string` just concatenates two strings and returns the result.
5. Print every string returned by the Map function.
6. The output you can expect when you run the program.

Now actually go to the command-line and run `jig -v`. Use the verbose flag `-v` because otherwise `jig` will be silent.

```bash
$ jig -v
found 126 templates in package "rx" (github.com/reactivego/rx/generic)
found 16 templates in package "multicast" (github.com/reactivego/multicast/generic)
generating "FromString"
  Scheduler
  Subscriber
  StringObserveFunc
  zeroString
  ObservableString
  FromSliceString
  FromString
generating "ObservableString MapString"
  ObservableString MapString
generating "ObservableString Println"
  Schedulers
  ObservableString Println
writing file "rx.go"
```

Now we can try to run the code and see what it does.

```bash
$ go run *.go
Hello, You!
Hello, Gophers!
Hello, World!
```

Success! `jig` generated the code into the file `rx.go` and we were able to run the program.
Turns out the generated file [`rx.go`](example/rx/rx.go) contains less than 250 lines of code.

If you add additional code to the program that uses different generics of the `rx` library, then you should run `jig` again to generate specializations of those generics.

