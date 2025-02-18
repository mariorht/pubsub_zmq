#ifndef SUBSCRIBER_HPP
#define SUBSCRIBER_HPP

#include <zmq.hpp>
#include <string>

class Subscriber {
public:
    Subscriber(const std::string &address);
    std::string receive();
private:
    zmq::context_t context;
    zmq::socket_t socket;
};

#endif // SUBSCRIBER_HPP
