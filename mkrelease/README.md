# mkrelease - Guilde CLI Release Tool

This tool automates the release process for the Guilde CLI, handling both GitHub releases and Homebrew formula generation.

## Overview

The `mkrelease` tool performs the following tasks:

1. Validates release assets and release notes
2. Computes SHA256 checksums for all binary assets
3. Generates a Homebrew formula file with the correct URLs and checksums
4. Commits the formula to the repository
5. Creates and pushes a Git tag for the release
6. Creates a GitHub release
7. Uploads all assets to the GitHub release

## Homebrew Formula Generation

The tool generates a Homebrew formula file (`guilde-cli.rb`) that follows current Homebrew conventions for binary-only distributions.

### Formula File Naming

In Homebrew:

- The formula filename must be lowercase (`guilde-cli.rb`)
- The formula class name is automatically derived from the filename in CamelCase (`GuildeCli`)
- Filenames should match the project name as closely as possible
- Formula files should be named after the software they install

For more details on Homebrew naming conventions, see the [Formula Cookbook](https://docs.brew.sh/Formula-Cookbook#a-quick-word-on-naming).

### Formula Structure

The generated formula uses the recommended structure for cross-platform binary distributions:

```ruby
class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/guilde-cli-releases"
  version "1.2.3"  # Example version

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v1.2.3/guilde-cli-darwin-arm64.zip"
      sha256 "..."
    end
    on_intel do
      url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v1.2.3/guilde-cli-darwin-amd64.zip"
      sha256 "..."
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v1.2.3/guilde-cli-linux-arm64.zip"
      sha256 "..."
    end
    on_intel do
      url "https://github.com/pagerguild/guilde-cli-releases/releases/download/v1.2.3/guilde-cli-linux-amd64.zip"
      sha256 "..."
    end
  end

  def install
    bin.install "guilde"
  end
end
```

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook) - Official documentation on how to write Homebrew formulas
- [Formula API Documentation](https://rubydoc.brew.sh/Formula) - Ruby API documentation for the Formula class
- [Handling Different System Configurations](https://docs.brew.sh/Formula-Cookbook#handling-different-system-configurations) - Documentation on using `on_macos`, `on_arm`, and `on_intel` conditions

## Usage

```
mkrelease RELEASE_VERSION DIRECTORY_PATH REPO_PATH
```

- `RELEASE_VERSION`: The version number for the release (e.g., "1.2.3")
- `DIRECTORY_PATH`: Path to the directory containing release assets
- `REPO_PATH`: Path to the repository

The tool also requires the `GITHUB_PAT` environment variable to be set with a valid GitHub Personal Access Token. 