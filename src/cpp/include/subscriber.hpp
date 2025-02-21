#ifndef SUBSCRIBER_HPP
#define SUBSCRIBER_HPP

#include <zmq.hpp>
#include <string>
#include <variant>
#include <opencv2/opencv.hpp>



using ValueVariant = std::variant<std::string, int, double, bool, std::nullptr_t>;
using DataMap = std::unordered_map<std::string, ValueVariant>;

class Subscriber {
public:
    Subscriber(const std::string &address);
    std::pair<std::vector<cv::Mat>, DataMap> receive_message();
    int numpy_dtype_to_opencv(const std::string& dtype, int channels);

    
private:
    zmq::context_t context;
    zmq::socket_t socket;
};

#endif // SUBSCRIBER_HPP
