class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.1"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.1/guilde-cli-darwin-arm64.zip"
      sha256 "45d3984f80600a123303d2e3559ba56c741f918612cae6d832e7f7ae39ef8f31"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.1/guilde-cli-darwin-amd64.zip"
      sha256 "1e46d22e57680434bce31acebb62294d070329d0e9b5306012586ace4b968603"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.1/guilde-cli-linux-arm64.zip"
      sha256 "9f3c84651f349bba7f72bbe6bbf5b56599c3d09ad3aa5444185d0bfe2ebc8688"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.1/guilde-cli-linux-amd64.zip"
      sha256 "0b1d522524c8b3389ee6c720d0cda3fbe86e27aebb437afda04454f63feaeaf9"
    end
  end

  def install
    bin.install "guilde"
  end
end

