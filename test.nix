{ pkgs ? import <nixpkgs> {}, lib ? pkgs.lib }:
let
  nc4nix = pkgs.callPackage ./. {};
in
  (with nc4nix."23"; [
    groupfolders
  ])
