def main(ctx):
    return [
        pipeline(build, "build", []),
        pipeline(deploy, "deploy", ["build"]),
    ]

def pipeline(task, name, depends_on):
    return {
        "kind": "pipeline",
        "name": name,
        "platform": {
            "os": "linux",
            "arch": "arm64",
        },
        "trigger": {
            "branch": ["main"],
            "event": ["push", "custom"],
        },
        "steps": [
            task("gateway"),
            task("messenger"),
            task("posts"),
            task("users"),
        ],
        "depends_on": depends_on,
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

def deploy(service):
    image = "$REGISTRY/$DRONE_REPO_NAME-$DRONE_STEP_NAME:$DRONE_COMMIT_SHA"
    return {
        "name": service,
        "image": "kubesphere/kubectl:v1.22.9",
        "environment": {
            "REGISTRY": {
                "from_secret": "registry",
            },
            "KUBE_CONFIG": {
                "from_secret": "kube_config",
            },
        },
        "commands": [
            "mkdir /root/.kube && echo -n $KUBE_CONFIG | base64 -di > /root/.kube/config",
            "cat ./.kube/$DRONE_STEP_NAME.yaml | sed 's~{{IMAGE}}~'%s'~g' | kubectl apply -f -" % image,
            "kubectl wait --for=condition=available --timeout=60s deployment.apps/$DRONE_STEP_NAME",
        ],
        "depends_on": ["clone"],
        "volumes": [
            {
                "name": "dockersock",
                "path": "/var/run/docker.sock",
            },
        ],
    }
