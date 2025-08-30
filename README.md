# Response Meter

A lightweight tool to continuously hit a target server, collect responses, and calculate traffic distribution percentages. It's useful for A/B testing, canary releases, or any situation where you need to verify how traffic is being routed between multiple systems.

## Features

- 🚀 **Concurrent Load Testing**: Send multiple concurrent requests to measure traffic distribution
- 📊 **Real-time Statistics**: Live TUI dashboard showing response distribution and request rates (updates every 2 seconds)
- 🎯 **HTTP Status Code Tracking**: Monitor and visualize HTTP status code distribution
- ⚡ **High Performance**: Built with Go for efficient concurrent request handling
- 🔧 **Extensible Architecture**: Designed to support various collectors and reporters in the future
- 📦 **Zero Dependencies**: Built using only Go's standard library - no external dependencies required

## Installation

Requires Go 1.25+

```bash
go install github.com/phhphc/response-meter@latest
```

## Usage

### Basic Usage

```bash
response-meter -t="https://httpbin.org/status/200%2C404%2C500"
```

### Command Line Options

- `-t string`: HTTP(S) URL to probe. (required)
- `-i duration`: Interval between reports (e.g., 500ms, 2s, 1m) (default: 2s)
- `-c int`: Number of concurrent requests (default: 1)
- `-d duration`: Per-request timeout; 0 disables timeout (default: 0s)

### Sample Output

The tool provides a real-time TUI dashboard displaying:

```
                                 Response Meter

----------------------------------- ALL TIME -----------------------------------

Duration: 54s                  Requests/sec: 1.4        Total Requests: 75
Response Distribution:
  Timeout                           4.0% [█                             ] 3
  200                              30.7% [█████████                     ] 23
  500                              33.3% [█████████                     ] 25
  404                              32.0% [█████████                     ] 24

--------------------------------- LAST PERIOD ----------------------------------

Duration: 3s                   Requests/sec: 2.3        Total Requests: 7
Response Distribution:
  Timeout                           0.0% [                              ] 0
  200                              42.9% [████████████                  ] 3
  500                              42.9% [████████████                  ] 3
  404                              14.3% [████                          ] 1

Press Ctrl+C to exit.
```

## Use Cases

- Monitor traffic distribution between different versions of your application to ensure proper routing ratios.
- Verify that canary deployments are receiving the expected percentage of traffic and responding correctly.
- Test load balancer behavior and ensure traffic is being distributed as expected across backend services.

## Architecture

The project follows a modular design with three main components:

- **Meter**: Orchestrates concurrent requests and aggregates statistics
- **Collectors**: Gather data (currently supports HTTP status codes)
- **Reporters**: Display results (currently supports TUI output)

### Project Structure

```
response-meter/
├── main.go                              # CLI entry point
├── internal/                            # Internal packages
│   ├── meter/                           # Core measurement logic
│   ├── collector/                       # Data collection implementations
│   └── reporter/                        # Output formatting
└── pkg/                                 # Reusable packages
```

## Future Roadmap

The project is designed to be extensible. Planned enhancements include:

### Additional Collectors

- **Request Headers**: Collect and analyze custom response headers
- **Response Body**: Parse and categorize response body content
- **Response Time**: Measure and track request latency

### Additional Reporters

- **JSON Output**: Machine-readable output format

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development

1. Clone the repository
2. Make your changes
3. Test your changes
4. Submit a pull request

## License

This project is released into the public domain under [The Unlicense](LICENSE).

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/phhphc/response-meter/issues) on GitHub.
