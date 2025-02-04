from publisher import Publisher
import time

if __name__ == "__main__":
    pub = Publisher()
    try:
        while True:
            pub.publish_image(image_path="pong.png")
            time.sleep(.1)
    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
        pub.close()
