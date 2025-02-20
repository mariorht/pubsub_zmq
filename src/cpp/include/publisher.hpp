#ifndef PUBLISHER_HPP
#define PUBLISHER_HPP

#include <zmq.hpp>
#include <string>
#include <vector>
#include <map>
#include <opencv2/opencv.hpp>


class Publisher {
public:
    Publisher(const std::string& address, const std::string& topic = "camera/image_raw", size_t chunk_size = 100000);

    std::vector<std::string> build_message(const std::vector<cv::Mat>& frames, const std::map<std::string, std::string>& data);
    void publish_message(const std::vector<std::string>& message_chunks);

private:
    zmq::context_t context;
    zmq::socket_t socket;

    std::string topic;
    size_t chunk_size;
};

#endif // PUBLISHER_HPP
