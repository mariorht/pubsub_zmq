#include "publisher.hpp"
#include <thread>
#include <chrono>
#include <fstream>
#include <sstream>

using namespace std::chrono_literals;

std::string build_json(const std::string& key, int index) {
    std::stringstream json;
    json << "{"
         << "\"key\":\"" << key << "\","
         << "\"index\":" << index
         << "}";
    return json.str();
}

int main() {
    Publisher pub("tcp://*:5559");

    using namespace std::chrono_literals;  // Habilita el sufijo "ms"
    std::this_thread::sleep_for(2000ms);

    std::string lastMessage;

    for (int i = 0; i < 1; i++) {
        lastMessage = build_json("probando", i);
        pub.publish(lastMessage);

        // Escribir el último mensaje enviado en /shared/result_publisher.json
        std::ofstream file("/shared/result_publisher.json");
        if (file.is_open()) {
            file << lastMessage;
            file.close();
        }

        std::this_thread::sleep_for(1s);
    }

    return 0;
}
