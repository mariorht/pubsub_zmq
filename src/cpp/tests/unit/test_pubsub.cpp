#include "publisher.hpp"
#include "subscriber.hpp"
#include <gtest/gtest.h>
#include <thread>
#include <chrono>
#include <sstream>

using namespace std::chrono_literals;

TEST(PubSubTest, MensajeUnico) {
    const std::string endpoint = "tcp://127.0.0.1:5559";
    const std::string topic = "test";
    const size_t chunk_size = 1024;

    Publisher pub(endpoint, topic, chunk_size);
    Subscriber sub(endpoint);

    std::this_thread::sleep_for(500ms); // Espera para asegurar conexión

    std::string mensajeEsperado = "Mensaje de prueba";
    std::map<std::string, std::string> data = {{"msg", mensajeEsperado}};

    std::thread publicador([&]() {
        std::this_thread::sleep_for(100ms);
        auto chunks = pub.build_message(data);
        pub.publish_message(chunks);
    });

    std::string mensajeRecibido;
    bool recibido = false;

    for (int i = 0; i < 5; i++) {
        try {
            mensajeRecibido = sub.receive();
            recibido = true;
            break;
        } catch (const std::exception& e) {
            std::this_thread::sleep_for(100ms);
        }
    }

    publicador.join();

    ASSERT_TRUE(recibido) << "No se recibió el mensaje después de varios intentos";
    ASSERT_NE(mensajeRecibido.find(mensajeEsperado), std::string::npos) << "El mensaje recibido no contiene el esperado";
}

TEST(PubSubTest, MultiplesMensajes) {
    const std::string endpoint = "tcp://127.0.0.1:5560";
    const std::string topic = "test";
    const size_t chunk_size = 1024;

    Publisher pub(endpoint, topic, chunk_size);
    Subscriber sub(endpoint);

    std::this_thread::sleep_for(500ms);

    const int totalMensajes = 5;
    std::vector<std::string> mensajesRecibidos(totalMensajes);

    std::thread suscriptor([&]() {
        for (int i = 0; i < totalMensajes; i++) {
            try {
                mensajesRecibidos[i] = sub.receive();
            } catch (...) {
                mensajesRecibidos[i] = "";
            }
        }
    });

    for (int i = 0; i < totalMensajes; i++) {
        std::ostringstream ss;
        ss << "Mensaje " << i;
        std::map<std::string, std::string> data = {{"msg", ss.str()}};

        auto chunks = pub.build_message(data);
        pub.publish_message(chunks);

        std::this_thread::sleep_for(100ms);
    }

    suscriptor.join();

    for (int i = 0; i < totalMensajes; i++) {
        ASSERT_FALSE(mensajesRecibidos[i].empty()) << "El mensaje " << i << " no fue recibido";
        std::ostringstream esperado;
        esperado << "Mensaje " << i;
        ASSERT_NE(mensajesRecibidos[i].find(esperado.str()), std::string::npos) << "El mensaje recibido no contiene el esperado";
    }
}
