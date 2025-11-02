#!/bin/bash
set -e

echo "--- Compilando proyecto Go ---"
ls -R
echo "------------------------------------"

# Construir y ejecutar el backend en Go
go build -o main .
