class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.7"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.7/guilde-cli-darwin-arm64.zip"
      sha256 "34fef46a8dfaba40e044dd4e4b3d10c7d3390977f314119f8fbfe6077eeebb40"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.7/guilde-cli-darwin-amd64.zip"
      sha256 "d0b530131271a9207caa676056d8480ef8bb8cd6f27cd0c0322c94a48ee9c8b7"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.7/guilde-cli-linux-arm64.zip"
      sha256 "607df248dbef688495da61985458d795f9c380101f8b0efa21d61c2bd0474968"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.7/guilde-cli-linux-amd64.zip"
      sha256 "0409e509f8775e747eb42989c5a9ee8c6ba4b3fe14fcbea24ef65ba2241524ad"
    end
  end

  def install
    bin.install "guilde"
  end
end

