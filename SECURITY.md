# Security Policy

## Supported Versions

We actively support the following versions of ViberCode CLI with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of ViberCode CLI seriously. If you discover a security vulnerability, please follow these steps:

### 🚨 For Security Issues

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please:

1. **Email us directly**: security@vibercode.com
2. **Include the following information**:
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact assessment
   - Any suggested fixes (if you have them)

### 📋 What to Include

When reporting a security vulnerability, please provide:

- **Vulnerability Type**: What kind of vulnerability is it? (e.g., code injection, authentication bypass, etc.)
- **Location**: Where in the codebase is the vulnerability?
- **Impact**: What could an attacker accomplish by exploiting this vulnerability?
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Environment**: OS, Go version, ViberCode CLI version where you discovered the issue

### ⏱️ Response Timeline

- **Acknowledgment**: We'll acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We'll provide an initial assessment within 5 business days
- **Status Updates**: We'll keep you informed of our progress every 10 business days
- **Resolution**: We aim to resolve critical security issues within 30 days

### 🔄 Our Process

1. **Confirmation**: We'll work to confirm and understand the vulnerability
2. **Assessment**: We'll assess the severity and impact
3. **Fix Development**: We'll develop and test a fix
4. **Coordinated Disclosure**: We'll coordinate the release timing with you
5. **Public Disclosure**: After the fix is released, we'll publicly disclose the issue

### 🏆 Security Acknowledgments

We believe in recognizing security researchers who help improve our security:

- **Hall of Fame**: We maintain a security researchers hall of fame
- **Credit**: With your permission, we'll credit you in our security advisories
- **Swag**: Security researchers who report valid vulnerabilities receive ViberCode swag

### 🛡️ Security Best Practices

When using ViberCode CLI:

#### For Generated APIs

- **Environment Variables**: Always use environment variables for sensitive configuration
- **Database Credentials**: Never commit database credentials to version control
- **API Keys**: Rotate API keys regularly
- **HTTPS**: Always use HTTPS in production
- **Input Validation**: Validate all user inputs
- **Authentication**: Implement proper authentication and authorization

#### For AI Features

- **API Keys**: Keep your Anthropic API key secure
- **Chat Logs**: Be aware that chat logs may contain sensitive information
- **Code Review**: Review AI-generated code before deployment

#### For Development

- **Dependencies**: Keep dependencies up to date
- **Build Environment**: Use secure build environments
- **Code Review**: Conduct security-focused code reviews

### 🔗 Related Resources

- **OWASP Top 10**: [https://owasp.org/www-project-top-ten/](https://owasp.org/www-project-top-ten/)
- **Go Security Checklist**: [https://github.com/securego/gosec](https://github.com/securego/gosec)
- **Node.js Security**: [https://nodejs.org/en/security/](https://nodejs.org/en/security/)

### 📞 Contact Information

- **Security Email**: security@vibercode.com
- **General Contact**: team@vibercode.com
- **PGP Key**: Available upon request

---

Thank you for helping us keep ViberCode CLI secure! 🔒
