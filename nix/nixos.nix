{ pkgs, lib, config, ... }:
import ./shell-plugins.nix {
  inherit pkgs;
  inherit lib;
  inherit config;
  is-home-manager = false;
}
