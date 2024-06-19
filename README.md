## VSCode setup
In directory __.vscode__ create file __settings.json__ and save this as its contents:
```JSON
{
    "go.buildTags": "test"
}
```
This will be needed in order to be able to compile test files.

## Test files

### Test files in package \<name\>_test
This allows for the package code to be tested as a black box. Only exposed code can be accessed. This approach is a problem if interfaces used are not exposed.

### Test files in same package as developed code
This allows for testing of code that is not exposed. In order to avoid having test code ending up in production, conditional build tag must be used at the start of the file. 
```Go
//go:build test
```