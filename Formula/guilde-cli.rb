class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.53.0"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.53.0/guilde-cli-darwin-arm64.zip"
      sha256 "c56e98f1e14dab025299f318fee21e74dc1ee50a26e95e95833999c1e32468be"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.53.0/guilde-cli-darwin-amd64.zip"
      sha256 "1069d82280a1e1a93f0d005332221b197f553d6df8f468ae7325a110c44a3da3"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.53.0/guilde-cli-linux-arm64.zip"
      sha256 "31c18414f1591aa18276aabd118a8710b9ecac3b781dd1fa035143d3479a911d"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.53.0/guilde-cli-linux-amd64.zip"
      sha256 "5dc0578d02c59756e6b65a32632d90a4f0eca0d13bce6ec071ffbf8bb641029f"
    end
  end

  def install
    bin.install "guilde"
  end
end

