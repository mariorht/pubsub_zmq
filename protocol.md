# PubSub ZMQ Protocol

This document describes the protocol used by PubSub ZMQ to encode and decode messages containing multiple images and additional data.

## Message Structure

A message consists of two main parts:
1. JSON Metadata
2. Binary Image Data

These parts are concatenated with a null byte (`\x00`) separator.

### JSON Metadata

The JSON metadata contains information about the images and any additional data. It has the following structure:

```json
{
    "type": "images",
    "count": <number_of_images>,
    "images": [
        {
            "metadata": {
                "format": "jpg",
                "width": <image_width>,
                "height": <image_height>,
                "size": <image_size_in_bytes>
            }
        },
        ...
    ],
    "data": <additional_data>
}
```

### Binary Image Data

The binary image data is a concatenation of the encoded images. Each image is encoded in the format specified in the JSON metadata (e.g., JPEG).

## Message Fragmentation

To handle large messages, the message is fragmented into smaller chunks. Each chunk is sent as a multipart message with the following parts:

1. Topic
2. Chunk Index
3. Total Number of Chunks
4. Chunk Data

### Example

For a message with 3 images and additional data, the process is as follows:

1. Construct the JSON metadata.
2. Encode the images and concatenate them.
3. Combine the JSON metadata and binary image data with a null byte separator.
4. Fragment the combined message into chunks.
5. Send each chunk as a multipart message.

### Receiving Messages

The subscriber receives the chunks and reconstructs the original message by:

1. Collecting all chunks.
2. Ordering the chunks by their index.
3. Combining the ordered chunks into a single byte stream.
4. Splitting the byte stream into JSON metadata and binary image data.
5. Decoding the JSON metadata.
6. Extracting and decoding each image based on the metadata.

## Error Handling

If the JSON metadata cannot be decoded, the message is discarded, and an error is logged.

## Example Code in python

### Publisher

```python
def build_message(self, frames, data):
    images_metadata = []
    images_data = []
    for frame in frames:
        _, buffer = cv2.imencode(".jpg", frame)
        image_bytes = buffer.tobytes()
        image_metadata = {
            "format": "jpg",
            "width": frame.shape[1],
            "height": frame.shape[0],
            "size": len(image_bytes)
        }
        images_metadata.append({"metadata": image_metadata})
        images_data.append(image_bytes)

    message = {
        "type": "images",
        "count": len(frames),
        "images": images_metadata,
        "data": data
    }
    message_json = json.dumps(message).encode('utf-8')
    return message_json + b'\x00' + b''.join(images_data)
```

### Subscriber

```python
def receive_message(self):
    while True:
        topic, index, total, chunk = self.socket.recv_multipart()
        index = int(index.decode())
        total = int(total.decode())

        self.image_chunks[index] = chunk

        if len(self.image_chunks) == total:
            ordered_chunks = [self.image_chunks[i] for i in sorted(self.image_chunks.keys())]
            message_bytes = b"".join(ordered_chunks)

            json_part, images_bytes = message_bytes.split(b'\x00', 1)

            try:
                message = json.loads(json_part.decode('utf-8'))
                images = []
                offset = 0
                for image_info in message["images"]:
                    size = image_info["metadata"]["size"]
                    image_bytes = images_bytes[offset:offset + size]
                    image_array = np.frombuffer(image_bytes, dtype=np.uint8)
                    frame = cv2.imdecode(image_array, cv2.IMREAD_COLOR)
                    images.append(frame)
                    offset += size
                self.image_chunks.clear()
                return images, message["data"]
            except json.JSONDecodeError:
                self.image_chunks.clear()
                return None, None
```