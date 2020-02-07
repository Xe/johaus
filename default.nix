{ pkgs ? import <nixpkgs> {} }:
with pkgs;

assert lib.versionAtLeast go.version "1.13";

buildGoPackage rec {
  name = "johaus";
  version = "1.1.1-dev";
  goPackagePath = "within.website/johaus";
  src = ./.;
  nativeBuildInputs = [ makeWrapper ];

  goDeps = ./deps.nix;
  allowGoReference = false;
  CGO_ENABLED = "0";
}
