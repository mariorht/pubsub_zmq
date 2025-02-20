#include "subscriber.hpp"
#include <fstream>
#include <iostream>
#include <sstream>

int main() {
    Subscriber sub("tcp://integration_cpp_pub:5559");

    std::string lastMessage;

    for (int i = 0; i < 1; i++) {
        lastMessage = sub.receive();
    }

    // Guardar el último mensaje recibido en /shared/result.json
    std::ofstream file("/shared/result.json");
    if (file.is_open()) {
        file << lastMessage;
        file.close();
        std::cout << "✅ Mensaje recibido guardado en /shared/result.json" << std::endl;
    } else {
        std::cerr << "❌ Error al escribir /shared/result.json" << std::endl;
    }

    return 0;
}
