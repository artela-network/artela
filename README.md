<h1 align="center"> Artela </h1>

<div align="center">
  <a href="https://t.me/artela_official" target="_blank">
    <img alt="Telegram Chat" src="https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fpolaris_devs">
  </a>
  <a href="https://twitter.com/Artela_Network" target="_blank">
    <img alt="Twitter Follow" src="https://img.shields.io/twitter/follow/Artela_Network">
  <a href="https://discord.gg/artela">
   <img src="https://img.shields.io/badge/discord-join%20chat-blue.svg" alt="Discord">
  </a>
</div>

**artela** is a blockchain built using Cosmos SDK and Cometbft and created with [Ignite CLI](https://ignite.com/cli).

## Build the source

1). Set Up Your Go Development Environment<br />
Make sure you have set up your Go programming language development environment.

2). Download the Source Code<br />
Obtain the project source code using the following method:

```
git clone https://github.com/artela-network/artela.git
```

3). Compile<br />
Compile the source code and generate the executable using the Go compiler:

```
cd artela
make clean && make
```

## Executables

The artela project comes with executable found in the `build` directory.
|  Command   | Description|
| :--------: | ----------------------------------------------------------------------------------------------------------------|
| **`artelad`** | artelad is the core node software of the Artela network, responsible for running and managing the Artela blockchain network. |

## Running Testnet

### Setting Up a Single-Node Ethermint Testnet

Initialize the testnet by running a simple script<br />

```
sh init.sh
```

Running the node by

```
artelad start
```

### Setting Up a 4-Validator Testnet

Run the following command to initialize a 4-validator testnet:

```
make start-testnet
```

To view node logs, use `tail -f _testnet/node2/artelad/node.log`.<br />

This command compiles the current repository code and generates Docker images. It launches Docker containers named artela0 to artela3, each running artelad as a node. Disk mapping is configured to map directories artea/_testnet/node0 to node3 to Docker containers artela0 to artela3. <br />

More options about the testnet:
| Command           | Description                                                                                     |
| ----------------------- | ----------------------------------------------------------------------------------------------- |
| <span style="white-space:nowrap">`build-testnet`</span> | Build Docker images for the testnet and create a configuration for 4-validator nodes.           |
| <span style="white-space:nowrap">`create-testnet`</span> | Remove a previously built testnet, build it again using `build-testnet`, and start Docker containers. |
| <span style="white-space:nowrap">`stop-testnet`</span> | Stop the running Docker containers for the testnet.                                             |
| <span style="white-space:nowrap">`start-testnet`</span> | Start the previously stopped Docker containers for the testnet.                                 |
| <span style="white-space:nowrap">`remove-testnet`</span> | Stop the Docker containers and remove all components created by the `build-testnet` command.    |

### Hardware Requirements

Minimum:

* CPU with 2+ cores
* 4GB RAM
* 1TB free storage space to sync the Tetnet
* 8 MBit/sec download Internet service

Recommended:

* Fast CPU with 4+ cores
* 16GB+ RAM
* High-performance SSD with at least 1TB of free space
* 25+ MBit/sec download Internet service

## Artela Testnet

Presistent nodes of artela testnet:

* `https://testnet-rpc1.artela.network`
* `https://testnet-rpc2.artela.network`
* `https://testnet-rpc3.artela.network`
* `https://testnet-rpc4.artela.network`

---
Find more about Artela in <https://www.artela.network/>
