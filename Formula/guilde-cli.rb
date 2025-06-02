class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.52.0"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.52.0/guilde-cli-darwin-arm64.zip"
      sha256 "be486e54ba53c145d443b698a3ce131cddc8b213139db524e71b9bdfa42be292"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.52.0/guilde-cli-darwin-amd64.zip"
      sha256 "e26e085e2ed2abc6a84021df49df5255ce8c5a98db8f7247d2693bbf23923c0d"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.52.0/guilde-cli-linux-arm64.zip"
      sha256 "f749198ed696174f2932ae333f4d2990424e4f578d2b552d6c3aa9b652f2eb4b"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.52.0/guilde-cli-linux-amd64.zip"
      sha256 "77dfd4b5e102c5ea12349102ddd9cba02b19603eee7f91ac4f1ff34660229974"
    end
  end

  def install
    bin.install "guilde"
  end
end

