class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.19"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.19/guilde-cli-darwin-arm64.zip"
      sha256 "780f7c8b5c265982380c6f15128f68aa983462e4fab251a983b1df140a618fb9"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.19/guilde-cli-darwin-amd64.zip"
      sha256 "d9bcca509e8d2cc6a3da199bc17d321d3bacff6a4ddc10bc0fdd57dc6d3375f9"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.19/guilde-cli-linux-arm64.zip"
      sha256 "ead63eaa3d6b693ead81689d61033999d2d49d35290cd485ab721ed5ccdeac02"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.19/guilde-cli-linux-amd64.zip"
      sha256 "a89e7038213f2bff8264e3766fac32ed4b30e9b7b7783aa05bd2cfc22b0d892e"
    end
  end

  def install
    bin.install "guilde"
  end
end

