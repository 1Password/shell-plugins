# Test nix file to trigger workflow
{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation {
  name = "test";
  src = ./.;
  buildInputs = [ ];
}