#include "publisher.hpp"
#include <iostream>
#include <sstream>
#include <nlohmann/json.hpp>
#include <cstring>

using json = nlohmann::json;

Publisher::Publisher(const std::string& address, const std::string& topic, size_t chunk_size)
    : context(1), socket(context, zmq::socket_type::pub), topic(topic), chunk_size(chunk_size) {
    socket.bind(address);
    std::cout << "✅ Publisher C++ conectado en " << address << " con topic: " << topic << std::endl;
}

std::vector<std::string> Publisher::build_message(const std::map<std::string, std::string>& data) {
    json message_json;
    message_json["type"] = "images";  // De momento seguimos diciendo que es "images" para mantener el protocolo
    message_json["count"] = 0;        // No hay imágenes
    message_json["images"] = json::array();  // Lista vacía
    message_json["data"] = data;

    std::string message_str = message_json.dump();
    std::vector<std::string> chunks;

    for (size_t i = 0; i < message_str.size(); i += chunk_size) {
        chunks.push_back(message_str.substr(i, chunk_size));
    }

    return chunks;
}

void Publisher::publish_message(const std::vector<std::string>& message_chunks) {
    size_t num_chunks = message_chunks.size();

    for (size_t i = 0; i < num_chunks; ++i) {
        zmq::message_t topic_msg(topic.begin(), topic.end());
        zmq::message_t index_msg(std::to_string(i));
        zmq::message_t num_chunks_msg(std::to_string(num_chunks));
        zmq::message_t chunk_msg(message_chunks[i].begin(), message_chunks[i].end());

        socket.send(topic_msg, zmq::send_flags::sndmore);
        socket.send(index_msg, zmq::send_flags::sndmore);
        socket.send(num_chunks_msg, zmq::send_flags::sndmore);
        socket.send(chunk_msg, zmq::send_flags::none);
    }

    std::cout << "📤 Mensaje fragmentado en " << num_chunks << " partes enviado." << std::endl;
}
