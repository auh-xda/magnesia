# Magnesia

Magnesia is an open-source **Remote Monitoring and Management (RMM)** agent written in **Go**. It collects system information, including CPU, memory, disk, network interfaces, and sends it to a server via NATS.

## Features

* Cross-platform agent (Linux first, expandable to Windows/macOS)

* Collects detailed system info:

  * Hostname, OS, OS version, uptime, boot time

  * Network interfaces, MAC addresses, public IP

  * Memory usage

  * Disk usage per mountpoint

  * CPU information (model, cores, speed, usage)

* Sends data via WebSocket or NATS

* Minimal dependencies, lightweight agent

## Installation

1. Clone the repository:

```
git clone [https://github.com/auh-xda/magnesia.git](https://github.com/auh-xda/magnesia.git)
cd magnesia

```

2. Build the agent:

```
go build -o magnesia main.go

```

3. Run the agent:

```
sudo ./magnesia-agent \
  -action install \
  -auth_token f8a7c3d9b1e2f6a4 \
  -api_key z9x8y7w6v5u4t3s2 \
  -client_id 12873 \
  -client_secret 5Fh8jK2qLp9Z


```

## Usage

The agent automatically collects system information and sends it to the configured server.

### Example payload:

```
{
  "version": "0.1.0",
  "hostname": "lwspc43.localhost",
  "public_ip": "49.36.193.184",
  "os": "linux",
  "os_version": "22.04",
  "uptime": 132573,
  "boot_time": 1755447485,
  "host_id": "20211112-70a6-cccc-6983-70a6cccc6987",
  "family": "debian",
  "interfaces": [
    {
      "name": "wlp0s20f3",
      "mac": "70:a6:cc:cc:69:83",
      "ip_addresses": ["192.168.3.53"]
    }
  ],
  "memory": {
    "total": 16542625792,
    "used": 7348465664,
    "free": 688787456,
    "usage_percent": 44.42
  },
  "disks": [
    {
      "device": "/dev/nvme0n1p11",
      "mountpoint": "/",
      "fstype": "ext4",
      "total": 68319354880,
      "used": 57195327488,
      "free": 7606738944,
      "usage_percent": 88.26
    }
  ],
  "cpu": {
    "manufacturer": "GenuineIntel",
    "cpu_speed_mhz": 4400,
    "cores": 4,
    "model": "11th Gen Intel(R) Core(TM) i5-11300H @ 3.10GHz",
    "sockets": 1,
    "cores_per_socket": 4,
    "logical_processors": 8,
    "hyperthread": true
  }
}

```

## Configuration

* **WebSocket:** set the target URL and channel in the agent.

* **NATS:** configure the host and port in nats.conf.
