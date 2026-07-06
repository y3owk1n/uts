{ config, pkgs, lib, ... }:

let
  cfg = config.programs.uts;
in {
  options.programs.uts = {
    enable = lib.mkEnableOption "uts - all-in-one utility toolkit";
    package = lib.mkOption {
      type = lib.types.package;
      default = pkgs.uts;
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];
  };
}
