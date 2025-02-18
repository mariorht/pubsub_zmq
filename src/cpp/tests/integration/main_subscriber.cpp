#include "subscriber.hpp"
#include <fstream>
#include <sstream>
#include <iostream>

std::string format_json(const std::string& rawMessage) {
    // No estamos usando una librería, asumimos que el mensaje ya es un JSON válido
    return rawMessage;
}

int main() {
    Subscriber sub("tcp://integration_cpp_pub:5559");

    std::string lastMessage;

    for (int i = 0; i < 1; i++) {
        lastMessage = sub.receive();
    }

    // Guardar el último mensaje recibido en /shared/result.json
    std::ofstream file("/shared/result.json");
    if (file.is_open()) {
        file << format_json(lastMessage);
        file.close();
    } else {
        std::cerr << "❌ Error al escribir /shared/result.json" << std::endl;
    }

    return 0;
}
