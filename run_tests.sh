#!/bin/bash
set -e

# ============================
# CONFIGURACIÓN DE COLORES
# ============================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # Sin color

# ============================
# FUNCIONES AUXILIARES
# ============================
print_banner() {
  local color="$1"
  local text="$2"
  local border=$(printf '%*s' "${#text}" '' | tr ' ' '#')
  echo -e "\n${color}${border}${NC}"
  echo -e "${color}${text}${NC}"
  echo -e "${color}${border}${NC}\n"
}

compare_json() {
  pub_json=$(jq -S . ./shared/result_publisher.json)
  sub_json=$(jq -S . ./shared/result.json)

  if [ "$pub_json" == "$sub_json" ]; then
    echo -e "${GREEN}✅ Los mensajes coinciden entre Publisher y Subscriber${NC}"
  else
    echo -e "${RED}❌ Los mensajes NO coinciden entre Publisher y Subscriber${NC}"
    echo -e "${YELLOW}📄 Contenido de result_publisher.json:${NC}"
    echo "$pub_json"
    echo -e "${YELLOW}📄 Contenido de result.json:${NC}"
    echo "$sub_json"
    exit 1
  fi
}

clean_shared() {
  rm -f ./shared/result.json ./shared/result_publisher.json
}

# ============================
# TESTS UNITARIOS
# ============================
print_banner "${BLUE}" "INICIANDO TESTS UNITARIOS DE GO"
docker-compose run --rm go_unit_tests
print_banner "${GREEN}" "TESTS UNITARIOS DE GO COMPLETADOS"

print_banner "${BLUE}" "INICIANDO TESTS UNITARIOS DE PYTHON"
docker-compose run --rm python_unit_tests
print_banner "${GREEN}" "TESTS UNITARIOS DE PYTHON COMPLETADOS"

print_banner "${BLUE}" "INICIANDO TESTS UNITARIOS DE C++"
docker-compose build cpp_build
docker-compose run --rm cpp_build
docker-compose run --rm cpp_unit_tests
print_banner "${GREEN}" "TESTS UNITARIOS DE C++ COMPLETADOS"

# ============================
# TESTS DE INTEGRACIÓN
# ============================

# ----------------------------
# PYTHON -> GO
# ----------------------------
print_banner "${YELLOW}" "INICIANDO TESTS DE INTEGRACIÓN PYTHON -> GO"
clean_shared

PUBSUB_ENDPOINT=tcp://integration_python_pub:5555 docker-compose up integration_python_pub integration_go_sub

if [ ! -f ./shared/result.json ]; then
  echo -e "${RED}❌ No se generó el archivo de resultado en la integración PYTHON -> GO${NC}"
  exit 1
fi

compare_json
clean_shared
print_banner "${GREEN}" "TESTS DE INTEGRACIÓN PYTHON -> GO COMPLETADOS"

# ----------------------------
# GO -> PYTHON
# ----------------------------
print_banner "${YELLOW}" "INICIANDO TESTS DE INTEGRACIÓN GO -> PYTHON"
clean_shared

PUBSUB_ENDPOINT=tcp://integration_go_pub:5555 docker-compose up integration_python_sub integration_go_pub

if [ ! -f ./shared/result.json ]; then
  echo -e "${RED}❌ No se generó el archivo de resultado en la integración GO -> PYTHON${NC}"
  exit 1
fi

compare_json
clean_shared
print_banner "${GREEN}" "TESTS DE INTEGRACIÓN GO -> PYTHON COMPLETADOS"

# ----------------------------
# C++ -> C++
# ----------------------------
print_banner "${YELLOW}" "INICIANDO TESTS DE INTEGRACIÓN C++ -> C++"
clean_shared

PUBSUB_ENDPOINT=tcp://integration_cpp_pub:5555 docker-compose up integration_cpp_pub integration_cpp_sub

if [ ! -f ./shared/result.json ]; then
  echo -e "${RED}❌ No se generó el archivo de resultado en la integración C++ -> C++${NC}"
  exit 1
fi

compare_json
clean_shared
print_banner "${GREEN}" "TESTS DE INTEGRACIÓN C++ -> C++ COMPLETADOS"

# ============================
# FINALIZACIÓN
# ============================
print_banner "${GREEN}" "TODOS LOS TESTS FINALIZADOS CORRECTAMENTE 🎉"
