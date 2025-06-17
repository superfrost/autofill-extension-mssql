<#
.SYNOPSIS
Build script for autofillserver project
#>

# Configuration
$CRX_PACKAGE = "crx"
$NODE_VERSION = "v16.0.0"
$BINARY_NAME = "autofillserver"
$EXTENSION_DIR = "extension"
$OUTPUT_FILE_CRX = "autofiller_ext.crx"
$OUTPUT_FILE_XPI = "autofiller_ext.xpi"
$KEY_FILE = "./key.pem"
$DIST_DIR = "dist"
$ZIP_FILE = "$DIST_DIR/autofill_ext.zip"
$SRC_DIR = "server"

function Invoke-Check-Node {
    Write-Host "[?] Checking Node.js installation..."
    try {
        $nodeVersionOutput = node -v 2>&1
        if (-not $nodeVersionOutput -or $nodeVersionOutput -notmatch '^v\d+\.\d+\.\d+') {
            throw "Node.js is not installed"
        }
        Write-Host "[+] Node.js is installed: $nodeVersionOutput"
    } catch {
        Write-Host "[-] Node.js is not installed. Please install Node.js version $NODE_VERSION or higher."
        exit 1
    }
}

function Invoke-Check-Crx {
    Write-Host "[?] Checking installation of package $CRX_PACKAGE..."
    try {
        $crxCheck = npm list -g $CRX_PACKAGE --depth=0 2>&1
        if ($crxCheck -match "empty") {
            throw "Package is not installed"
        }
        Write-Host "[+] Package $CRX_PACKAGE is installed globally"
    } catch {
        Write-Host "[-] Package $CRX_PACKAGE is not installed globally. Install it using: npm install -g $CRX_PACKAGE"
        exit 1
    }
}

function Invoke-Check-Signkey {
    Write-Host "[?] Checking signature key..."
    if (-not (Test-Path $KEY_FILE)) {
        Write-Host "[!] Key not found. Generating a new key..."
                
        crx keygen ./
        Write-Host "[+] New key generated: $KEY_FILE"
    } else {
        Write-Host "[+] Key already exist: $KEY_FILE"
    }
}

function Invoke-PackageExtension-Crx-Xpi {
    Write-Host "Building the extension..."
    if (-not (Test-Path $DIST_DIR)) {
        New-Item -ItemType Directory -Path $DIST_DIR | Out-Null
    }

    crx pack $EXTENSION_DIR -p $KEY_FILE -o "$DIST_DIR/$OUTPUT_FILE_CRX"

    if (Test-Path "$DIST_DIR/$OUTPUT_FILE_CRX") {
        Write-Host "[+] Extension build OK: $DIST_DIR/$OUTPUT_FILE_CRX"
    } else {
        Write-Host "[-] Error building the extension"
        exit 1
    }

    Copy-Item -Path "$DIST_DIR/$OUTPUT_FILE_CRX" -Destination "$DIST_DIR/$OUTPUT_FILE_XPI"
    Write-Host "[+] Extension build OK: $DIST_DIR/$OUTPUT_FILE_XPI"
}

function Invoke-PackageExtension-Zip {
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

    Remove-Item "$DIST_DIR\popup"
    Remove-Item "$DIST_DIR\icons"

    Write-Host "[+] Pack extension to ZIP OK"
}

# List of extension files to package
$EXT_FILES = @("manifest.json", "style.css", "content.js", "background.js", "popup/popup.html", "popup/popup.js", "icons/icon16.png", "icons/icon48.png", "icons/icon128.png")

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
    Write-Host "[+] Cleaning OK"
}

function Invoke-BuildLinux {
    Write-Host "Building server for Linux..."
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Push-Location $SRC_DIR
    go build -o "../$DIST_DIR/${BINARY_NAME}_linux_amd64" main.go
    Pop-Location
    Write-Host "[+] Building server for Linux OK"
}

function Invoke-BuildWindows {
    Write-Host "Building server for Windows..."
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Push-Location $SRC_DIR
    go build -o "../$DIST_DIR/${BINARY_NAME}_windows_amd64.exe" main.go
    Pop-Location
    Write-Host "[+] Building server for Windows OK"
}

function Invoke-BuildServer {
    Invoke-Tidy
    Invoke-BuildWindows
    # Invoke-BuildLinux
}

function Invoke-PackageExtension {
    Invoke-Check-Node
    Invoke-Check-Crx
    Invoke-Check-Signkey
    Invoke-PackageExtension-Crx-Xpi
    Invoke-PackageExtension-Zip
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
