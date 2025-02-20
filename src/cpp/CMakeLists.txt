cmake_minimum_required(VERSION 3.10)
project(PubSubZMQ)

set(CMAKE_CXX_STANDARD 17)

# Buscar ZeroMQ
include_directories(/usr/include)
link_libraries(zmq)

# Incluir directorios
include_directories(include)

# Añadir ejecutables
add_executable(main_publisher tests/integration/main_publisher.cpp src/publisher.cpp)
add_executable(main_subscriber tests/integration/main_subscriber.cpp src/subscriber.cpp)

# Linkear con ZeroMQ
target_link_libraries(main_publisher zmq)
target_link_libraries(main_subscriber zmq)


enable_testing()

add_executable(test_pubsub tests/unit/test_pubsub.cpp src/publisher.cpp src/subscriber.cpp)
target_link_libraries(test_pubsub zmq gtest gtest_main pthread)

add_test(NAME PubSubTests COMMAND test_pubsub)

