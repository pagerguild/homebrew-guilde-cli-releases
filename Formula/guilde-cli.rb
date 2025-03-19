class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.9"

  on_macos do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-darwin-arm64.zip"
      sha256 "ba8abc40d6fd38ef1c318092464b35860e90850fe50b5ced776ffc78cd649a7b"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-darwin-amd64.zip"
      sha256 "3ad7545325ae405af68bcdffa43ef638709028826a5f0474f2ff300caa096698"
    end
  end

  on_linux do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-linux-arm64.zip"
      sha256 "01f986fe0509342f27d33b3d89faba32be333493e088b512f37bd68df9dde680"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-linux-amd64.zip"
      sha256 "f6275b61eb8978fac666a96e8b5edd408c90bcbec0bef89b2fc80d517131a81e"
    end
  end

  def install
    bin.install "guilde"
  end
end
