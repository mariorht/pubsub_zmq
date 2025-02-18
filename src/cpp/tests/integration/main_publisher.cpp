#include "publisher.hpp"
#include <thread>
#include <chrono>

int main() {
    Publisher pub("tcp://*:5559");

    using namespace std::chrono_literals;
    std::this_thread::sleep_for(2000ms);

    for (int i = 0; i < 1; i++) {
        std::string msg = "Mensaje C++ " + std::to_string(i);
        pub.publish(msg);
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    return 0;
}
