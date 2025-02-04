from publisher import Publisher
import time
import json  # Add import for JSON handling
import cv2  # Add import for OpenCV

if __name__ == "__main__":
    pub = Publisher()
    angle = 0  # Initialize rotation angle
    try:
        # Load and rotate image
        frame = cv2.imread("pong.png")
        assert frame is not None, "❌ No se pudo cargar la imagen: pong.png"
            
        while True:

            angle = (angle + 10) % 360  # Increment angle
            center = (frame.shape[1] // 2, frame.shape[0] // 2)
            matrix = cv2.getRotationMatrix2D(center, angle, 1.0)
            rotated_frame = cv2.warpAffine(frame, matrix, (frame.shape[1], frame.shape[0]))

            # Publish rotated image
            pub.publish_image(rotated_frame)
            
            # Publish JSON
            json_message = json.dumps({"key": "probando"}).encode('utf-8')  # Encode JSON message to bytes
            pub.publish_message(json_message)

            time.sleep(.01)

            
    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
        pub.close()
