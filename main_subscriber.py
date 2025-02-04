from subscriber import Subscriber
import time

if __name__ == "__main__":
    sub = Subscriber()
    try:
        while True:
            start_time = time.time()  # Reset start time for each image
            sub.total_bytes_received = 0  # Reset total bytes received for each image
            frame = sub.receive_image()
            if frame is not None:
                elapsed_time = time.time() - start_time
                mbps = (sub.total_bytes_received * 8) / (elapsed_time * 1_000_000)
                print(f"🚀 Velocidad de recepción media: {mbps:.2f} Mbps")
    except KeyboardInterrupt:
        print("\n🛑 Suscriptor detenido.")
        sub.close()
