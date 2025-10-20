### Related Repositories
[Digital](https://github.com/rtfmtom/Digital)  
[CPU](https://github.com/rtfmtom/CPU)  
[Assembler](https://github.com/rtfmtom/Assembler)  

# Life-Hack

This application demonstrates remotely connecting to Digital circuit simulator and running Conway's Game of Life on the Hack [CPU](https://github.com/rtfmtom/CPU), with a GUI for visualization.

Currently requires my custom fork of [Digital](https://github.com/rtfmtom/Digital).

<img width="1321" height="861" alt="Screenshot 2025-10-18 at 2 02 15 PM" src="https://github.com/user-attachments/assets/534c771c-e2b8-486b-ba49-3dcc65041405" />

## Prerequisites

### Digital

As mentioned above, you will need my custom fork of Digital. This is because I modified the remote interface to allow remotely accessing memory on the `HackComputer` chip.[^1]

[^1]: My current implementation specifically targets the RAMDualPort chip on the running circuit, so it's tightly coupled with this demo—but could easily be generalized.

Clone my custom fork of Digital from [https://github.com/rtfmtom/Digital](https://github.com/rtfmtom/Digital):
```bash
git clone https://github.com/rtfmtom/Digital.git && cd Digital
mvn clean install
```

Set the `DIGITAL_PATH` environment variable:
```bash
export DIGITAL_PATH="/path/to/Digital/target/Digital.jar"
```

Or add it to your shell profile (`.bashrc`, `.zshrc`, etc.) to make it permanent:
```bash
echo 'export DIGITAL_PATH="/path/to/Digital/target/Digital.jar"' >> ~/.bashrc
```

### Hack CPU

Once you have the modified Digital installed, you'll need the `.dig` files for the `HackComputer` circuit.

Clone the Hack CPU from [https://github.com/rtfmtom/CPU](https://github.com/rtfmtom/CPU):
```bash
git clone https://github.com/rtfmtom/CPU
```

This repository contains all the `.dig` circuit files for the simulated Hack computer.

Set the `CIRCUIT_PATH` environment variable:
```bash
export CIRCUIT_PATH="/path/to/CPU/circuits/HackComputer.dig"
```

Or add it to your shell profile (`.bashrc`, `.zshrc`, etc.) to make it permanent:
```bash
echo 'export CIRCUIT_PATH="/path/to/CPU/circuits/HackComputer.dig"' >> ~/.bashrc
```

### Life-Hack

Clone the repository and download dependencies:
```bash
git clone https://github.com/rtfmtom/Life-Hack.git
cd Life-Hack
go mod download
```

## Usage

Run the application with default settings (uses `DIGITAL_PATH` and `CIRCUIT_PATH` environment variables if set, or specify paths via command-line flags):
```bash
go run .
```

The default configuration performs the following steps:

1. Launches Digital with the specified circuit file
2. Loads `example/Conway32.hex` (relative to the repository root) into the Hack CPU's ROM
3. Sends the start command to begin program execution
4. Opens the GUI and continuously reads memory addresses `0x2000-0x23FF` to visualize the Game of Life state

### Building and Using Custom Paths

You can build the application and run it with custom paths to override the defaults:
```bash
go build
./life-hack -d /path/to/Digital.jar -c /path/to/Circuit.dig -p /path/to/program.hex
```

## Command-line Flags

| Flag        | Short | Description                               | Default                          |
|-------------|-------|-------------------------------------------|----------------------------------|
| `--digital` | `-d`  | Path to Digital executable (.jar or .exe) | `$DIGITAL_PATH` env var          |
| `--circuit` | `-c`  | Path to `.dig` circuit file to simulate   | `$CIRCUIT_PATH` env var          |
| `--program` | `-p`  | Path to `.hex` file to run                | `example/Conway32.hex`           |
