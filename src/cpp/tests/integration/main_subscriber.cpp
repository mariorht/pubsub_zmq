#include "subscriber.hpp"
#include <iostream>
#include <fstream>
#include <nlohmann/json.hpp>
#include <opencv2/opencv.hpp>

int main() {
    Subscriber sub("tcp://integration_cpp_pub:5559");

    nlohmann::json lastMessage;

    for (int i = 0; i < 1; ++i) {
        auto [images, data] = sub.receive_message();

        std::cout << "📥 Recibido -> key: " << data["key"] << ", index: " << data["index"] << std::endl;

        lastMessage["type"] = "images";
        lastMessage["count"] = images.size();
        lastMessage["images"] = nlohmann::json::array();
        lastMessage["data"] = data;

        for (const auto& img : images) {
            nlohmann::json img_meta;
            img_meta["width"] = img.cols;
            img_meta["height"] = img.rows;
            img_meta["channels"] = img.channels();
            img_meta["size"] = img.total() * img.elemSize();

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
