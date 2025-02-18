#include "subscriber.hpp"

int main() {
    Subscriber sub("tcp://integration_cpp_pub:5559");

    for (int i = 0; i < 1; i++) {
        sub.receive();
    }

    return 0;
}
