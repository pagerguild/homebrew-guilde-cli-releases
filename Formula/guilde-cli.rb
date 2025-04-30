class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.44.2"

  on_macos do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.2/guilde-cli-darwin-arm64.zip"
      sha256 "3137b4487d68305758bc83540f82cca9453b3137f6d0fe16cef6a9612127270f"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.2/guilde-cli-darwin-amd64.zip"
      sha256 "553ced590f94ae752b57cab69674e3847167b46c1a51c2bb3dcc285bbb6606e9"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.2/guilde-cli-linux-arm64.zip"
      sha256 "4debdec74ff4e8869b720ca9010e31fbfda913b1ee7ce0da638d6bff3985ff9b"
    end
    on_intel do
      url "https://github.com/pagerguild/homebrew-guilde-cli-releases/releases/download/v0.44.2/guilde-cli-linux-amd64.zip"
      sha256 "54e6417aa4b9a6bbe6b628a1d24c8a76999e8ebe7d2feef294416e9fde958bb9"
    end
  end

  def install
    bin.install "guilde"
  end
end

