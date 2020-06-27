# damb - Docker As a Monorepo Build-system

If it's dumb but it works it ain't damb.

![Lint Status](https://github.com/stiletto/damb/workflows/Lint/badge.svg)
![Tests Status](https://github.com/stiletto/damb/workflows/Tests/badge.svg)

## Description

**damb** is a dependency management tool for Dockerfiles in a monorepo.

Before using **damb** in you repo, you have to create `.damb.yml` at the top level of your repository:

```yaml
build_cmd: ["sudo", "docker", "build"] # a command used to build docker images. damb also works with podman
args: # args are passed as --build-arg's to build_cmd
  damb_prefix: # this arg is special, it's used as a prefix for all images built by damb
    val: example.com/johndoe/awesomerepo/
  damb_tag: # this arg is also special, it's used as image tag for all images built by damb
    cmd: ["git", "symbolic-ref", "--short", "HEAD"]
  # build arguments may be static (like damb_prefix) or dynamic (like damb_tag). dynamic args are evaluated on first use.
```

Every Dockerfile in your repository gets assigned an image name in the form of `${damb_prefix}/path/to/directory[.suffix]:${damb_tag}`.

For example if you have a dockerfile named `projecta/hello/Dockerfile` in a monorepo configured using example .damb.yml, the name of its image will be `example.com/johndoe/awesomerepo/projecta/hello:master`.
A dockerfile named `projectb/Dockerfile.test` will get built into `example.com/johndoe/awesomerepo/projectb.test:master`.

When a Dockerfile refers to an image in a `FROM` directive, **damb** checks if this image starts with `${damb_prefix}` (example.com/johndoe/awesomerepo/) and ends with `:${damb_tag}` (master).
If it does, the image is considered to be an *internal* dependency, if it doesn't - an *external* one.

But how do you refer to an image using `damb_tag` and `damb_prefix` if `damb_tag` is evaluated at run time and `damb_prefix` may be changed via `.damb.yml`?
By defining build arguments before `FROM` of course:

```
ARG damb_prefix
ARG damb_tag
FROM ${damb_prefix}projecta/foo:${damb_tag}
```

Build arguments defined before the first `FROM` directive [are considered to be outside of any stages and cannot be referenced from ordinary build commands](https://docs.docker.com/engine/reference/builder/#from). If you'd like to use an argument to parametrize `FROM` and to parametrize a build step at the same time - define it multiple times.

All the weight of the damb is concentrated on two commands:

* `damb resolve` - resolves target dependencies and outputs them to stdout
* `damb build` - recursively builds images and their dependencies

Check out [examples](examples/) and try to build them with `damb build --no-capture all`!

## License

Copyright 2020 Stiletto <blasux@blasux.ru>

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
