# Commet ☄️

AI-powered commit message generator and Git utilities to help you write better commits.

## Features

- 🤖 **Smart Commit Messages**: Generate meaningful commit messages using AI
- 🔧 **Multiple AI Providers**: Support for OpenAI, Claude, Google Gemini, and Groq
- ⚙️ **Interactive Configuration**: Easy-to-use TUI for setup and configuration
- 🌐 **Git Integration**: Configurable git workflows and automation
- 📝 **Flexible**: Command-line flags and interactive modes

## Installation

### From Source

```bash
git clone https://github.com/yourusername/commet.git
cd commet
make build
```

The binary will be available at `./bin/commet`.

### Install to System

```bash
make install
```

This will install `commet` to your `$GOPATH/bin` directory.

## Quick Start

1. **Configure Commet**: Set up your AI provider and API key
   ```bash
   commet config set
   ```

2. **Generate commit messages**: (Feature coming soon)
   ```bash
   commet commit
   ```

3. **View current configuration**:
   ```bash
   commet config show
   ```

## Configuration

### Interactive Configuration (TUI)

Run the interactive configuration interface:

```bash
commet config set
```

This will open a terminal user interface where you can:
- Select AI provider (OpenAI, Claude, Google, or Groq)
- Set your API key (supports clipboard paste with Ctrl+V/Cmd+V)
- Configure AI model from predefined options
- Adjust Git settings

### Command Line Configuration

You can also configure settings using command-line flags:

```bash
# Set AI provider
commet config set --provider openai

# Set API key
commet config set --api-key "your-api-key-here"

# Set model
commet config set --model "gpt-4"
```

### Keyboard Shortcuts (TUI)

When editing text fields (like API keys):
- **Ctrl+V / Cmd+V**: Paste from clipboard
- **Ctrl+A / Cmd+A**: Clear current input
- **Ctrl+U**: Clear entire field
- **Backspace**: Delete one character
- **Enter**: Save and continue
- **Esc**: Cancel and go back

### Configuration Options

#### AI Settings
- **Provider**: Choose from `openai`, `claude`, `google`, or `groq`
- **API Key**: Your AI provider's API key
- **Model**: Select from available models for your chosen provider

#### Git Settings
- **Auto Stage**: Automatically stage changes before committing
- **Show Diff**: Display diff before generating commit message
- **Confirm Push**: Ask for confirmation before pushing

### Configuration File

Settings are stored in `~/.commet.yaml`:

```yaml
ai:
  provider: openai
  api_key: your-api-key
  model: gpt-4
git:
  auto_stage: false
  show_diff: true
  confirm_push: true
```

## API Keys

You'll need an API key from one of the supported providers:

- **OpenAI**: Get your API key from [OpenAI Platform](https://platform.openai.com/api-keys)
- **Claude**: Get your API key from [Anthropic Console](https://console.anthropic.com/)
- **Google**: Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
- **Groq**: Get your API key from [Groq Console](https://console.groq.com/keys)

## Usage

### Available Commands

```bash
commet                 # Show help
commet config          # Manage configuration
commet config show     # Display current settings
commet config set      # Interactive configuration
commet commit          # Generate commit messages (coming soon)
commet --help          # Show detailed help
```

### Development Commands

Use the included Makefile for development:

```bash
make build    # Build the binary
make clean    # Clean build artifacts
make install  # Install to system
make test     # Run tests
make lint     # Run linter
make fmt      # Format code
make run      # Build and run
make dev      # Development mode
```

## Examples

### Setting up OpenAI

```bash
commet config set --provider openai --api-key "sk-..." --model "gpt-4"
```

### Setting up Claude

```bash
commet config set --provider claude --api-key "sk-ant-..." --model "claude-3-sonnet-20240229"
```

### Interactive Setup

```bash
commet config set
# Follow the interactive prompts to configure all settings
```

## Project Structure

```
commet/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command and setup
│   └── config.go       # Configuration commands
├── internal/
│   ├── config/         # Configuration management
│   │   ├── config.go   # Main config logic
│   │   ├── ai.go       # AI provider config
│   │   └── git.go      # Git settings config
│   └── ui/             # Terminal user interface
│       └── config.go   # TUI for configuration
├── Makefile            # Build and development commands
├── go.mod              # Go module definition
├── go.sum              # Go module dependencies
└── README.md           # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make fmt`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development

### Prerequisites

- Go 1.19 or later
- Make (optional, for using Makefile commands)

### Building

```bash
# Build for development
make build

# Build and run
make run ARGS="config show"

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run linter
make lint

# Format code
make fmt
```

## Roadmap

- [x] Configuration management (TUI and CLI)
- [x] Multiple AI provider support
- [x] Git settings integration
- [ ] Commit message generation
- [ ] Git hooks integration
- [ ] Commit message templates
- [ ] Branch name integration
- [ ] Multi-language support

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/yourusername/commet/issues) page
2. Create a new issue with detailed information
3. Include your configuration and error messages

---

**Made with ☄️ by the Commet team**