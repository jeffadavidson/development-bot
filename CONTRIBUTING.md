# Contributing to Development Bot

Thank you for your interest in contributing to the Development Bot! This project helps communities stay informed about local development activity through automated RSS feeds.

## 🎯 Ways to Contribute

### 🐛 Bug Reports
- Use GitHub Issues to report bugs
- Include steps to reproduce the issue
- Provide your operating system and Go version
- Include relevant log output or error messages

### 💡 Feature Requests
- Use GitHub Issues with the "enhancement" label
- Describe the use case and expected behavior
- Consider how it might benefit other communities

### 🔧 Code Contributions
- Fork the repository
- Create a feature branch (`git checkout -b feature/amazing-feature`)
- Make your changes
- Add tests for new functionality
- Ensure all tests pass (`go test ./...`)
- Commit with clear messages
- Push to your branch
- Create a Pull Request

## 🏗️ Development Setup

### Prerequisites
- Go 1.20 or later
- Git

### Local Development
```bash
# Clone your fork
git clone https://github.com/yourusername/development-bot.git
cd development-bot

# Install dependencies
go mod download

# Run tests
go test ./...

# Run the application
go run main.go
```

## 📝 Code Style

### Go Standards
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

### Testing
- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Follow existing test patterns

## 🌍 Adapting for Other Cities

This bot is designed to be adaptable for any city with open data APIs:

### Configuration Changes
1. Update `config.yaml` with your neighborhood boundaries
2. Modify API endpoints in `interactions/calgaryopendata/`
3. Adjust data parsing for your city's JSON structure
4. Update permit status mappings if needed

### Common Adaptations
- **API Endpoints**: Change base URLs for your city's open data
- **Field Mappings**: Adjust JSON field names to match your data
- **Status Values**: Update permit status strings
- **Geographic Bounds**: Set your neighborhood's lat/lng coordinates

## 🚀 Deployment

### GitHub Pages
- Fork the repository
- Enable GitHub Actions in your fork
- Set up GitHub Pages in repository settings
- The bot will run daily and publish RSS feeds automatically

### Custom Deployment
- Adapt the GitHub Actions workflow for your CI/CD platform
- Configure RSS feed hosting (any web server works)
- Set up scheduled execution (cron jobs, etc.)

## 📋 Pull Request Guidelines

### Before Submitting
- [ ] Tests pass locally
- [ ] Code follows Go standards
- [ ] Changes are documented
- [ ] Feature works end-to-end

### PR Description
- Describe what the change does
- Explain why the change is needed
- Include any breaking changes
- Reference related issues

## 🤝 Community Guidelines

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers get started
- Share knowledge and resources

## 📞 Getting Help

- **Questions**: Use GitHub Discussions
- **Bug Reports**: Use GitHub Issues
- **Feature Ideas**: Use GitHub Issues with "enhancement" label

## 🎉 Recognition

Contributors will be recognized in:
- GitHub Contributors list
- Project documentation
- Release notes for significant contributions

Thank you for helping make local government data more accessible! 🏙️ 