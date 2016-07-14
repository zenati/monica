## Monica
[Monica](https://www.youtube.com/watch?v=OY1xxhlq4RU) is a [Go](https://golang.org) project that helps avoid repeating commands by defining a structured `.monica.yml` config file. Monica dynamically generates needed arguments and validates them.

### Installation for OSX—amd64
```
sudo curl -sSo /usr/bin/monica https://raw.githubusercontent.com/zenati/monica/master/main && sudo chmod 777 /usr/bin/monica
```

### Manual download
Every new Monica version is released using Github Releases and the latest release download links are available here:
```
https://github.com/zenati/monica/releases/latest
```

Here are all available plateforms:
```
### Darwin (Apple Mac)

 * [monica\_0.1\_darwin\_386.zip](monica_0.1_darwin_386.zip)
 * [monica\_0.1\_darwin\_amd64.zip](monica_0.1_darwin_amd64.zip)

### FreeBSD

 * [monica\_0.1\_freebsd\_386.zip](monica_0.1_freebsd_386.zip)
 * [monica\_0.1\_freebsd\_amd64.zip](monica_0.1_freebsd_amd64.zip)
 * [monica\_0.1\_freebsd\_arm.zip](monica_0.1_freebsd_arm.zip)

### Linux

 * [monica\_0.1\_amd64.deb](monica_0.1_amd64.deb)
 * [monica\_0.1\_armhf.deb](monica_0.1_armhf.deb)
 * [monica\_0.1\_i386.deb](monica_0.1_i386.deb)
 * [monica\_0.1\_linux\_386.tar.gz](monica_0.1_linux_386.tar.gz)
 * [monica\_0.1\_linux\_amd64.tar.gz](monica_0.1_linux_amd64.tar.gz)
 * [monica\_0.1\_linux\_arm.tar.gz](monica_0.1_linux_arm.tar.gz)

### MS Windows

 * [monica\_0.1\_windows\_386.zip](monica_0.1_windows_386.zip)
 * [monica\_0.1\_windows\_amd64.zip](monica_0.1_windows_amd64.zip)

### NetBSD

 * [monica\_0.1\_netbsd\_386.zip](monica_0.1_netbsd_386.zip)
 * [monica\_0.1\_netbsd\_amd64.zip](monica_0.1_netbsd_amd64.zip)
 * [monica\_0.1\_netbsd\_arm.zip](monica_0.1_netbsd_arm.zip)

### OpenBSD

 * [monica\_0.1\_openbsd\_386.zip](monica_0.1_openbsd_386.zip)
 * [monica\_0.1\_openbsd\_amd64.zip](monica_0.1_openbsd_amd64.zip)

### Other files

 * [control.tar.gz](.goxc-temp/control.tar.gz)
 * [data.tar.gz](.goxc-temp/data.tar.gz)
 * [LICENSE.md](LICENSE.md)
 * [README.md](README.md)

### Plan 9

 * [monica\_0.1\_plan9\_386.zip](monica_0.1_plan9_386.zip)
```

### Example of use
Let's say we need to type almost everyday the following commands in the same directory:
```
$ rake assets:clobber assets:precompile
$ git add -A
$ git commit -m 'Commit message'
$ git push origin master
```

And these too:

```
$ goxc -d=dist -pv=12.31
$ touch src/var/debian/file
$ rm -rf dist/debian-tmp
```

What we also need here is to be able to change the `branch`, `commit message`, `pv` and `debian` using command line arguments.

Here is what defining `actions` in the `.monica.yml` file would look like:

```yaml
actions:
  - name: push
    desc: Pushing current branch to Github
    content:
      - command: rake assets:clobber assets:precompile
      - command: git add -A
      - command: git commit -m "${m}"
      - command: git push ${r} ${b}
    default:
      - m: no-commit-message
      - r: origin
      - b: master

  - name: c #short for compile
    desc: Compiling latest version for all plateforms
    content:
      - command: goxc -d=dist -pv=${pv}
      - command: touch src/var/${a}/file
      - command: rm -rf dist/${a}-tmp
    default:
      - a: debian
```

The config file should be placed inside the directory in which you want to run these commands to be detected and parsed by `monica`. If you use the curl command above to install `monica`, the executable will be named `monica`. Once done, you can call the following command :

```
monica push -b master -m "commit message"
```

Or to use the Goxc example with `pv=16.32` and `a=debian`:

```
monica c --pv 16.32 -a debian
```

And here is the output for the `push` reaction:
```
computer:dir zenati$ m push -b master -m "commit message"
executing: push
-> rake assets:clobber assets:precompile
-> git add -A
-> git commit -m '...'
-> git push origin master
```

## Dynamic arguments
`Monica` detects the config file and dynamically creates the needed mandatory options.
Here is an example of `m --help` output for the example above:
```
computer:dir zenati$ monica --help
usage: monica [<flags>] <command> [<args> ...]

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
      --debug    Enable debug mode.
      --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  push -m=M -b=B
    Pushing current branch to Github

  c --pv=PV -a=A
    Compiling latest version for all plateforms
```

## License
```
The MIT License (MIT)

Copyright (c) 2016 Yassine Zenati

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
