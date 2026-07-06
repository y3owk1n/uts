{ pkgs, version, usePrebuilt, commitHash ? null }:

let
  system = pkgs.stdenv.hostPlatform.system;
  # Pre-built binary URLs (replace with actual release URLs)
  prebuilt = {
    "aarch64-darwin" = {
      url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-darwin-arm64";
      sha256 = "";
    };
    "x86_64-darwin" = {
      url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-darwin-amd64";
      sha256 = "";
    };
    "aarch64-linux" = {
      url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-linux-arm64";
      sha256 = "";
    };
    "x86_64-linux" = {
      url = "https://github.com/y3owk1n/uts/releases/download/v${version}/uts-linux-amd64";
      sha256 = "";
    };
  };
in

if usePrebuilt && builtins.hasAttr system prebuilt then
  pkgs.stdenv.mkDerivation {
    pname = "uts";
    inherit version;
    src = pkgs.fetchurl {
      url = prebuilt.${system}.url;
      sha256 = prebuilt.${system}.sha256;
    };
    dontUnpack = true;
    installPhase = ''
      install -m755 $src $out/bin/uts
    '';
    meta = {
      description = "All-in-one utility toolkit for media compression and conversion";
      homepage = "https://github.com/y3owk1n/uts";
      license = pkgs.lib.licenses.mit;
      platforms = builtins.attrNames prebuilt;
    };
  }
else
  pkgs.buildGoModule {
    pname = "uts";
    inherit version;
    src = pkgs.lib.cleanSource ./.;
    vendorHash = "";
    subPackages = [ "." ];
    ldflags = [ "-s -w -X github.com/y3owk1n/uts/cmd.Version=${version}" ];
    meta = {
      description = "All-in-one utility toolkit for media compression and conversion";
      homepage = "https://github.com/y3owk1n/uts";
      license = pkgs.lib.licenses.mit;
      platforms = builtins.attrNames prebuilt;
    };
  }
