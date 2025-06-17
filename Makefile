# Имя выходных файлов
CRX_PACKAGE := crx
NODE_VERSION := v16.0.0
BINARY_NAME = autofillserver
EXTENSION_DIR = extension
OUTPUT_FILE_CRX := autofiller_ext.crx
OUTPUT_FILE_XPI := autofiller_ext.xpi
KEY_FILE := ./key.pem
DIST_DIR = dist
ZIP_FILE = $(DIST_DIR)/autofill_ext.zip

# Папка с исходным кодом сервера
SRC_DIR = server

# Список файлов для упаковки расширения
EXT_FILES = manifest.json style.css content.js background.js popup/popup.html popup/popup.js icons/icon16.png icons/icon48.png icons/icon128.png popup/ icons/

# Цели для сборки
all: clean tidy build_extensions build_server

tidy:
	@echo "go tidy..."
	@cd $(SRC_DIR) && go mod tidy

dev: tidy
	@cd $(SRC_DIR) && go run main.go

# Упаковка расширения в ZIP
package_extension_zip:
	@mkdir -p $(DIST_DIR)
	@echo "📦 Упаковка файлов в ZIP..."
	@cp -r $(EXTENSION_DIR)/* $(DIST_DIR)
	@cd $(DIST_DIR) && zip -r ../$(ZIP_FILE) .
	@cd $(DIST_DIR) && rm -fd $(EXT_FILES)
	@echo "✅ Файлы расширения упакованы в zip"

build_crx_xpi: check-node check-crx generate-key
	@echo "🔧 Сборка расширений crx, xpi..."
	@mkdir -p $(DIST_DIR)
	crx pack $(EXTENSION_DIR) -p $(KEY_FILE) -o $(DIST_DIR)/$(OUTPUT_FILE_CRX)
	@echo "✅ Расширение успешно собрано: $(DIST_DIR)/$(OUTPUT_FILE_CRX)"
	@cp $(DIST_DIR)/$(OUTPUT_FILE_CRX) $(DIST_DIR)/$(OUTPUT_FILE_XPI)
	@echo "✅ Расширение успешно собрано: $(DIST_DIR)/$(OUTPUT_FILE_XPI)"

check-node:
	@echo "🔍 Проверка установки Node.js..."
	@which node > /dev/null || (echo "❌ Node.js не установлен. Установите Node.js версии $(NODE_VERSION) или выше."; exit 1)
	@node -v | grep -q '^v' || (echo "❌ Node.js не установлен корректно"; exit 1)
	@echo "✅ Node.js установлен: $$(node -v)"

check-crx:
	@echo "🔍 Проверка установки пакета crx..."
	@npm list -g $(CRX_PACKAGE) --depth=0 > /dev/null 2>&1 || (echo "❌ Пакет crx не установлен глобально. Установите: npm install -g crx"; exit 1)
	@echo "✅ Пакет crx установлен"

generate-key:
	@echo "🔑 Проверка ключа подписи..."
	@if [ ! -f $(KEY_FILE) ]; then \
		echo "⚠️ Ключ не найден. Генерация нового ключа..."; \
		openssl genrsa -out $(KEY_FILE) 2048; \
		echo "✅ Новый ключ сгенерирован: $(KEY_FILE)"; \
	else \
		echo "✅ Ключ уже существует: $(KEY_FILE)"; \
	fi

build_extensions: package_extension_zip build_crx_xpi

# Сборка для Linux
build_linux:
	@echo "Сборка сервера для Linux..."
	@cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_linux_amd64 main.go
	@echo "✅ Сборка сервера для Linux завершена"

# Сборка для Windows
build_windows:
	@echo "Сборка сервера для Windows..."
	@cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_windows_amd64.exe main.go
	@echo "✅ Сборка сервера для Windows завершена"

# Сборка для macOS intel
build_darwin:
	@echo "Сборка сервера для MacOS..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_darwin_amd64 main.go
	@echo "✅ Сборка сервера для MacOS завершена"

# Сборка для macOS arm
build_darwin_arm64:
	@echo "Сборка сервера для MacOS arm64..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../$(DIST_DIR)/$(BINARY_NAME)_darwin_arm64 main.go
	@echo "✅ Сборка сервера для MacOS arm64 завершена"

# Сборка сервера для всех платформ
build_server: build_linux build_windows build_darwin build_darwin_arm64

# Очистка выходных файлов
clean:
	@echo "🧹 Очистка..."
	rm -rf $(DIST_DIR)
	@echo "✓ Директория dist очищена"

# Цель по умолчанию
.PHONY: all tidy clean build_extensions package_extension_zip build_crx_xpi check-node check-crx generate-key build_server build_linux build_windows build_darwin build_darwin_arm64
