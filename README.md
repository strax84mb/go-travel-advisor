# Table of contents
1. [VSCode setup](#vscode-setup)
2. [Controllers](#controllers)
    1. [DTOs using FFJSON](#dtos-using-ffjson)
    2. [Using generics to extract and validate query and path variables](#using-generics-to-extract-and-validate-query-and-path-variables)
    3. [JWT](#jwt)
3. [Services](#services)
    1. [Logging](#logging)
    2. [Destructors](#destructors)
4. [Test files](#test-files)
    1. [Test files in package \<name\>_test](#test-files-in-package-name_test)
    2. [Test files in same package as developed code](#test-files-in-same-package-as-developed-code)
    3. [Mocking interfaces](#mocking-interfaces)

## VSCode setup
In directory __.vscode__ create file __settings.json__ and save this as its contents:
```JSON
{
    "go.buildTags": "test"
}
```
This will be needed in order to be able to compile test files.

## Controllers

Router used is `Gorilla`.

### DTOs using FFJSON

DTO serialization/deserialization is performed using FFJSON. That way, reflection is avoided as much as possible resulting in 2x~3x increase in performance.

In order to use FFJSON it must be installed onto the system (not just project):
```Bash
go get -u github.com/pquerna/ffjson
ffjson [-nodecoder] [-noencoder] myfile.go
```

This will create file `myfile_ffjson.go`. That file will contain logic for marshalling and unmarshalling without the need for reflection.

### Using generics to extract and validate query and path variables

As Go supports generics now (although not fully using monomorphism), getting parameters from query or path doesn't need to be done with a separate function for each return type. We can create generic functions and types. Cool thing about Go generics is that stating every single generic parameter is not needed. We can state only one as long as the rest can be derived from it. Even that one parameter is not needed if it can be implied by one of the regular parameters.

In this project, a separate parser is defined for every expected type. That parser is one of the parameters in `Path` or `Query` generic function.
```Go
type Value[V int | int64 | string] struct {
	val V
}

func (v *Value[V]) Val() V {
	return v.val
}

type PrimitiveValue[
	V int | int64 | string,
	T Value[V],
] interface {
	Val() V
	*T
}

func Query[
	V int | int64 | string,
	T Value[V],
	PV PrimitiveValue[V, T],
](
	r *http.Request,
	parse func(string) (V, error),
	param string,
	mandatory bool,
	defaultValue V,
	validators ...AtomicValidator[V, T, PV],
) (T, error) {
    ***
}
```
So, because parser for int64 is used, it is implied that entire flow should be used to return parsed int64. It may seem complicated, but everything is connected. Only first generic parameter needs to be defined/implied. All others can be derived from it. When needing to actually use this, it can be done like so:
```Go
beginID, err := handler.Query(r, handler.Int64, "begin", true, 0, handler.IsPositive)
```
`Int64` is important in this as it is defined like this:
```Go
func Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}
```
and second parameter in Query is defined like:
```Go
parse func(string) (V, error)
```
Because of this, `V` is defined as `int64` and everything else is easily derived from this. So, no need to actually state `int64` generic parameter when calling `Query`.

### JWT

JWT claims used are:

|Claim|Explanation|
|-----|-----------|
|**sub**|User name here but usually it is some ID.|
|**username**|User name again as it is good practice to not rely on sub claim for username or user ID.|
|**user-id**|ID of user|
|**roles**|Comma separated list of roles.|
|**nbf**|Epoch time after which JWT is considered valid.|
|**iat**|Epoch time at which JWT is generated.|
|**exp**|Epoch time at which JWT expires. After that time JWT is invalid.|

## Services

### Logging

Logging is done by `logrus` in JSON format.

### Destructors

It is recommended to write destructors for large and complicated structures. When executing a destructor, all arrays, maps and pointers (immediate or nested) should be assigned a value `nil`. This will help out garbage collector as scanning and collection will be faster.

## Test files

### Test files in package \<name\>_test
This allows for the package code to be tested as a black box. Only exposed code can be accessed. This approach is a problem if interfaces used are not exposed.

### Test files in same package as developed code
This allows for testing of code that is not exposed. In order to avoid having test code ending up in production, conditional build tag must be used at the start of the file.
```Go
//go:build test
```

### Mocking interfaces

When using interfaces as fields of structures we need some way to mock these interfaces. It is trivial to write such mocks but why not automate the task?

We can't use testify for mocking as it can only mock structures. It is fine to use testify for assertions, though. For our purposes we can use mockgen.
```Bash
go install github.com/golang/mock/mockgen@v1.6.0
```
Now, if interfaces being used are not exposed mockgen needs to be used in source mode like so:
```Bash
mockgen -source=foo.go
```
Mock structs will be generated and those can be used as the implementation of interfaces in `foo.go`. These mocks can be used like this:
```Go
func TestFoo(t *testing.T) {
  ctrl := gomock.NewController(t)
  defer ctrl.Finish()

  m := NewMockFoo(ctrl)

  // Does not make any assertions. Executes the anonymous functions and returns
  // its result when Bar is invoked with 99.
  m.
    EXPECT().
    Bar(gomock.Eq(99)).
    DoAndReturn(func(_ int) int {
      time.Sleep(1*time.Second)
      return 101
    }).
    AnyTimes()

  // Does not make any assertions. Returns 103 when Bar is invoked with 101.
  m.
    EXPECT().
    Bar(gomock.Eq(101)).
    Return(103).
    AnyTimes()

  SUT(m)
}
```
