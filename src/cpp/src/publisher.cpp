#include "publisher.hpp"
#include <nlohmann/json.hpp>
#include <iostream>
#include <opencv2/imgcodecs.hpp>

using json = nlohmann::json;

Publisher::Publisher(const std::string& address, const std::string& topic, size_t chunk_size)
    : context(1), socket(context, zmq::socket_type::pub), topic(topic), chunk_size(chunk_size) {
    socket.bind(address);
    std::cout << "✅ Publisher C++ conectado en " << address << " con topic: " << topic << std::endl;
}

std::vector<std::string> Publisher::build_message(const std::vector<cv::Mat>& frames, const std::map<std::string, std::string>& data, const std::string& format) {
    json message_json;
    std::vector<std::string> images_data;
    json images_metadata = json::array();

    for (const auto& frame : frames) {
        std::vector<uchar> buffer;
        
        if (format == "jpeg") {
            cv::imencode(".jpg", frame, buffer);
        } else if (format == "png") {
            cv::imencode(".png", frame, buffer);
        } else { // RAW
            buffer.assign(frame.data, frame.data + frame.total() * frame.elemSize());
        }

        json meta;
        meta["format"] = format;
        meta["width"] = frame.cols;
        meta["height"] = frame.rows;
        meta["channels"] = frame.channels();
        meta["dtype"] = mat_type_to_dtype_string(frame.depth());
        meta["size"] = buffer.size();
        images_metadata.push_back({{"metadata", meta}});

        images_data.push_back(std::string(buffer.begin(), buffer.end()));
    }

    message_json["type"] = "images";
    message_json["count"] = frames.size();
    message_json["images"] = images_metadata;
    message_json["data"] = data;

    std::string header = message_json.dump();
    std::vector<std::string> chunks;

    std::string messageBytes = header;
    messageBytes.push_back('\0');
    for (const auto& img : images_data) {
        messageBytes += img;
    }

    for (size_t i = 0; i < messageBytes.size(); i += chunk_size) {
        chunks.push_back(messageBytes.substr(i, chunk_size));
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



std::string Publisher::mat_type_to_dtype_string(int mat_type) {
    switch (mat_type) {
        case CV_8U: return "uint8";
        case CV_8S: return "int8";
        case CV_16U: return "uint16";
        case CV_16S: return "int16";
        case CV_32S: return "int32";
        case CV_32F: return "float32";
        case CV_64F: return "float64";
        default: throw std::runtime_error("Tipo de cv::Mat no soportado: " + std::to_string(mat_type));
    }
}
