from publisher import Publisher
import time
import json  # Add import for JSON handling

if __name__ == "__main__":
    pub = Publisher()
    try:
        while True:
            # Publish image
            pub.publish_image(image_path="pong.png")
            time.sleep(.1)
            
            # Publish JSON
            json_message = json.dumps({"key": "probando"})
            pub.publish_message(json_message)
            time.sleep(.1)
    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
        pub.close()
