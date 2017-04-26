# karaberry
Karaoke for Raspberry Pi

## install dependencies
* go-lang (https://github.com/tgogos/rpi_golang)
```bash
sudo apt-get install curl git make binutils bison gcc build-essential
# install gvm
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /home/pi/.gvm/scripts/gvm

# install go
gvm install go1.4 
gvm use go1.4 
export GOROOT_BOOTSTRAP=$GOROOT 
gvm install go1.7
gvm use go1.7 --default
# Configuration of $PATH, $GOPATH variables
# Add `/usr/local/go/bin` to the PATH environment variable.
# You can do this by adding this line to your `/etc/profile` (for a system-wide installation)
# or `$HOME/.profile`: 

export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
```

* omxplayer
```bash
sudo apt-get install omxplayer
```
* youtube-dl #latest version
```bash
sudo -s
curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
chmod a+rx /usr/local/bin/youtube-dl
exit
```
## build & run
```bash
go build
go get <dependency> #install any dependencies if required
./karaberry #point browser to localhost:8080
```
