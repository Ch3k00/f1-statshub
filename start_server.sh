#!/bin/bash
# Script para iniciar el servidor asegurando CGO_ENABLED=1 para evitar errores

export CGO_ENABLED=1
go run server.go