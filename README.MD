pkgparser
=========

Pkgparser is a simple to translate the pacman like output into a dialog checklist, 
for picking up packages conveniently.
The selected packages will be copied into clipboard.

Usage
=========
```bash
sudo pacman -Ss neovim | pkgparser
sudo powerpill -Ss neovim | pkgparser
```

Installation
=========
### Go
```bash
 go get "github.com/ssfdust/pkgparser"
```
