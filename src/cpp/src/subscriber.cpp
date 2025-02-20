#include "subscriber.hpp"
#include <iostream>
#include <nlohmann/json.hpp>

using json = nlohmann::json;


Subscriber::Subscriber(const std::string &address)
    : context(1), socket(context, zmq::socket_type::sub) {
    socket.connect(address);
    socket.set(zmq::sockopt::subscribe, "");
    std::cout << "✅ Subscriber C++ conectado a " << address << std::endl;
}

std::pair<std::vector<cv::Mat>, std::map<std::string, std::string>> Subscriber::receive_message() {
    zmq::message_t topic;
    zmq::message_t index;
    zmq::message_t num_chunks;
    zmq::message_t chunk;

    auto res1 = socket.recv(topic, zmq::recv_flags::none);
    auto res2 = socket.recv(index, zmq::recv_flags::none);
    auto res3 = socket.recv(num_chunks, zmq::recv_flags::none);
    auto res4 = socket.recv(chunk, zmq::recv_flags::none);
    
    if (!res1 || !res2 || !res3 || !res4) {
        throw std::runtime_error("❌ Error recibiendo partes del mensaje multipart");
    }
    

    std::string chunk_data(static_cast<char*>(chunk.data()), chunk.size());

    // Aquí sería donde podrías reconstruir el mensaje si fueran múltiples partes
    // Suponemos que es una sola parte
    size_t pos = chunk_data.find('\0');
    std::string header_json = chunk_data.substr(0, pos);
    std::string images_data = chunk_data.substr(pos + 1);

    json message_json = json::parse(header_json);

    std::map<std::string, std::string> data;
    for (auto& [k, v] : message_json["data"].items()) {
        data[k] = v;
    }

    std::vector<cv::Mat> images;
    int offset = 0;
    for (const auto& img_meta : message_json["images"]) {
        int width = img_meta["metadata"]["width"];
        int height = img_meta["metadata"]["height"];
        int channels = img_meta["metadata"]["channels"];
        int size = img_meta["metadata"]["size"];

        cv::Mat img(height, width, CV_8UC(channels));
        std::memcpy(img.data, images_data.data() + offset, size);
        offset += size;
        images.push_back(img);
    }

    return {images, data};
}
