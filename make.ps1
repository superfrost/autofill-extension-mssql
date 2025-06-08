<#
.SYNOPSIS
Build script for autofillserver project
#>

# Configuration
$BINARY_NAME = "autofillserver"
$EXTENSION_DIR = "extension"
$DIST_DIR = "dist"
$ZIP_FILE = "$DIST_DIR/extension.zip"
$SRC_DIR = "server"

# List of extension files to package
$EXT_FILES = @("manifest.json", "popup.html", "popup.js", "icon16.png", "icon48.png", "icon128.png")

function Invoke-Tidy {
    Write-Host "go tidy..."
    Set-Location $SRC_DIR
    go mod tidy
    Set-Location ..
}

function Invoke-Clean {
    Write-Host "Cleaning..."
    if (Test-Path $DIST_DIR) {
        Remove-Item -Recurse -Force $DIST_DIR
    }
}

function Invoke-PackageExtension {
    Write-Host "Pack extension to ZIP..."
    
    if (!(Test-Path $DIST_DIR)) {
        New-Item -ItemType Directory -Path $DIST_DIR | Out-Null
    }
    
    # Copy extension files to dist directory
    Copy-Item -Path "$EXTENSION_DIR\*" -Destination $DIST_DIR -Recurse
    
    # Create zip archive
    Compress-Archive -Path "$DIST_DIR\*" -DestinationPath $ZIP_FILE -Force
    
    # Remove copied files
    foreach ($file in $EXT_FILES) {
        $filePath = Join-Path -Path $DIST_DIR -ChildPath $file
        if (Test-Path $filePath) {
            Remove-Item $filePath
        }
    }
}

function Invoke-BuildLinux {
    Write-Host "Building server for Linux..."
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Push-Location $SRC_DIR
    go build -o "../$DIST_DIR/${BINARY_NAME}_linux_amd64" main.go
    Pop-Location
}

function Invoke-BuildWindows {
    Write-Host "Building server for Windows..."
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Push-Location $SRC_DIR
    go build -o "../$DIST_DIR/${BINARY_NAME}_windows_amd64.exe" main.go
    Pop-Location
}

function Invoke-BuildServer {
    Invoke-Tidy
    Invoke-BuildWindows
    # Invoke-BuildLinux
}

function Invoke-All {
    Invoke-Clean
    Invoke-PackageExtension
    Invoke-BuildServer
}

# Main execution
switch ($args[0]) {
    "clean" { Invoke-Clean }
    "package_extension" { Invoke-PackageExtension }
    "build_linux" { Invoke-BuildLinux }
    "build_windows" { Invoke-BuildWindows }
    "build_server" { Invoke-BuildServer }
    default { Invoke-All }
}
