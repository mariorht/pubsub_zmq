from publisher import Publisher
import time
import json  # Add import for JSON handling
import cv2  # Add import for OpenCV

if __name__ == "__main__":
    pub = Publisher()
    angle = 0  # Initialize rotation angle
    try:
        # Load image
        frame = cv2.imread("../../assets/pong.png")
        assert frame is not None, "❌ No se pudo cargar la imagen: pong.png"
            
        while True:
            frames = []
            for _ in range(3):  # Add 3 rotated images
                angle = (angle + 10) % 360  # Increment angle
                center = (frame.shape[1] // 2, frame.shape[0] // 2)
                matrix = cv2.getRotationMatrix2D(center, angle, 1.0)
                rotated_frame = cv2.warpAffine(frame, matrix, (frame.shape[1], frame.shape[0]))
                frames.append(rotated_frame)

            # Additional data
            data = {"key": "probando"}

            # Build and publish message with multiple images and additional data
            message_bytes = pub.build_message(frames, data)
            pub.publish_message(message_bytes)

            time.sleep(.1)
            
    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
        pub.close()
