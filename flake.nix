{
  description = "All-in-one utility toolkit for media compression and conversion";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs =
    { self, nixpkgs, ... }:
    let
      eachSystem = nixpkgs.lib.genAttrs [
        "aarch64-darwin"
        "x86_64-darwin"
        "aarch64-linux"
        "x86_64-linux"
      ];
      latestVersion = "1.0.0";
      makeUtsPackage =
        pkgs: version: usePrebuilt: commitHash:
        pkgs.callPackage ./package.nix {
          inherit version usePrebuilt commitHash;
        };
    in
    {
      overlays.default = final: prev: {
        uts = makeUtsPackage final latestVersion true null;
        uts-source = makeUtsPackage final "main" false (self.rev or self.dirtyRev or "unknown");
      };
      packages = eachSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ self.overlays.default ];
          };
        in
        {
          default = makeUtsPackage pkgs latestVersion true null;
          source = makeUtsPackage pkgs "main" false (self.rev or self.dirtyRev or "unknown");
        }
      );

      homeManagerModules.default = import ./home-module.nix;
    };
}
