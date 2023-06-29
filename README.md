## How to use

Either install Go or run the provided binary. Go binaries are (mostly) self-contained. (A majore exception would be glibc).
Binaries can be found in the current release.

Options:
--username, Docker username
--registry_prefix, Prefix appended to all images' name
--modules_path, Path to explore for finding Docker Bake targets
--log, Log level

## Install Go

```bash
rm -rf /usr/local/go
wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz

# For bash
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
source $HOME/.profile
# For fish
fish_add_path /usr/local/go/bin

go version
```
## Developement

### Managing Go dependencies

```bash
go get --upgrade <package>
# Import in your package
go mod tidy # will take care of adding it to go.mod and go.sum
```