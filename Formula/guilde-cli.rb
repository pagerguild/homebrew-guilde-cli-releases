class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.43.0"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.0/guilde-cli-darwin-arm64.zip"
      sha256 "af69cf260bd5345643d082978890bd285c02ff9928dcf682026500b6c0ee9660"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.0/guilde-cli-darwin-amd64.zip"
      sha256 "e7cbc804ab9dc41513c08cbaa974bf899fe7d12ee7c55d85932e63d38ba00fbc"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.0/guilde-cli-linux-arm64.zip"
      sha256 "2ca90fbe5aae381aa897fac524e25d883523f01fd6d2e5c937cf9ccf89600c24"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.0/guilde-cli-linux-amd64.zip"
      sha256 "063ca450d0c6942b6c3e267a27659d2322c6980bb7e755c0ec3138f8669f2107"
    end
  end

  def install
    bin.install "guilde"
  end
end

