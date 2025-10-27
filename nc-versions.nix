# this file is used to figure out which versions of nextcloud we have in nixpkgs
{ pkgs ? import <unstable> {}, lib ? pkgs.lib }:
let
  n = lib.mapAttrsToList (_: v: v.version) (
      lib.filterAttrs (k: v: builtins.match "nextcloud[0-9]+" k != null && (builtins.tryEval v.version).success)
    pkgs);
in {
  inherit n;
  e = lib.concatStringsSep "," n;
}
