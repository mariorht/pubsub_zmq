#include "publisher.hpp"
#include "subscriber.hpp"
#include <gtest/gtest.h>
#include <thread>
#include <chrono>
#include <sstream>
#include <vector>
#include <unordered_map>
#include <opencv2/opencv.hpp>
#include <nlohmann/json.hpp>
#include <variant>

using namespace std::chrono_literals;
using json = nlohmann::json;

// Definiciones para simplificar el acceso a std::variant
using ValueVariant = std::variant<std::string, int, double, bool, std::nullptr_t>;
using DataMap = std::unordered_map<std::string, ValueVariant>;

// Función auxiliar para obtener std::string de ValueVariant
std::string get_string(const DataMap& data, const std::string& key) {
    auto it = data.find(key);
    if (it == data.end()) {
        throw std::runtime_error("Clave no encontrada: " + key);
    }
    if (!std::holds_alternative<std::string>(it->second)) {
        throw std::runtime_error("Valor de '" + key + "' no es string");
    }
    return std::get<std::string>(it->second);
}

// ======================================
// TEST: Envío y recepción de mensaje único con data
// ======================================
TEST(PubSubTest, MensajeUnico) {
    const std::string endpoint = "tcp://127.0.0.1:5555";
    const std::string topic = "test";
    const size_t chunk_size = 1024;

    Publisher pub(endpoint, topic, chunk_size);
    Subscriber sub(endpoint, topic);

    std::this_thread::sleep_for(500ms); // Espera para asegurar conexión

    std::string mensajeEsperado = "Mensaje de prueba";
    std::map<std::string, std::string> data = {{"msg", mensajeEsperado}};

    std::thread publicador([&]() {
        std::this_thread::sleep_for(100ms);
        auto chunks = pub.build_message({}, data); // Sin imágenes
        pub.publish_message(chunks);
    });

    std::pair<std::vector<cv::Mat>, DataMap> mensajeRecibido;
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
    ASSERT_EQ(get_string(mensajeRecibido.second, "msg"), mensajeEsperado);
}

// ======================================
// TEST: Envío y recepción de múltiples mensajes
// ======================================
TEST(PubSubTest, MultiplesMensajes) {
    const std::string endpoint = "tcp://127.0.0.1:5560";
    const std::string topic = "test";
    const size_t chunk_size = 1024;

    Publisher pub(endpoint, topic, chunk_size);
    Subscriber sub(endpoint, topic);

    std::this_thread::sleep_for(500ms);

    const int totalMensajes = 5;
    std::vector<std::string> mensajesRecibidos(totalMensajes);

    std::thread suscriptor([&]() {
        for (int i = 0; i < totalMensajes; i++) {
            try {
                auto mensaje = sub.receive_message(3000);
                mensajesRecibidos[i] = get_string(mensaje.second, "msg");
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
        ASSERT_EQ(mensajesRecibidos[i], esperado.str());
    }
}

// ======================================
// TEST: Enviar y recibir imágenes reales
// ======================================
TEST(PubSubTest, EnviarYRecibirImagenReal) {
    Publisher pub("tcp://127.0.0.1:5555", "test");
    Subscriber sub("tcp://127.0.0.1:5555", "test");

    std::this_thread::sleep_for(500ms);


    // Imagen 10x10 con un degradado
    cv::Mat imagen_real(10, 10, CV_8UC3);
    for (int i = 0; i < imagen_real.rows; ++i) {
        for (int j = 0; j < imagen_real.cols; ++j) {
            imagen_real.at<cv::Vec3b>(i, j) = cv::Vec3b(i * 25, j * 25, (i + j) * 12);
        }
    }

    std::vector<cv::Mat> frames = {imagen_real};
    std::map<std::string, std::string> data = {{"tipo", "imagen_real_test"}};

    std::thread publicador([&]() {
        std::this_thread::sleep_for(1000ms);
        auto chunks = pub.build_message(frames, data);
        pub.publish_message(chunks);
    });

    std::pair<std::vector<cv::Mat>, DataMap> mensajeRecibido;
    bool recibido = false;


    try {
        mensajeRecibido = sub.receive_message();
        recibido = true;
    } catch (std::runtime_error e) {
        std::cerr << "Error al recibir mensaje: " << e.what() << std::endl;
    }
    

    publicador.join();

    ASSERT_TRUE(recibido) << "No se recibió el mensaje con imágenes después de varios intentos";
    ASSERT_EQ(mensajeRecibido.first.size(), 1);
    ASSERT_EQ(mensajeRecibido.first[0].cols, imagen_real.cols);
    ASSERT_EQ(mensajeRecibido.first[0].rows, imagen_real.rows);
    ASSERT_EQ(mensajeRecibido.first[0].channels(), imagen_real.channels());
    ASSERT_EQ(cv::norm(mensajeRecibido.first[0], imagen_real, cv::NORM_L1), 0)
        << "La imagen reconstruida no coincide con la original";

    // Si quisieras comprobar el campo 'tipo' en data:
    ASSERT_EQ(get_string(mensajeRecibido.second, "tipo"), "imagen_real_test");
}


TEST(PubSubTest, EnviarYRecibirImagenPNG) {
    Publisher pub("tcp://127.0.0.1:5551", "test",1024);
    Subscriber sub("tcp://127.0.0.1:5551", "test");

    // Crear imagen de prueba
    cv::Mat imagen_png(20, 20, CV_8UC3, cv::Scalar(50, 100, 150));

    std::vector<cv::Mat> frames = {imagen_png};
    std::map<std::string, std::string> data = {{"formato", "png"}};

    std::thread publicador([&]() {
        std::this_thread::sleep_for(100ms);
        auto chunks = pub.build_message(frames, data, "png"); // PNG
        pub.publish_message(chunks);
    });

    std::pair<std::vector<cv::Mat>, DataMap> mensajeRecibido;
    bool recibido = false;

    for (int i = 0; i < 5; ++i) {
        try {
            mensajeRecibido = sub.receive_message();
            recibido = true;
            break;
        } catch (...) {
            std::this_thread::sleep_for(100ms);
        }
    }

    publicador.join();

    ASSERT_TRUE(recibido) << "No se recibió el mensaje PNG después de varios intentos";
    ASSERT_EQ(mensajeRecibido.first.size(), 1);
    ASSERT_EQ(mensajeRecibido.first[0].cols, imagen_png.cols);
    ASSERT_EQ(mensajeRecibido.first[0].rows, imagen_png.rows);
    ASSERT_EQ(mensajeRecibido.first[0].channels(), imagen_png.channels());
    ASSERT_EQ(cv::norm(mensajeRecibido.first[0], imagen_png, cv::NORM_L1), 0)
        << "La imagen PNG reconstruida no coincide con la original";

    ASSERT_EQ(get_string(mensajeRecibido.second, "formato"), "png");
}


TEST(PubSubTest, EnviarYRecibirImagenJPEG) {
    Publisher pub("tcp://127.0.0.1:5555", "test");
    Subscriber sub("tcp://127.0.0.1:5555", "test");

    // Crear imagen de prueba
    cv::Mat imagen_jpeg(20, 20, CV_8UC3, cv::Scalar(100, 150, 200));

    std::vector<cv::Mat> frames = {imagen_jpeg};
    std::map<std::string, std::string> data = {{"formato", "jpeg"}};

    std::thread publicador([&]() {
        std::this_thread::sleep_for(100ms);
        auto chunks = pub.build_message(frames, data, "jpeg"); // JPEG
        pub.publish_message(chunks);
    });

    std::pair<std::vector<cv::Mat>, DataMap> mensajeRecibido;
    bool recibido = false;

    for (int i = 0; i < 5; ++i) {
        try {
            mensajeRecibido = sub.receive_message();
            recibido = true;
            break;
        } catch (...) {
            std::this_thread::sleep_for(100ms);
        }
    }

    publicador.join();

    ASSERT_TRUE(recibido) << "No se recibió el mensaje JPEG después de varios intentos";
    ASSERT_EQ(mensajeRecibido.first.size(), 1);
    ASSERT_EQ(mensajeRecibido.first[0].cols, imagen_jpeg.cols);
    ASSERT_EQ(mensajeRecibido.first[0].rows, imagen_jpeg.rows);
    ASSERT_EQ(mensajeRecibido.first[0].channels(), imagen_jpeg.channels());

    // Dado que JPEG es con pérdida, permitimos una ligera diferencia en la reconstrucción
    double diff = cv::norm(mensajeRecibido.first[0], imagen_jpeg, cv::NORM_L2);
    ASSERT_LT(diff, 50.0) << "La imagen JPEG reconstruida tiene demasiada diferencia respecto a la original";

    ASSERT_EQ(get_string(mensajeRecibido.second, "formato"), "jpeg");
}
