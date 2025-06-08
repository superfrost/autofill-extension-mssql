# Имя выходных файлов
$BINARY_NAME = "autofillserver"
$EXTENSION_DIR = "extension"
$DIST_DIR = "dist"
$ZIP_FILE = Join-Path $DIST_DIR "extension.zip"

# Папка с исходным кодом сервера
$SRC_DIR = "server"

# Список файлов для упаковки расширения
$EXT_FILES = @("manifest.json", "popup.html", "popup.js", "icon16.png", "icon48.png", "icon128.png")

# Основная функция
function Main {
    Clean
    Package-Extension
    Build-Server
}

# Упаковка расширения в ZIP
function Package-Extension {
    if (-Not (Test-Path $DIST_DIR)) {
        New-Item -ItemType Directory -Path $DIST_DIR | Out-Null
    }
    Write-Host "Упаковка файлов в ZIP..."
    Copy-Item -Path (Join-Path $EXTENSION_DIR "*") -Destination $DIST_DIR -Recurse
    Compress-Archive -Path (Join-Path $DIST_DIR "*") -DestinationPath $ZIP_FILE
    Remove-Item -Path (Join-Path $DIST_DIR $EXT_FILES) -Force
}

# Сборка для Linux
function Build-Linux {
    Write-Host "Сборка сервера для Linux..."
    Set-Location $SRC_DIR
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go build -o (Join-Path ".." (Join-Path $DIST_DIR "$BINARY_NAME`_linux_amd64")) main.go
    Set-Location ..
}

# Сборка для Windows
function Build-Windows {
    Write-Host "Сборка сервера для Windows..."
    Set-Location $SRC_DIR
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o (Join-Path ".." (Join-Path $DIST_DIR "$BINARY_NAME`_windows_amd64.exe")) main.go
    Set-Location ..
}

# Сборка для macOS intel
function Build-Darwin {
    Write-Host "Сборка сервера для macOS..."
    Set-Location $SRC_DIR
    $env:GOOS = "darwin"
    $env:GOARCH = "amd64"
    go build -o (Join-Path ".." (Join-Path $DIST_DIR "$BINARY_NAME`_darwin_amd64")) main.go
    Set-Location ..
}

# Сборка для macOS arm
function Build-Darwin-Arm64 {
    Write-Host "Сборка сервера для macOS arm64..."
    Set-Location $SRC_DIR
    $env:GOOS = "darwin"
    $env:GOARCH = "arm64"
    go build -o (Join-Path ".." (Join-Path $DIST_DIR "$BINARY_NAME`_darwin_arm64")) main.go
    Set-Location ..
}

# Сборка сервера для всех платформ
function Build-Server {
    Build-Linux
    Build-Windows
    Build-Darwin
    Build-Darwin-Arm64
}

# Очистка выходных файлов
function Clean {
    Write-Host "Очистка..."
    Remove-Item -Path $DIST_DIR -Recurse -Force -ErrorAction SilentlyContinue
}

# Запуск основной функции
Main
