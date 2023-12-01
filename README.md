# Ntụrụ CLI (IheNtụrụ)

Ntụrụ is a powerful CLI tool designed to simplify the generation of microservices for your next Go lang application. With Ntụrụ, you can quickly set up the foundational structure for your microservices architecture, allowing you to focus on building and scaling your services.

## Features

- **Microservice Boilerplate:** Generate a boilerplate for your microservices, including essential directory structure and configuration files.
- **Custom Templates:** Easily customize templates to suit your specific microservices requirements.
- **Efficient and Fast:** Ntụrụ is designed to be fast and efficient, making it easy to kickstart your microservices projects.

## Installation

### Via Go Get

Ensure you have Go installed on your machine. Run the following command to install nturu globally:

```bash
go get -u github.com/CeoFred/nturu
```

### Via Binary Download

Download the latest release binary for your platform from the [GitHub Releases](https://github.com/CeoFred/nturu/releases) page. Make it executable and move it to a directory in your PATH.

```bash
curl -LJO https://github.com/CeoFred/nturu/releases/download/v1.0.2/nturu
chmod +x nturu
mv nturu /usr/local/bin/nturu
```

## Usage

### Generate Microservice Boilerplate

```bash
nturu generate 
```

This command generates a boilerplate from the input you would enter.

### Customize Templates

You can customize templates by modifying the template files in the `templates` directory. Run the generation command with the `-t` flag to use your custom templates:

```bash
nturu generate  -t path/to/custom/templates 
```

For more detailed information, run:

```bash
nturu --help
```

## Contributing

Contributions are welcome! Please read the [Contributing Guidelines](CONTRIBUTING.md) for more information.

## License

This project is licensed under the [MIT License](LICENSE).

---