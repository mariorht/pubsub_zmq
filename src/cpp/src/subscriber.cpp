#include "subscriber.hpp"
#include <iostream>
#include <nlohmann/json.hpp>
#include <unordered_map>
#include <variant>

using json = nlohmann::json;


Subscriber::Subscriber(const std::string &address)
    : context(1), socket(context, zmq::socket_type::sub) {
    socket.connect(address);
    socket.set(zmq::sockopt::subscribe, "");
    std::cout << "✅ Subscriber C++ conectado a " << address << std::endl;
}

std::pair<std::vector<cv::Mat>, DataMap> Subscriber::receive_message() {
    std::string topic_str;
    std::vector<std::string> chunks;
    int num_chunks_int = -1;

    while (true) {
        zmq::message_t topic, index, num_chunks, chunk;

        auto res1 = socket.recv(topic, zmq::recv_flags::none);
        auto res2 = socket.recv(index, zmq::recv_flags::none);
        auto res3 = socket.recv(num_chunks, zmq::recv_flags::none);
        auto res4 = socket.recv(chunk, zmq::recv_flags::none);

        if (!res1 || !res2 || !res3 || !res4) {
            throw std::runtime_error("❌ Error recibiendo partes del mensaje multipart");
        }

        // Extraer los valores como strings
        std::string index_str(static_cast<char*>(index.data()), index.size());
        std::string num_chunks_str(static_cast<char*>(num_chunks.data()), num_chunks.size());

        int index_int = std::stoi(index_str);
        num_chunks_int = std::stoi(num_chunks_str);

        if (topic_str.empty()) {
            topic_str.assign(static_cast<char*>(topic.data()), topic.size());
            chunks.resize(num_chunks_int);
        }

        // Verificar que el mensaje es consistente (mismo topic y tamaño esperado)
        if (topic_str != std::string(static_cast<char*>(topic.data()), topic.size())) {
            throw std::runtime_error("❌ Error: Cambio de topic en mitad del mensaje.");
        }

        if (index_int < 0 || index_int >= num_chunks_int) {
            throw std::runtime_error("❌ Error: Índice fuera de rango.");
        }

        // Guardar el chunk en su posición
        chunks[index_int].assign(static_cast<char*>(chunk.data()), chunk.size());

        // Si ya tenemos todos los chunks, salir del bucle
        bool mensaje_completo = true;
        for (int i = 0; i < num_chunks_int; ++i) {
            if (chunks[i].empty()) {
                mensaje_completo = false;
                break;
            }
        }

        if (mensaje_completo) {
            break;
        }
    }

    // Concatenar todos los chunks en uno solo
    std::string messageBytes;
    for (const auto& chunk : chunks) {
        messageBytes += chunk;
    }

    // Ahora sigue igual que antes...
    size_t pos = messageBytes.find('\0');
    if (pos == std::string::npos) {
        throw std::runtime_error("❌ Error: separador '\\0' no encontrado en el mensaje");
    }

    std::string header_json = messageBytes.substr(0, pos);
    std::string images_data = messageBytes.substr(pos + 1);

    json message_json = json::parse(header_json);

    DataMap data;
    for (auto& [k, v] : message_json["data"].items()) {
        if (v.is_string()) {
            data[k] = v.get<std::string>();
        } else if (v.is_number_integer()) {
            data[k] = v.get<int>();
        } else if (v.is_number_float()) {
            data[k] = v.get<double>();
        } else if (v.is_boolean()) {
            data[k] = v.get<bool>();
        } else if (v.is_null()) {
            data[k] = nullptr;
        } else {
            data[k] = v.dump();
        }
    }

    std::vector<cv::Mat> images;
    int offset = 0;

    for (const auto& img_meta : message_json["images"]) {
        int width = img_meta["metadata"]["width"];
        int height = img_meta["metadata"]["height"];
        int channels = img_meta["metadata"]["channels"];
        int size = img_meta["metadata"]["size"];
        std::string dtype = img_meta["metadata"]["dtype"];

        int type = numpy_dtype_to_opencv(dtype, channels);

        if (offset + size > images_data.size()) {
            throw std::runtime_error("❌ Error: Datos de imagen incompletos.");
        }

        cv::Mat img(height, width, type);
        std::memcpy(img.data, images_data.data() + offset, size);
        offset += size;

        images.push_back(img);
    }

    return {images, data};
}

int Subscriber::numpy_dtype_to_opencv(const std::string& dtype, int channels) {
    int type;
    if (dtype == "uint8") type = CV_8UC(channels);
    else if (dtype == "int8") type = CV_8SC(channels);
    else if (dtype == "uint16") type = CV_16UC(channels);
    else if (dtype == "int16") type = CV_16SC(channels);
    else if (dtype == "float32") type = CV_32FC(channels);
    else if (dtype == "float64") type = CV_64FC(channels);
    else throw std::runtime_error("Tipo de imagen no soportado: " + dtype);
    return type;
}
