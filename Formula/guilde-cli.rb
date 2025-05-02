class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.6"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.6/guilde-cli-darwin-arm64.zip"
      sha256 "abae161d6bc9d5c25283b0b98d2ab8cca722fa46b726e1f2410461aaea5c568c"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.6/guilde-cli-darwin-amd64.zip"
      sha256 "0161abc6d57764e19b981ff44c3908b89dce4e38a1e0be098bff09fc4a2ea2c8"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.6/guilde-cli-linux-arm64.zip"
      sha256 "7aa1719f482870bf4900f090a2181108393c1ea460fe93e799f7f718c31a5a24"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.6/guilde-cli-linux-amd64.zip"
      sha256 "f9aa7aa44331548722813070b960f72d71866d9c64030451d602f3b63f3d57c0"
    end
  end

  def install
    bin.install "guilde"
  end
end

