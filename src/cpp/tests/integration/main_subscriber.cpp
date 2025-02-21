#include <iostream>
#include <fstream>
#include <nlohmann/json.hpp>
#include <opencv2/opencv.hpp>
#include <cstdlib> // Para getenv

#include "subscriber.hpp"
#include "publisher.hpp"

using ValueVariant = std::variant<std::string, int, double, bool, std::nullptr_t>;
using DataMap = std::unordered_map<std::string, ValueVariant>;

// Función auxiliar para convertir DataMap a JSON
nlohmann::json variant_map_to_json(const DataMap& data) {
    nlohmann::json json_data;

    for (const auto& [key, value] : data) {
        std::visit([&](auto&& arg) {
            json_data[key] = arg;
        }, value);
    }

    return json_data;
}

int main() {
    const char* env_endpoint = std::getenv("PUBSUB_ENDPOINT");
    std::string endpoint = env_endpoint ? env_endpoint : "tcp://localhost:5555";
    Subscriber sub(endpoint, "test");

    nlohmann::json lastMessage;

    for (int i = 0; i < 1; ++i) {
        auto [images, data] = sub.receive_message();

        for (const auto& [k, v] : data) {
            std::cout << k << ": ";
            std::visit([](auto&& arg) { std::cout << arg; }, v);
            std::cout << std::endl;
        }

        lastMessage["type"] = "images";
        lastMessage["count"] = images.size();
        lastMessage["images"] = nlohmann::json::array();

        // Aquí está el cambio clave:
        lastMessage["data"] = variant_map_to_json(data);

        for (const auto& img : images) {
            nlohmann::json img_meta;
            img_meta["width"] = img.cols;
            img_meta["height"] = img.rows;
            img_meta["channels"] = img.channels();
            img_meta["size"] = img.total() * img.elemSize();
            img_meta["dtype"] = Publisher::mat_type_to_dtype_string(img.depth());
            lastMessage["images"].push_back({{"metadata", img_meta}});
        }
    }

    // Guarda el último mensaje recibido en JSON
    std::ofstream file("/shared/result.json");
    if (file.is_open()) {
        file << lastMessage.dump(4);
        std::cout << "✅ Mensaje recibido guardado en /shared/result.json" << std::endl;
    } else {
        std::cerr << "❌ Error al escribir /shared/result.json" << std::endl;
    }

    return 0;
}
