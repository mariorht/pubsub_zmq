#!/bin/bash
set -e

# Definir colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # Sin color

print_banner() {
  local color="$1"
  local text="$2"
  local border=$(printf '%*s' "${#text}" '' | tr ' ' '#')
  echo -e "${color}${border}${NC}"
  echo -e "${color}${text}${NC}"
  echo -e "${color}${border}${NC}"
}

echo -e "\n"
print_banner "${BLUE}" "INICIANDO TESTS UNITARIOS DE GO"
echo -e "\n"

docker-compose run --rm go

echo -e "\n"
print_banner "${GREEN}" "TESTS UNITARIOS DE GO COMPLETADOS"
echo -e "\n"

print_banner "${YELLOW}" "INICIANDO TESTS UNITARIOS DE PYTHON"
echo -e "\n"

docker-compose run --rm python

echo -e "\n"
print_banner "${GREEN}" "TESTS UNITARIOS DE PYTHON COMPLETADOS"
echo -e "\n"

