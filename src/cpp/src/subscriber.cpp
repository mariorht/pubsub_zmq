#include "subscriber.hpp"
#include <iostream>

Subscriber::Subscriber(const std::string &address)
    : context(1), socket(context, zmq::socket_type::sub) {
    socket.connect(address);
    socket.set(zmq::sockopt::subscribe, "");
    std::cout << "✅ Subscriber C++ conectado a " << address << std::endl;
}

std::string Subscriber::receive() {
    zmq::message_t topic;
    zmq::message_t index;
    zmq::message_t num_chunks;
    zmq::message_t chunk;

    auto result = socket.recv(topic, zmq::recv_flags::none);
    if (!result) throw std::runtime_error("❌ Error al recibir topic");

    result = socket.recv(index, zmq::recv_flags::none);
    if (!result) throw std::runtime_error("❌ Error al recibir index");

    result = socket.recv(num_chunks, zmq::recv_flags::none);
    if (!result) throw std::runtime_error("❌ Error al recibir num_chunks");

    result = socket.recv(chunk, zmq::recv_flags::none);
    if (!result) throw std::runtime_error("❌ Error al recibir chunk");

    std::string received(static_cast<char*>(chunk.data()), chunk.size());
    std::cout << "📥 Chunk recibido: " << received << std::endl;

    return received; // Devuelve el contenido del chunk (parte del mensaje real)
}
