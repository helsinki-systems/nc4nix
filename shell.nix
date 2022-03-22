with import <nixpkgs> { };

stdenv.mkDerivation {
  name = "go";
  buildInputs = [
    delve
    libcap
    go
    gcc
  ];
  shellHook = ''
    export GOPATH=$PWD/gopath
  '';
}
