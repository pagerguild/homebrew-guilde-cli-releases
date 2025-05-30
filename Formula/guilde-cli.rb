class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.51.0"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.51.0/guilde-cli-darwin-arm64.zip"
      sha256 "7aa9ab0cd703c7c1d8d558dcca757c4d370435d09a817d206756082ab3419e00"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.51.0/guilde-cli-darwin-amd64.zip"
      sha256 "c66aaf785d5f639f4990d06614ff5cab870f08c0ca3e12760a95b05a5d6b94d6"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.51.0/guilde-cli-linux-arm64.zip"
      sha256 "686e1be89c394af3ff7e9b4e0be565362456ed96bc4b4126f581073100da0b54"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.51.0/guilde-cli-linux-amd64.zip"
      sha256 "04f59e149d893c4a4cbd0ac1baccbccdfa0a4755876ffce5a6f79143e85d4f88"
    end
  end

  def install
    bin.install "guilde"
  end
end

