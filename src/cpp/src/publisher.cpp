#include "publisher.hpp"
#include <iostream>

Publisher::Publisher(const std::string &address)
    : context(1), socket(context, ZMQ_PUB) {
    socket.bind(address);
    std::cout << "✅ Publisher C++ conectado en " << address << std::endl;
}

void Publisher::publish(const std::string &message) {
    zmq::message_t zmqMessage(message.begin(), message.end());
    socket.send(zmqMessage, zmq::send_flags::none);
    std::cout << "📤 Mensaje enviado: " << message << std::endl;
}
