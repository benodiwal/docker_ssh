# docker_ssh
This app allows users to SSH into a system where the app is running and provides direct access to Docker on that machine.

## Features
- **SSH Access**: Securely SSH into the host system.
- **Docker Management**: Run Docker commands directly after connecting via SSH.

## Prerequisites
- **Go**: Ensure Go is installed on your system.
- **Docker**: Docker must be installed and running.

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/benodiwal/docker_ssh.git
   cd docker_ssh
   ```
   
3. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
   
5. Build the project:
   ```bash
   go build -o docker_ssh cmd/docker_ssh/main.go
   ```
   
7. Run the application:
   ```bash
   ./docker_ssh
   ```

   ## Usage
   1. Run the docker_ssh binary.
   2. SSH into the system using:
      ```bash
      ssh user@host
      ```

   ## Contributing
   Contributions are welcome! Please submit any issues or pull requests.

   ## License
   This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
