{ lib, fetchurl, runCommand, callPackage, overrides ? (self: super: {}) }:
let apps = (self:
  let
    appJson = version: builtins.fromJSON (builtins.readFile (./. + "/${version}.json"));
    versions = (callPackage ./nc-versions.nix {}).n;
    apps = builtins.listToAttrs (map (v: let
      majorVer = lib.versions.major v;
    in {
      name = majorVer;
      value = builtins.mapAttrs mkApp (appJson majorVer);
    }) versions);

    mkApp = name: value: runCommand "nc-app-${name}-${value.version}" {
      src = fetchurl {
        inherit (value) url sha256;
      };
      inherit (value) version;
    } /* sh */ ''
      mkdir _unp
      tar -xpf "$src" -C _unp
      if [ $(find _unp -mindepth 1 -maxdepth 1 -type d | wc -l) != 1 ]; then
        echo "error: zip file must contain a single directory"
        exit 1
      fi
      mkdir -p $out
      cp -R _unp/*/. $out/

      if [ ! -f "$out/appinfo/info.xml" ]; then
        echo "appinfo/info.xml doesn't exist in $out, aborting!"
        exit 2
      fi
    '';
  in
    apps
  );
in lib.fix' (lib.extends overrides apps)
