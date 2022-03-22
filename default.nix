{ fetchurl, libarchive, lib, stdenvNoCC, callPackage, overrides ? (self: super: {}) }:
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

    mkApp = pname: value: stdenvNoCC.mkDerivation {
      inherit pname;
      inherit (value) version;
      src = fetchurl {
        inherit (value) url sha256;
        name = "${pname}-${value.version}.zip";
      };
      buildInputs = [ libarchive ];
      dontUnpack = true;
      installPhase = ''
        mkdir -p $out/apps/${pname}
        bsdtar xf "$src" -C $out/apps/${pname}/
      '';
    };
  in
    apps
  );
in lib.fix' (lib.extends overrides apps)
