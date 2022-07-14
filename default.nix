{ lib, fetchNextcloudApp, callPackage, overrides ? (self: super: {}) }:
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

    mkApp = name: value: fetchNextcloudApp {
      inherit name;
      inherit (value) version url sha256;
    };
  in
    apps
  );
in lib.fix' (lib.extends overrides apps)
