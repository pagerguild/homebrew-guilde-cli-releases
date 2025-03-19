class GuildeCli < Formula
  desc "Guilde Command Line Interface"
  homepage "https://github.com/pagerguild/homebrew-guilde-cli-releases"
  version "0.41.11"

  on_macos do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.11-darwin-arm64.zip"
      sha256 "4bb0459ef619cf576aacad12e4bb03826518f48f90dd46a870d97f4704c4d1ed"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.11-darwin-amd64.zip"
      sha256 "12422c73c895e5fec39885140ce046d630c1546295055f5afe3804f0ee6f4d75"
    end
  end

  on_linux do
    on_arm do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.11-linux-arm64.zip"
      sha256 "359d2f4b59db75cd486b256a55fe79d8f0d867da2138ce8baa5ceb8b927e0080"
    end
    on_intel do
      url "https://guilde-cli-releases.vercel.app/guilde-cli-0.41.11-linux-amd64.zip"
      sha256 "4b2f2a831997351bf663afe3b427cad852e7f80fc67a9d04b51d91f56237f60a"
    end
  end

  def install
    bin.install "guilde"
  end
end
