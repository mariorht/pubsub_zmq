#include "publisher.hpp"
#include "subscriber.hpp"
#include <gtest/gtest.h>
#include <thread>
#include <chrono>
#include <sstream>
#include <vector>
#include <map>
#include <opencv2/opencv.hpp>
#include <nlohmann/json.hpp>


using namespace std::chrono_literals;
using json = nlohmann::json;


TEST(PubSubTest, MensajeUnico) {
    const std::string endpoint = "tcp://127.0.0.1:5555";
    const std::string topic = "test";
    const size_t chunk_size = 1024;

    Publisher pub(endpoint, topic, chunk_size);
    Subscriber sub(endpoint);

    std::this_thread::sleep_for(500ms); // Espera para asegurar conexión

    std::string mensajeEsperado = "Mensaje de prueba";
    std::map<std::string, std::string> data = {{"msg", mensajeEsperado}};

    std::thread publicador([&]() {
        std::this_thread::sleep_for(100ms);
        auto chunks = pub.build_message({}, data); // Sin imágenes
        pub.publish_message(chunks);
    });

    std::pair<std::vector<cv::Mat>, std::map<std::string, std::string>> mensajeRecibido;
    bool recibido = false;

    for (int i = 0; i < 5; i++) {
        try {
            mensajeRecibido = sub.receive_message();
            recibido = true;
            break;
        } catch (const std::exception& e) {
            std::this_thread::sleep_for(100ms);
        }
    }

    publicador.join();

    ASSERT_TRUE(recibido) << "No se recibió el mensaje después de varios intentos";
    ASSERT_EQ(mensajeRecibido.second["msg"], mensajeEsperado) << "El mensaje recibido no coincide con el esperado";
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
                auto mensaje = sub.receive_message();
                mensajesRecibidos[i] = mensaje.second["msg"];
            } catch (...) {
                mensajesRecibidos[i] = "";
            }
        }
    });

    for (int i = 0; i < totalMensajes; i++) {
        std::ostringstream ss;
        ss << "Mensaje " << i;
        std::map<std::string, std::string> data = {{"msg", ss.str()}};

        auto chunks = pub.build_message({}, data); // Sin imágenes
        pub.publish_message(chunks);

        std::this_thread::sleep_for(100ms);
    }

    suscriptor.join();

    for (int i = 0; i < totalMensajes; i++) {
        ASSERT_FALSE(mensajesRecibidos[i].empty()) << "El mensaje " << i << " no fue recibido";
        std::ostringstream esperado;
        esperado << "Mensaje " << i;
        ASSERT_EQ(mensajesRecibidos[i], esperado.str()) << "El mensaje recibido no coincide con el esperado";
    }
}





TEST(PubSubTest, EnviarYRecibirImagenReal) {
    Publisher pub("tcp://127.0.0.1:5555");
    Subscriber sub("tcp://127.0.0.1:5555");

    // Imagen 10x10 con un degradado
    cv::Mat imagen_real(10, 10, CV_8UC3);
    for (int i = 0; i < imagen_real.rows; ++i) {
        for (int j = 0; j < imagen_real.cols; ++j) {
            imagen_real.at<cv::Vec3b>(i, j) = cv::Vec3b(i * 25, j * 25, (i + j) * 12);
        }
    }

    std::vector<cv::Mat> frames = {imagen_real};
    std::map<std::string, std::string> data = {{"tipo", "imagen_real_test"}};

    // Enviar mensaje
    auto chunks = pub.build_message(frames, data);

    // Simular recepción de multipart (topic, index, num_chunks, chunk)
    // Vamos a juntar todo el mensaje de golpe (simulación simplificada)
    std::string recibido;
    for (const auto& chunk : chunks) {
        recibido += chunk;
    }

    // Procesar como si fuera recibido en el Subscriber
    size_t pos = recibido.find('\0');
    ASSERT_NE(pos, std::string::npos) << "No se encontró el separador \\0 en el mensaje";

    std::string header_json = recibido.substr(0, pos);
    std::string images_data = recibido.substr(pos + 1);

    // Deserializar el header
    nlohmann::json message_json = nlohmann::json::parse(header_json);

    // Simular reconstrucción de imágenes
    std::vector<cv::Mat> images;
    int offset = 0;
    for (const auto& img_meta : message_json["images"]) {
        int width = img_meta["metadata"]["width"];
        int height = img_meta["metadata"]["height"];
        int channels = img_meta["metadata"]["channels"];
        int size = img_meta["metadata"]["size"];

        cv::Mat img(height, width, CV_8UC(channels));
        std::memcpy(img.data, images_data.data() + offset, size);
        offset += size;
        images.push_back(img);
    }

    // Validaciones
    ASSERT_EQ(images.size(), 1);
    ASSERT_EQ(images[0].cols, imagen_real.cols);
    ASSERT_EQ(images[0].rows, imagen_real.rows);
    ASSERT_EQ(images[0].channels(), imagen_real.channels());

    // Compara el contenido pixel a pixel
    ASSERT_EQ(cv::norm(images[0], imagen_real, cv::NORM_L1), 0) << "La imagen reconstruida no coincide con la original";
}