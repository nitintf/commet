# Commet 🚀

**AI-powered commit message generator and Git utilities to help you write better commits.**

---

## ✨ Features

- 🤖 **AI Commit Messages**: Instantly generate meaningful commit messages using advanced AI models
- 🧠 **Multiple AI Providers**: Supports OpenAI, Claude, Google Gemini, and Groq
- 🛠️ **Interactive TUI**: Easy-to-use terminal UI for setup and configuration
- 🔄 **Seamless Git Integration**: Works with your existing Git workflow, automates best practices
- ⚡ **Flexible Usage**: Command-line flags and interactive modes for every workflow

---

## 📝 Planned Features / TODOs

- **Project-specific settings:**
  - Allow defining extra rules and configuration per project folder
- **Local models:**
  - Support running with local LLMs (offline or self-hosted)
- **More Git utilities:**
  - Add additional helpful git-related commands and automations

---

## 🏄‍♂️ Installation

### Using Homebrew (Recommended)

```sh
brew tap nitintf/homebrew-commet
brew install commet
```

### From Source (Go 1.21+ required)

```sh
git clone https://github.com/nitintf/commet.git
cd commet
make install
```

---

## 🚦 Quick Start

1. **Configure your AI provider:**
   ```sh
   commet config set
   ```
2. **Generate a commit message:**
   ```sh
   commet commit
   ```
3. **Enjoy smarter, faster commits!**

---

## 💡 Why Commet?

- Save time and mental energy on writing commit messages
- Enforce consistent, high-quality commit history
- Integrate seamlessly with your favorite AI providers
- Supercharge your Git workflow with automation and best practices

---

## 📚 More

- [Documentation](https://github.com/nitintf/commet)
- [Issues & Feedback](https://github.com/nitintf/commet/issues)

---

> Made with ❤️ by [Nitin Panwar](https://github.com/nitintf)

---

## 🛠️ Troubleshooting (macOS)

If you see a security error after installing with Homebrew (e.g., "commet cannot be opened because the developer cannot be verified"), run:

```sh
xattr -d com.apple.quarantine /opt/homebrew/bin/commet
```

This removes the Apple quarantine attribute from the binary.

---
