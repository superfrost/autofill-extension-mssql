# Имя выходных файлов
BINARY_NAME = autofillserver
EXTENSION_DIR = extension
DIST_DIR = dist
ZIP_FILE = $(DIST_DIR)/extension.zip

# Папка с исходным кодом сервера
SRC_DIR = server

# Список файлов для упаковки расширения
EXT_FILES = manifest.json popup.html popup.js icon16.png icon48.png icon128.png

# Цели для сборки
all: clean tidy package_extension build_server

tidy:
	@echo "go tidy..."
	@cd $(SRC_DIR) && go mod tidy

# Упаковка расширения в ZIP
package_extension:
	@mkdir -p $(DIST_DIR)
	@echo "Упаковка файлов в ZIP..."
	@cp -r $(EXTENSION_DIR)/* $(DIST_DIR)
	@cd $(DIST_DIR) && zip -r ../$(ZIP_FILE) .
	@cd $(DIST_DIR) && rm -f $(EXT_FILES)

# Сборка для Linux
build_linux:
	@echo "Сборка сервера для Linux..."
	@cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_linux_amd64 main.go

# Сборка для Windows
build_windows:
	@echo "Сборка сервера для Windows..."
	@cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_windows_amd64.exe main.go

# Сборка для macOS intel
build_darwin:
	@echo "Сборка сервера для MacOS..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_darwin_amd64 main.go

# Сборка для macOS arm
build_darwin_arm64:
	@echo "Сборка сервера для MacOS arm64..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_darwin_arm64 main.go

# Сборка сервера для всех платформ
build_server: build_linux build_windows build_darwin build_darwin_arm64

# Очистка выходных файлов
clean:
	@echo "Очистка..."
	rm -rf $(DIST_DIR)

# Цель по умолчанию
.PHONY: all tidy clean package_extension build_server build_linux build_windows build_darwin build_darwin_arm64
