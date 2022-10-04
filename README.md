# Lomba

Analyze Kubernetes cluster logs in an instant. Powered by Grafana and Loki.

## What does it do?

Lomba does the following things:

1. Spins up a local loki container and grafana container.
2. Reads all of the pod logs from a kubernetes cluster and ingests them into loki. 

This empowers you to quickly analyze the current state of a cluster without making any modifications to the cluster state.

## Prerequisites

* kubeconfig with access to logs on a Kubernetes cluster.
* linux or macos machine with docker setup.

## Install Lomba

install from the latest [release artifacts](https://github.com/zawachte/lomba/releases):

* Linux

  ```sh
  curl -LO https://github.com/zawachte/lomba/releases/download/v0.0.1/lomba
  mv lomba /usr/local/bin/
  ```

* macOS

  ```sh
  curl -LO https://github.com/zawachte/lomba/releases/download/v0.0.1/lomba-darwin
  mv lomba-darwin /usr/local/bin/lomba
  ```

## Usage

If you have a kubeconfig at `~/.kube/config`

```sh
lomba run
```

To pass a kubeconfig

```
lomba run --kubeconfig kubeconfig
```

Now go to your favorite browser at `http://localhost:3000` and analyze your pod logs!

## Building 

Running make should be enough to get you going:

```
make
```

## Roadmap
