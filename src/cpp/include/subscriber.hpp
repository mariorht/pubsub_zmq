#ifndef SUBSCRIBER_HPP
#define SUBSCRIBER_HPP

#include <zmq.hpp>
#include <string>
#include <opencv2/opencv.hpp>


class Subscriber {
public:
    Subscriber(const std::string &address);
    std::pair<std::vector<cv::Mat>, std::map<std::string, std::string>> receive_message();
    
private:
    zmq::context_t context;
    zmq::socket_t socket;
};

#endif // SUBSCRIBER_HPP
