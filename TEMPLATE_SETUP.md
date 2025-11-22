# GitHub Template Repository Setup Guide

This guide explains how to enable the **"Use this template"** button on GitHub.

## Enabling Template Repository

### Step 1: Push to GitHub

First, push this repository to GitHub:

```bash
# Initialize git (if not already done)
git init
git add .
git commit -m "Initial commit: Go Clean Architecture Template"

# Create repository on GitHub, then:
git remote add origin https://github.com/yourusername/go-scaffolding.git
git branch -M main
git push -u origin main
```

### Step 2: Enable Template Repository Feature

1. Go to your repository on GitHub: `https://github.com/yourusername/go-scaffolding`

2. Click on **Settings** (top right)

3. In the **General** section (left sidebar), scroll down to the **Template repository** section

4. Check the box that says **"Template repository"**

   ![Template Repository Checkbox](https://docs.github.com/assets/cb-25570/images/help/repository/template-repository-checkbox.png)

5. Click **Save** (if needed)

### Step 3: Verify Template is Enabled

Once enabled, you'll see a **"Use this template"** button appear on the repository home page, next to the "Code" button.

## Using the Template

### For Template Users

When someone wants to use your template:

1. Go to `https://github.com/yourusername/go-scaffolding`
2. Click the green **"Use this template"** button
3. Select **"Create a new repository"**
4. Fill in:
   - Repository name
   - Description
   - Public/Private
5. Click **"Create repository"**

### What Gets Copied

When creating from template:
- ✅ All files and directories
- ✅ Branch structure
- ✅ README and documentation
- ✅ Configuration files
- ❌ Git history (starts fresh)
- ❌ Issues and Pull Requests
- ❌ GitHub Actions workflow runs
- ❌ Excluded files (see `.github/template.yml`)

## After Creating from Template

Users should follow these steps after creating a repository from the template:

### 1. Update Module Name

```bash
# Replace the module name
find . -type f -name '*.go' -exec sed -i 's|github.com/yourusername/go-scaffolding|github.com/yourorg/yourproject|g' {} +

# Update go.mod
go mod edit -module github.com/yourorg/yourproject
go mod tidy
```

### 2. Update Configuration

Edit `config.yaml`:
- Change `app.name` to your project name
- Update database credentials
- Configure your environment

### 3. Update Documentation

- Edit `README.md` with your project details
- Update repository URLs in documentation
- Customize the architecture for your needs

### 4. Initialize Infrastructure

```bash
# Start infrastructure
task docker:up

# Run migrations
task migrate:up

# Build and run
task build:api
task run:api
```

### 5. Verify Everything Works

```bash
# Run tests
task test
task test:integration

# Check API
curl http://localhost:8080/health/live
```

## Template Maintenance

### Keeping Template Updated

As the template maintainer:

```bash
# Update dependencies
go get -u ./...
go mod tidy

# Update Go version
# Edit go.mod and Dockerfile

# Test everything
task test:all

# Commit and push
git add .
git commit -m "chore: update dependencies"
git push
```

### Version Tags

Consider tagging releases:

```bash
git tag -a v1.0.0 -m "Release v1.0.0: Initial stable release"
git push origin v1.0.0
```

Users can then choose specific versions when creating from template.

## Template Features

This template includes:

✅ **Hexagonal Architecture** - Production-ready structure
✅ **Multiple Protocols** - REST API ready, gRPC/CLI planned
✅ **Database Support** - PostgreSQL with GORM, migrations
✅ **Testing** - Unit, integration tests (83.8% coverage)
✅ **Developer Tools** - Taskfile, Docker, hot reload
✅ **Documentation** - Comprehensive guides and ADRs
✅ **CI/CD Ready** - GitHub Actions workflows
✅ **Docker** - Multi-stage Dockerfile with scratch base

## Support

- 📚 [Documentation](./README.md)
- 🏗️ [Architecture Guide](./docs/architecture/overview.md)
- 📝 [ADRs](./docs/adr/)
- 💬 [GitHub Discussions](https://github.com/yourusername/go-scaffolding/discussions)
- 🐛 [Report Issues](https://github.com/yourusername/go-scaffolding/issues)

## License

This template is released under the MIT License. See [LICENSE](./LICENSE) file.

---

**Note**: Replace `yourusername/go-scaffolding` with your actual GitHub username/repository throughout this guide.
