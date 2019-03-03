# Bilivedown
A bilibili live stream recording tool

# Usage
To make it easy to execute by some task programs like cron, this program **will not** accept any parameters. 
All parameters will be stored in Setting.ini file. 

# Install/Compile
## via go tool
Just type following command in shell:
``` bash
# Clone the repositories via git
go get github.com/IceSandwich/Bilivedown
# After `go get` command, you can't run it directly. It need Setting.ini file.
cp $GOPATH/src/github.com/IceSandwich/Bilivedown/Setting.ini $GOPATH/bin/Setting.ini
# Edit some parameters.
vim $GOPATH/Setting.ini
# Now you can run this program and recording live stream.
Bilivedown
```
## via binary package
The binary package will provide for windows, linux, android. 
If there is no correct version of your system, you can compile it by yourself.
## via source code
It's simple to build using go tool, type:
``` bash
git clone https://github.com/IceSandwich/Bilivedown.git
cd Bilivedown
go build
```

# How to do with the ts files?
You can merge them like this:
``` bash
# Windows
copy /b *.ts merge.ts
# Linux
cat *.ts > merge.ts
```
Or you can use ffmpeg to merge them.
