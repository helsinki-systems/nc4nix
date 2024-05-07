{ pkgs ? import <nixpkgs> {}, lib ? pkgs.lib }:
let
  nc4nix = pkgs.callPackage ./. {};
  ncVersions = map lib.versions.major (import ./nc-versions.nix { inherit pkgs; }).n;
in builtins.map (v: lib.recurseIntoAttrs {
    inherit (nc4nix."nextcloud-${v}")
      # groupfolders FIXME: https://github.com/nextcloud/groupfolders/issues/2940
      onlyoffice
      spreed   # aka talk
      twofactor_webauthn
    ;
  }
) ncVersions
