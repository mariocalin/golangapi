# Library web api

Pet project to translate java springboot knowledge into golang

## Generate mocks

````powershell
PS path\to\folder> mockery --name=MyInterface --outpkg=package --dir=. --output=. --filename=whatever_mock_test.go --inpackage
```

It may be needed to change the name of the type that it is generated