## NewChainDump 

`NewChainDump` project contains the following:
* Extract NewChain blocks and transactions into database

## QuickStart

### Download from releases

Binary archives are published at https://release.cloud.diynova.com/newton/NewChainDump/.

### Building the source

install:

```bash
git clone https://github.com/newtonproject/newchain-dump.git && cd newchain-dump && make install
```

run NewChainDump:

```bash
%GOPATH%/bin/newchain-dump
```

## Usage

#### Help

Use command `NewChainDump help` to display the usage.

```bash
Usage:
  NewChainDump [flags]
  NewChainDump [command]

Available Commands:
  help        Help about any command
  init        Initialize config file
  run         Get NewChain blocks and store in database
  version     Get version of NewChainDump CLI

Flags:
  -c, --config path       The path to config file (default "./config.toml")
      --database string   The name of database
  -h, --help              help for NewChainDump
      --host string       The host for database (default "127.0.0.1:3306")
  -l, --log string        The path of log file (default "./error.log")
      --password string   The password for database
  -i, --rpcURL url        Geth json rpc or ipc url (default "https://rpc1.newchain.newtonproject.org")
      --user string       The user for database

Use "NewChainDump [command] --help" for more information about a command.
```

#### Use config.toml

You can use a configuration file to simplify the command line parameters.

One available configuration file `config.toml` is as follows:


```conf
log = "./error.log"
rpcurl = "https://rpc1.newchain.newtonproject.org"

[mysql]
  database = "newchaindb"
  host = "127.0.0.1"
  password = "password"
  user = "newchain"
```

#### Initialize config file

```bash
# Initialize config file
$ NewChainDump init
Initialize config file
Enter file in which to save (./config.toml):
Enter path of log file (./error.log):
Enter geth json rpc or ipc url (https://rpc1.newchain.newtonproject.org):
Configure MySQL database or not: [Y/n]
Enter database host(127.0.0.1:3306): 127.0.0.1
Enter database name: newchain
Enter the username to connect to the database: newchain
Enter the password for user:
Your configuration has been saved in  ./config.toml
```

#### Run NewChainDump

```bash
# Extract NewChain blocks, start with 0 and end with 100
NewChainDump run --start 0 --end 100

# Extract NewChain blocks without stop
NewChainDump run --loop

# Extract NewChain blocks without stop and with 5 block delay
NewChainDump run --loop --delay 5
```