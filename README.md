# nc4nix

So you want to roll out Nextcloud on NixOS but you also want to use Nix to manage your apps instead of the builtin apps system?
You came to the right place.

This default.nix expression contains the code to handle all apps Nextcloud has to offer.
It does that by parsing pre-generated JSON files with all apps.
The files are pre-generated using the `main.go` script.

Apps are provided for all Nextcloud versions currently in nixpkgs.

## Generating the JSONs

The main.go script (by default) parses **all** apps for Nextcloud.

There also is an environment varaible, called `COMMIT_LOG`.
If set to `1`, logs are generated.
This is used by the `ci` script.

## About
We develop this software we made this software for our own usage.
You are free to use it and open issues. We will look through them and decide if this is an issue to our use case, thus we are not able to address all of them.
But do not hesitate to send a pull request!
If you need this software but do not find the time to the development in house, we also offer professional commerical nixOS support - contact us by mail via [kunden@helsinki-systems.de](mailto:kunden@helsinki-systems.de)!


---

The `ci` script is run daily by our CI and updates all apps.
It basically runs the `main.go` script and generates a commit message.
