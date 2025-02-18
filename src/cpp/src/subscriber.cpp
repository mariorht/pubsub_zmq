#include "subscriber.hpp"
#include <iostream>
#include <stdexcept>

Subscriber::Subscriber(const std::string &address)
    : context(1), socket(context, zmq::socket_type::sub) {
    socket.connect(address);
    socket.set(zmq::sockopt::subscribe, "");
    std::cout << "✅ Subscriber C++ conectado a " << address << std::endl;
}

std::string Subscriber::receive() {
    zmq::message_t message;
    auto result = socket.recv(message, zmq::recv_flags::none);
    if (!result) {
        throw std::runtime_error("❌ Error al recibir el mensaje con ZeroMQ");
    }
    std::string received(static_cast<char*>(message.data()), message.size());
    std::cout << "📥 Mensaje recibido: " << received << std::endl;
    return received;
}
