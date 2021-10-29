# ecs-exec-pf

Port forwarding using the ECS task container. (aws-cli wrapper)

## Usage

```
ecs-exec-pf - Port forwarding using the ECS task container. (aws-cli wrapper)

  Flags:
       --version      Displays the program version string.
    -h --help         Displays help with available flag, subcommand, and positional value parameters.
    -c --cluster      ECS cluster name.
    -t --task         ECS task ID.
    -n --container    Container name in ECS task.
    -p --port         Target remote port. (default: 0)
    -l --local-port   Client local port. (default: 0)
```

## Installation

```sh
brew tap winebarrel/ecs-exec-pf
brew isntall ecs-exec-pf
```

## Execution Example

```sh
$ ecs-exec-pf -c my-cluster -t 0113f61a4b1044d99c627daeee8c0d0c -p 80 -l 8080
Starting session with SessionId: root-03f56652a5f120d48
Port 8080 opened for sessionId root-03f56652a5f120d48.
Waiting for connections...
```

```
$ curl -s localhost:8080 | grep title
<title>Welcome to nginx!</title>
```
