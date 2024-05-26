param (
    [string]$Task,
    [string]$MAIN_PACKAGE_PATH = "./cmd/server",
    [string]$MAIN_SERVER_FILE = "api.go",
    [string]$BINARY_NAME = "server",
    [string]$API_DOCS = "./docs"
)

function Build {
    # Crear el directorio bin si no existe
    $binPath = "bin"
    if (-Not (Test-Path -Path $binPath)) {
        New-Item -ItemType Directory -Path $binPath
    }

    # Construir el binario
    $binaryPath = Join-Path -Path $binPath -ChildPath $BINARY_NAME
    $mainFilePath = Join-Path -Path $MAIN_PACKAGE_PATH -ChildPath $MAIN_SERVER_FILE

    if (Test-Path -Path $mainFilePath) {
        go build -o $binaryPath $mainFilePath
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Build failed."
            exit $LASTEXITCODE
        }
    } else {
        Write-Error "Main file not found: $mainFilePath"
        exit 1
    }
}

function Run {
    # Ejecutar el binario
    $binaryPath = Join-Path -Path "bin" -ChildPath $BINARY_NAME
    if (Test-Path -Path $binaryPath) {
        Start-Process -NoNewWindow -FilePath $binaryPath
    } else {
        Write-Error "Binary file not found: $binaryPath"
        exit 1
    }
}

function TestUnit {
    # Pruebas unitarias
    go test -v ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Unit tests failed."
        exit $LASTEXITCODE
    }
}

function TestUnitCover {
    # Pruebas unitarias con cobertura
    $tmpPath = "tmp"
    if (-Not (Test-Path -Path $tmpPath)) {
        New-Item -ItemType Directory -Path $tmpPath
    }
    $coveragePath = Join-Path -Path $tmpPath -ChildPath "coverage.out"

    go test -v -coverprofile=$coveragePath ./...
    if ($LASTEXITCODE -eq 0) {
        go tool cover -func=$coveragePath
    } else {
        Write-Error "Unit tests with coverage failed."
        exit $LASTEXITCODE
    }
}

function TestAll {
    # Todas las pruebas
    $env:RUN_INTEGRATION_TESTS = 1
    go test -v ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "All tests failed."
        exit $LASTEXITCODE
    }
}

function TestAllCover {
    # Todas las pruebas con cobertura
    $tmpPath = "tmp"
    if (-Not (Test-Path -Path $tmpPath)) {
        New-Item -ItemType Directory -Path $tmpPath
    }
    $env:RUN_INTEGRATION_TESTS = 1
    $coveragePath = Join-Path -Path $tmpPath -ChildPath "coverage.out"

    go test -v -coverprofile=$coveragePath ./...
    if ($LASTEXITCODE -eq 0) {
        go tool cover -func=$coveragePath
    } else {
        Write-Error "All tests with coverage failed."
        exit $LASTEXITCODE
    }
}

function ApiSpec {
    # Generar documentación de la API
    $mainFilePath = Join-Path -Path $MAIN_PACKAGE_PATH -ChildPath $MAIN_SERVER_FILE
    if (Test-Path -Path $mainFilePath) {
        swag init -g $mainFilePath -o $API_DOCS
        if ($LASTEXITCODE -ne 0) {
            Write-Error "API specification generation failed."
            exit $LASTEXITCODE
        }
    } else {
        Write-Error "Main file not found: $mainFilePath"
        exit 1
    }
}

switch ($Task) {
    "build" { Build }
    "run" { Run }
    "test-unit" { TestUnit }
    "test-unit-cover" { TestUnitCover }
    "test-all" { TestAll }
    "test-all-cover" { TestAllCover }
    "api-spec" { ApiSpec }
    default { Write-Error "Invalid task specified. Valid tasks are: build, run, test-unit, test-unit-cover, test-all, test-all-cover, api-spec." }
}
