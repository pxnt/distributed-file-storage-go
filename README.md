# Distributed File Storage (DFS)

A robust, peer-to-peer distributed file storage system written in Go, enabling reliable file storage and retrieval across a network of nodes.

## Features

- **Peer-to-peer Architecture**: Seamlessly distribute and replicate files across multiple nodes
- **Content-Addressable Storage**: Uses SHA-1 hashing for efficient file addressing and retrieval
- **Bootstrap Network**: Easily join existing networks by connecting to bootstrap nodes
- **Stream-based File Transfer**: Efficiently handles files of any size through streaming
- **Multiple Transport Options**: Built with TCP transport, extensible for other protocols

## Architecture

DFS follows a modular architecture:

- **Transport Layer**: Handles network communication between nodes
- **Storage Layer**: Manages local file persistence with content-addressable storage
- **Server Layer**: Coordinates file operations and peer communication
- **Domain Layer**: Defines message formats for inter-node communication

## Quick Start

### Prerequisites

- Go 1.22 or higher

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/distributed-file-storage.git
cd distributed-file-storage

# Build the project
make build
```

### Running Nodes

Start a bootstrap node:

```bash
./dfs s2
```

Start a regular node that connects to the bootstrap node:

```bash
./dfs s1
```

### Testing the Network

The default setup automatically tests file storage and retrieval:

1. Node 1 (s1) will store a test file with key "test1"
2. Node 2 (s2) will store a test file with key "test2"
3. The system will verify retrieval across nodes

## Usage

```go
// Create a file server
server := createServer(":3000", ":4000") // Connecting to bootstrap node on port 4000
server.Start()

// Store data
data := bytes.NewBuffer([]byte("Hello, distributed world!"))
server.StoreData("my_file_key", data)

// Retrieve data
reader, err := server.Get("my_file_key")
if err != nil {
    // Handle error
}
// Read file content from reader
```

## Development

### Project Structure

- **main.go**: Entry point and server initialization
- **server.go**: File server implementation
- **store.go**: Local storage implementation
- **p2p/**: Peer-to-peer networking components
- **domain/**: Message definitions
- **codec/**: Serialization tools

### Running Tests

```bash
make test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)

## Acknowledgments

- Inspired by modern distributed systems like IPFS and BitTorrent
- Built with pure Go with minimal dependencies 