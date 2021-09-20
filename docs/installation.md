# Installation

Binary distributions are currently only available for linux/amd64. They can be downloaded from the [releases](https://github.com/metno/mkharp/releases/) page. Download the latest, and unpack into a folder where your path is set, for example like this:

```bash
wget https://github.com/metno/mkharp/releases/download/v0.1.0/mkharp-v0.1.0-linux-amd64.tar.gz
tar xvzf mkharp-v0.1.0-linux-amd64.tar.gz
sudo cp linux-amd64/mkharp /usr/local/bin/
hash -r
```

It should then be possible to run mkharp from the command-line, like this:

```bash
mkharp -help
```


