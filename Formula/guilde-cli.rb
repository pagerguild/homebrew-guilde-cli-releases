class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.43.3"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.3/guilde-cli-darwin-arm64.zip"
      sha256 "ad588c8bc291be59e7fe7506084146d5d1b78f016783876e4fba35217409573a"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.3/guilde-cli-darwin-amd64.zip"
      sha256 "c02a2406bdbac615b1528868e3e0a61e5db7a8905be58dd4670212067347a72c"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.3/guilde-cli-linux-arm64.zip"
      sha256 "450beb54ba21de42e5eca9979faf72b1a4ba5176bcb02d989c1d93800a3024a2"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.43.3/guilde-cli-linux-amd64.zip"
      sha256 "69e03455c35ff91bde59dcf34f681184ea5ef1bb4a58083f0ff7ddce96065e08"
    end
  end

  def install
    bin.install "guilde"
  end
end

