class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.7"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.7/guilde-cli-darwin-arm64.zip"
      sha256 "d56629f53d5140a5904f80ab6bb3eaebaa845dc6f556fb17dcf7846acad855d4"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.7/guilde-cli-darwin-amd64.zip"
      sha256 "a9b9b8e6cec73ec05881aac11c1038a357c7f521a958eb96aba128384bf3c937"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.7/guilde-cli-linux-arm64.zip"
      sha256 "a8b14becef9423af1d8ac55cdd4a5dbb4be940cf93a01ec53eb3a4d7c2d89ae0"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.41.7/guilde-cli-linux-amd64.zip"
      sha256 "8ee96efb019f8c4b7bb44329fc3501649979322c48235eb900f0fe4931fc9e23"
    end
  end

  def install
    bin.install "bin/guilde-cli"
  end
end
