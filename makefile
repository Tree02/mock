# nombre de ejecutable
BINARY_NAME=mocklogin
MAIN_FILE=./cmd/main.go

GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Ejecutable generado: $(BINARY_NAME)"

# Linux
#.PHONY: build-linux
#build-linux:
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) $(MAIN_FILE)
#	@echo "Ejecutable generado para Linux: $(BINARY_NAME)"

# Windows
#.PHONY: build-windows
#build-windows:
#	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(MAIN_FILE)
#	@echo "Ejecutable generado para Windows: $(BINARY_NAME).exe"

# macOS
#.PHONY: build-macos
#build-macos:
#	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME) $(MAIN_FILE)
#	@echo "Ejecutable generado para macOS: $(BINARY_NAME)"

# binarios generados
.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@echo "Limpiados los ejecutables generados."

# Comando para ayuda en caso de que el usuario necesite
.PHONY: help
help:
	@echo "Comandos disponibles:"
	@echo "  make build           - Construye un binario estático para el sistema actual."
	@echo "  make build-linux     - Construye un binario estático para Linux."
	@echo "  make build-windows   - Construye un binario estático para Windows."
	@echo "  make build-macos     - Construye un binario estático para macOS."
	@echo "  make clean           - Elimina los binarios generados."
