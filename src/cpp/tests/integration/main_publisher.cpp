#include "publisher.hpp"
#include <opencv2/opencv.hpp> // Necesario para cargar imágenes
#include <thread>
#include <chrono>
#include <fstream>
#include <iostream>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

int main() {
    Publisher pub("tcp://*:5555", "test");

    zmq_sleep(5);

    json lastMessage;
    std::vector<json> images_metadata;

    for (int i = 0; i < 1; ++i) {
        std::map<std::string, std::string> data = {
            {"key", "probando"},
            {"index", std::to_string(i)}
        };

        std::vector<cv::Mat> frames;

        // 🔹 Intentar cargar una imagen real o generar una en memoria
        cv::Mat img = cv::imread("/assets/pong.png"); // Ruta de prueba
        if (img.empty()) {
            std::cerr << "⚠️ No se encontró /assets/test_image.png, generando una imagen artificial." << std::endl;
            img = cv::Mat(100, 100, CV_8UC3, cv::Scalar(255, 0, 0)); // Imagen azul 100x100
        }
        frames.push_back(img); // 🔹 Añadir imagen al mensaje

        // 🔹 Publicar mensaje con imagen
        auto message_chunks = pub.build_message(frames, data, "png"); // 🔹 Ahora enviamos una imagen PNG
        pub.publish_message(message_chunks);


        // 🔹 Generar metadatos de la imagen
        json imageMeta;
        imageMeta["metadata"]["width"] = img.cols;
        imageMeta["metadata"]["height"] = img.rows;
        imageMeta["metadata"]["channels"] = img.channels();
        imageMeta["metadata"]["dtype"] = "uint8"; // OpenCV usa uint8
        imageMeta["metadata"]["size"] = img.total() * img.elemSize();
        images_metadata.push_back(imageMeta);

        // Guardar el mensaje completo como lo construye build_message
        lastMessage["type"] = "images";
        lastMessage["count"] = frames.size();
        lastMessage["images"] = images_metadata; // 🔹 Ahora se guardan los metadatos
        lastMessage["data"] = data;

        std::ofstream file("/shared/result_publisher.json");
        if (file.is_open()) {
            file << lastMessage.dump(4); // Formateado bonito
            file.close();
            std::cout << "✅ Mensaje con imagen enviado guardado en /shared/result_publisher.json" << std::endl;
        } else {
            std::cerr << "❌ Error al escribir /shared/result_publisher.json" << std::endl;
        }

        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    return 0;
}
