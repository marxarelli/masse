# Examples

This example shows how one might use Massé to build a container image for a
basic Go microservice using.

To build the image and load it into docker for use, run the following in this
directory.

```console
$ docker buildx build -f masse.cue --target service -t localhost/cowsayd --load .
```

Now we can run our service!

```console
$ docker run -d --rm --name cowsayd -p 3000:3000 localhost/cowsayd
$ curl -d 'Hurray for cowsay!' http://localhost:3000
 ____________________
< Hurray for cowsay! >
 --------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
$ docker stop cowsayd
```
