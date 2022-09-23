# TODO

This is a punch list of the things that need to be done to this CLI for the
demo, in order of importance.

* Modify builder image (Jenkins) to start locally with a hard-coded user and
    an API key created.
* Modify this CLI to start the builder image with a directory mapped for data,
    such as `~/.atk/builder`
* Modify the CLI to provide credentials to the builder image
* Modify bifrost to create the projects with real code (e.g., `tf apply`)

And then test, test, test...
