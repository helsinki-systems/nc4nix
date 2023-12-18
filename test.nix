{ pkgs ? import <nixpkgs> {}, lib ? pkgs.lib }:
let
  nc4nix = pkgs.callPackage ./. {};
  ncVersions = map lib.versions.major (import ./nc-versions.nix { inherit pkgs; }).n;
in builtins.map (v: lib.recurseIntoAttrs {
    inherit (nc4nix."${v}")
      groupfolders
      onlyoffice
      spreed   # aka talk
      twofactor_webauthn
    ;
  }
) ncVersions
