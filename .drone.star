def main(ctx):
    return {
        "kind": "pipeline",
        "name": "build",
        "platform": {
            "os": "linux",
            "arch": "arm64",
        },
        "trigger": {
            "branch": ["main", "dev"],
            "event": ["push"],
        },
        "steps": [
            build("gateway"),
            build("messenger"),
            build("posts"),
            build("users"),
        ],
        "volumes": [
            {
                "name": "dockersock",
                "host": {
                    "path": "/var/run/docker.sock",
                },
            },
        ],
    }

def build(service):
    image = "$REGISTRY/$DRONE_REPO_NAME-$DRONE_STEP_NAME:$DRONE_COMMIT_SHA"
    return {
        "name": service,
        "image": "docker:20.10.17-dind",
        "environment": {
            "REGISTRY": {
                "from_secret": "registry",
            },
            "REGISTRY_USER": {
                "from_secret": "registry_user",
            },
            "REGISTRY_PASSWORD": {
                "from_secret": "registry_password",
            },
        },
        "commands": [
            "echo $REGISTRY_PASSWORD | docker login $REGISTRY -u $REGISTRY_USER --password-stdin",
            "docker build -f ./%s/Dockerfile -t %s ." % (service, image),
            "docker push %s" % image,
            "docker rmi %s" % image,
        ],
        "depends_on": ["clone"],
        "volumes": [
            {
                "name": "dockersock",
                "path": "/var/run/docker.sock",
            },
        ],
    }
