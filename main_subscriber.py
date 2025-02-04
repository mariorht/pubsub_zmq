from subscriber import Subscriber

if __name__ == "__main__":
    sub = Subscriber()

    try:
        sub.receive_image()  # Espera y recibe imágenes
    except KeyboardInterrupt:
        print("\n🛑 Suscriptor detenido.")
        sub.close()
