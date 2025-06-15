# 🌓 Mittnite: A Minimal Init System for Containers

![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)
![Latest Release](https://img.shields.io/github/v/release/fabro-acosta/mittnite)
![Downloads](https://img.shields.io/github/downloads/fabro-acosta/mittnite/total)

Welcome to **Mittnite**, a lightweight init system designed for use as a container entrypoint or an init system. This project allows you to manage your container's lifecycle with ease, using templated configuration files for a streamlined experience.

## 📦 Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Configuration](#configuration)
6. [Examples](#examples)
7. [Contributing](#contributing)
8. [License](#license)
9. [Contact](#contact)

## 📖 Introduction

In the world of containerization, having a reliable init system is crucial. **Mittnite** provides a simple solution for managing processes within your containers. With its templated configuration files, you can easily customize the behavior of your containers to suit your needs.

To get started, check out the [latest releases here](https://github.com/fabro-acosta/mittnite/releases). Download the necessary files and execute them to set up **Mittnite**.

## ✨ Features

- **Lightweight**: Designed to be minimal, ensuring quick startup times.
- **Templated Configurations**: Use templates to easily configure your container environments.
- **Easy Integration**: Works seamlessly with existing container orchestration tools.
- **Process Management**: Efficiently manage multiple processes within your container.
- **Cross-Platform**: Compatible with various container runtimes.

## ⚙️ Installation

To install **Mittnite**, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/fabro-acosta/mittnite.git
   ```

2. Navigate to the directory:

   ```bash
   cd mittnite
   ```

3. Build the project:

   ```bash
   make build
   ```

4. After building, you can find the executable in the `bin` directory.

5. For the latest version, visit the [releases page](https://github.com/fabro-acosta/mittnite/releases). Download the appropriate files, and execute them to install.

## 🚀 Usage

To use **Mittnite**, follow these steps:

1. Create a configuration file based on the provided templates.
2. Run the init system with the configuration file:

   ```bash
   ./mittnite -c path/to/config.yaml
   ```

3. Monitor the logs for any output:

   ```bash
   tail -f /var/log/mittnite.log
   ```

## ⚙️ Configuration

**Mittnite** uses YAML for configuration files. Below is an example configuration:

```yaml
services:
  - name: my-service
    command: /usr/bin/my-service
    restart: always
```

### Configuration Options

- **name**: The name of the service.
- **command**: The command to run the service.
- **restart**: Policy for restarting the service (e.g., `always`, `on-failure`).

## 🛠️ Examples

### Basic Example

Here’s a simple example to get you started:

1. Create a file named `my-service.yaml`:

   ```yaml
   services:
     - name: hello-world
       command: echo "Hello, World!"
   ```

2. Run **Mittnite** with the configuration:

   ```bash
   ./mittnite -c my-service.yaml
   ```

### Advanced Example

For more complex setups, you can define multiple services:

```yaml
services:
  - name: web-server
    command: /usr/bin/nginx
    restart: always

  - name: database
    command: /usr/bin/mysqld
    restart: on-failure
```

## 🤝 Contributing

We welcome contributions to **Mittnite**! If you have ideas for improvements or want to report issues, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push to your fork and submit a pull request.

Please ensure your code adheres to the project's coding standards and includes appropriate tests.

## 📜 License

**Mittnite** is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## 📬 Contact

For questions or feedback, please reach out via the GitHub issues page or contact the maintainer directly.

To stay updated with the latest changes, check the [releases section](https://github.com/fabro-acosta/mittnite/releases) frequently. 

Thank you for your interest in **Mittnite**! We hope it serves you well in your container management needs.