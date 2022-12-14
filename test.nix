{ pkgs ? import <nixpkgs> {}, lib ? pkgs.lib }:
let
  nc4nix = pkgs.callPackage ./. {};
in {
  inherit (nc4nix."24")
    breezedark
    drawio
    groupfolders
    onlyoffice;
}
