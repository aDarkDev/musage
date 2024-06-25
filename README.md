# musage
RSS Memory Measurement Tool

A real-time memory measurement tool for monitoring memory usage during program execution.

## installation:
```bash
wget https://github.com/aDarkDev/musage/releases/download/v0.1/musage-linux; mv musage-linux /usr/bin/musage; chmod +x musage
```

## usage:
```bash
usage: musage <command>
example: musage sleep 10; echo hi
```

### output [musage.log]:
```json
{"process_id": 60238,"memory_usage": 474.062500, "timestamp": 1719317687, "command": "/home/user/Desktop/Telegram"}
```
note `memory_usage` is MB
