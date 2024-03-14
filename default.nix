{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
}:

pkgs.buildGoApplication {
  pname = "nc4nix";
  version = "unstable-2023-06-06";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
}
