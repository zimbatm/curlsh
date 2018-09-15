{ pkgs ? import <nixpkgs> {} }:
pkgs.buildGoPackage {
  name = "curlsh";
  goPackagePath = "github.com/zimbatm/curlsh";
  src = ./.;
}
