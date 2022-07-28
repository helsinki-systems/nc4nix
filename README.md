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

---

The `ci` script is run daily by our CI and updates all apps.
It basically runs the `main.go` script and generates a commit message.
