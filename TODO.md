# TODO

This is a punch list of the things that need to be done to this CLI for the
demo, in order of importance.

* Implement the `itz solution deploy --solution xyz --reservation abc123` command
* Implement the `itz solution deploy --solution xyz --cluster-name abc123` command
* Get credentials from reservation web services to make `--reservation abc123` work
* Modify bifrost to create the pipeline with the parameters from TFVars from the file
* Modify the CLI to provide credentials to the builder image from the ocpnow configuration
* Modify the CLI to provide credentials to the builder image from the reservation
* Modify bifrost to create the projects with real code (e.g., `tf apply`)

# TO-DONE!

* Modify builder image (Jenkins) to start locally with a hard-coded user.
* Create an API key.
* Modify this CLI to start the builder image with a directory mapped for data,
  such as `~/.itz/builder`
* Modify the CLI to submit the file to local bifrost services
* Modify the CLI to start the pipeline

And then test, test, test...

