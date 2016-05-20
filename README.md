## Monica
Monica is a [Go](https://golang.org) project that helps developers avoid repeating commands by defining a structured `.monica.yml` config file using dynamic arguments generation and validation.

### Installation
```
sudo curl -sSo /usr/bin/m https://raw.githubusercontent.com/zenati/monica/master/monica && sudo chmod 777 /usr/bin/m
```

### Manual download
Every new Monica version is released using Github Releases and the latest release download links are available here:
```
https://github.com/zenati/monica/releases/latest
```

Here are all available plateforms:
```
i386
amd64
armhf
darwin_386
darwin_amd64
freebsd_386
freebsd_amd64
freebsd_arm
linux_386
linux_amd64
linux_arm
netbsd_386
netbsd_amd64
netbsd_arm
windows_386
windows_amd64
```

### Example of use
Let's say we need to type almost everyday the following commands:
- `rake assets:clobber assets:precompile`
- `git add -A`
- `git commit -m 'Commit message'`
- `git push origin master`

And these too:
- `goxc -d=dist -pv=12.31`
- `rm -rf /dist/debian-tmp`

What we need here, is also to be able to change the `branch`, `commit message` and `pv`.

Here is what the `.monica.yml` file would look like:

```yaml
engine: monica
reactions:
  - name: push
    desc: Pushing current branch to Github
    content:
      - command: rake assets:clobber assets:precompile
      - command: git add -A
      - command: git commit -m '${m}'
      - command: git push origin ${b}

  - name: c #short for compile
    desc: Compiling Goxc for all plateforms upon release
    content:
      - command: goxc -d=dist -pv=${pv}
      - command: rm -rf /dist/debian-tmp
```

The config file should be placed at the root of the git repository to be detected and parsed by `monica`.
Once done, you can call the following command :

```
m push -b master -m "commit message"
```

Or to use the Goxc example with `16.32` as `pv`:

```
m c -pv 16.32
```

And here is the output:
```
monica executing: push
monica 	-> rake assets:clobber assets:precompile
monica 	-> git add -A
monica 	-> git commit -m '...'
monica 	-> git push origin master
```

## Dynamic arguments
`Monica` detects the config file and dynamically creates the needed mandatory options.
Here is an example of `monica --help` output for the example above:
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

  push --m=M --b=B
    Pushing current branch to Github

  compile --pv=PV
    Compiling Goxc for all plateforms and cleaning files
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