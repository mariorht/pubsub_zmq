#include "publisher.hpp"
#include <thread>
#include <chrono>
#include <fstream>
#include <iostream>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

int main() {
    Publisher pub("tcp://*:5559");

    zmq_sleep(2);

    json lastMessage;

    for (int i = 0; i < 1; ++i) {
        std::map<std::string, std::string> data = {
            {"key", "probando"},
            {"index", std::to_string(i)}
        };

        auto message_chunks = pub.build_message(data);
        pub.publish_message(message_chunks);

        // Guardar el mensaje completo como lo construye build_message
        lastMessage["type"] = "images";
        lastMessage["count"] = 0;
        lastMessage["images"] = json::array();
        lastMessage["data"] = data;

        std::ofstream file("/shared/result_publisher.json");
        if (file.is_open()) {
            file << lastMessage.dump(4); // Bonito
            file.close();
            std::cout << "✅ Mensaje enviado guardado en /shared/result_publisher.json" << std::endl;
        } else {
            std::cerr << "❌ Error al escribir /shared/result_publisher.json" << std::endl;
        }

        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    return 0;
}
