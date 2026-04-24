# Claudes - Claude CLI Configuration Manager

A simple command-line tool to manage multiple Claude configurations using environment files. Switch between different API keys, models, and endpoints with a single command.

## Features

- **Profile Management**: Create, update, delete, and list multiple profiles (e.g., `kimi`, `deepseek`, `gptoss`)
- **Interactive Commands**: Add and update profiles interactively with prompts
- **Environment Variables**: Load different API keys, base URLs, and models per profile
- **Simple Commands**: List profiles, run Claude with a specific configuration, manage profiles
- **Auto-created Config Directory**: `~/.config/claudes/` is created automatically if it doesn't exist

## Installation

### Quick Install / Update (Recommended)

Run the following command to **install or update** claudes automatically:

```bash
curl -fsSL https://raw.githubusercontent.com/yasaricli/claudes/refs/heads/develop/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/yasaricli/claudes/refs/heads/develop/install.sh | bash
```

The script will:
- Check if Go 1.16+ is installed
- **Auto-detect** if claudes is already installed (install vs update mode)
- Clone the repository (develop branch)
- Build the binary
- Install/update it to `~/.local/bin` or `/usr/local/bin`

### Check Version

To see your current version:

```bash
claudes version
# or
claudes --version
```

### Manual Installation

1. **Install Go**: Make sure you have Go 1.16+ installed
2. **Clone or download the repository**
3. **Build and install**:

```bash
cd claudes
go build -o claudes .
sudo mv claudes /usr/local/bin/
```

Alternatively, you can install it directly:

```bash
go install claudes
```

Make sure your `$GOPATH/bin` is in your `$PATH`.

## Configuration

### Directory Structure

Claudes looks for `.env` files in `~/.config/claudes/` directory. Each `.env` file represents a profile.

```
~/.config/claudes/
├── anthropic.env
├── openai.env
├── deepseek.env
└── custom.env
```

### Profile Format

Each `.env` file should contain Claude CLI environment variables:

```env
# Example: anthropic.env (default Claude API)
ANTHROPIC_AUTH_TOKEN=sk-ant-xxx
ANTHROPIC_BASE_URL=https://api.anthropic.com/v1
ANTHROPIC_MODEL=claude-3-opus-20240229
```

```env
# Example: openai.env
ANTHROPIC_AUTH_TOKEN=sk-xxx
ANTHROPIC_BASE_URL=https://api.openai.com/v1
ANTHROPIC_MODEL=gpt-4
```

```env
# Example: deepseek.env
ANTHROPIC_AUTH_TOKEN=your-deepseek-api-token
ANTHROPIC_BASE_URL=https://api.deepseek.com/v1
ANTHROPIC_MODEL=deepseek-chat
```

### Required Variables

At minimum, your `.env` file should include:

- `ANTHROPIC_AUTH_TOKEN`: Your API key/token
- `ANTHROPIC_BASE_URL`: The API endpoint URL
- `ANTHROPIC_MODEL`: The model to use

## Usage

### Add a New Profile

```bash
claudes add
# or
claudes add <profile-name>
```

Interactive example:
```
Enter profile name: custom
Enter ANTHROPIC_AUTH_TOKEN: your-api-token
Enter ANTHROPIC_BASE_URL (e.g., https://api.anthropic.com/v1): https://api.example.com/v1
Enter ANTHROPIC_MODEL (e.g., claude-3-opus-20240229): custom-model
Profile 'custom' created successfully!
```

### List Profiles

```bash
claudes list
```

Output example:
```
Available profiles:
  - anthropic
  - openai
  - deepseek
  - custom
```

### Update an Existing Profile

```bash
claudes update
# or
claudes update <profile-name>
```

Interactive example:
```
Updating profile 'anthropic' (leave empty to keep current value)

Enter ANTHROPIC_AUTH_TOKEN (current: sk-a********1234):
Enter ANTHROPIC_BASE_URL (current: https://api.anthropic.com/v1):
Enter ANTHROPIC_MODEL (current: claude-3-opus-20240229):

Profile 'anthropic' updated successfully!
```

> 💡 All existing values are shown in parentheses. Just press Enter to keep the current value, or enter a new value to change it. 
Enter new ANTHROPIC_BASE_URL: https://api.newendpoint.com/v1
Enter new ANTHROPIC_MODEL: 
Profile 'kimi' updated successfully!
```

### Delete a Profile

```bash
claudes delete
# or
claudes delete <profile-name>
```

Example:
```
Are you sure you want to delete profile 'kimi'? (y/N): y
Profile 'kimi' deleted successfully!
```

### Run Claude with a Profile

```bash
claudes <profile> [claude-args]
```

Examples:

```bash
# Run Claude with kimi profile
claudes kimi

# Run Claude with deepseek profile and extra arguments
claudes deepseek --help
```

### Help

```bash
claudes help
```

Output:
```
A CLI tool to manage multiple Claude configurations using environment files.

Usage:
  claudes [flags]
  claudes [command]

Available Commands:
  add         Add a new profile interactively
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an existing profile
  help        Help about any command
  list        List all available profiles
  update      Update an existing profile interactively

Flags:
  -h, --help   help for claudes

Use "claudes [command] --help" for more information about a command.
```

## How It Works

1. When you run `claudes <profile>`:
   - Claudes looks for `<profile>.env` in `~/.config/claudes/`
   - Loads all environment variables from that file
   - Finds the `claude` executable in your `PATH`
   - Runs the Claude CLI with the loaded environment variables

2. The Claude CLI will then use the `ANTHROPIC_*` variables to connect to the specified API.

## Prerequisites

- **Claude CLI**: Make sure you have the official Claude CLI installed and in your `PATH`
- **Go**: For building from source (Go 1.16+)

## Dependencies

This project uses:
- [Cobra](https://github.com/spf13/cobra): For CLI command handling
- [godotenv](https://github.com/joho/godotenv): For loading `.env` files

## License

MIT

## Contributing

Feel free to open issues or submit pull requests.
