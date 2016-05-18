# go-limit
Simple tool to run command with limited resource. Only works on linux for now.

## Install

Your need libcgroup to compile [gp-cgroup](https://github.com/vbatts/go-cgroup). Install
it with your package manager:

On debian:

```bash
sudo apt-get install -y libcgroup-dev
```

On fedora:

```bash
sudo yum install -y libcgroup-devel
```

Install this project:

```bash
go get github.com/inoc603/go-limit
```

## Usage

```bash
go-limit --cpu CPU_USAGE_PERCENTAGE --memory MEMORY_LIMIT YOUR_COMMAND
```

## Example

```bash
# limit cpu usage to 50%
go-limit --cpu 50 example/cpu.py
# limit memory usage to 512MB
go-limit --memory 512M example/mem.py
```

## TODO

- Handle keyboard interruption signal

