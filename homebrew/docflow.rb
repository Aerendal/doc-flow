class Docflow < Formula
  desc "Documentation governance CLI"
  homepage "https://github.com/jerzy/docflow"
  version "0.1.0"
  if Hardware::CPU.arm?
    url "https://github.com/jerzy/docflow/releases/download/v0.1.0/docflow-darwin-arm64.tar.gz"
    sha256 "SKIP"
  else
    url "https://github.com/jerzy/docflow/releases/download/v0.1.0/docflow-darwin-amd64.tar.gz"
    sha256 "SKIP"
  end

  def install
    bin.install Dir["docflow*"].first => "docflow"
  end

  test do
    system "#{bin}/docflow", "--help"
  end
end
