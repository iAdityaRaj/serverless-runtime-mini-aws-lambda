<div align="center">

# ⚡ GoLambda serverless Runtime

### A Lightweight Serverless Function Execution Platform Built in Go

*Grounded in the Berkeley View on Serverless Computing — engineered from first principles.*

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Engine-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com)
[![AWS EC2](https://img.shields.io/badge/AWS-EC2-FF9900?style=flat-square&logo=amazonaws&logoColor=white)](https://aws.amazon.com/ec2/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey?style=flat-square)]()
[![Status](https://img.shields.io/badge/Status-Live-brightgreen?style=flat-square)]()

<img width="1536" height="1024" alt="ChatGPT Image May 27, 2026, 06_17_04 PM" src="https://github.com/user-attachments/assets/67c185b3-7e77-4251-9252-e2ee2c0ee070" />


</div>

---

## Overview

**GoLambda Serverless Runtime** is a production-grade, containerized function execution platform built entirely in Go. It implements the core mechanics of a serverless compute engine: dynamic function registration via source code deployment, isolated container-based invocation, concurrent worker scheduling, and real-time telemetry — all deployed publicly on AWS EC2.

This is not a toy wrapper around Docker. It is an engineered execution substrate that models the architectural primitives found in real FaaS platforms: a runtime abstraction layer, a goroutine-backed worker pool, queue-based scheduling with backpressure, and execution timeout enforcement — built to be extensible, observable, and cloud-deployable.

> **Live Deployment**: The platform runs publicly on AWS EC2 at `http://16.176.26.153:8080`. Functions can be deployed and invoked remotely via the HTTP API with no configuration required on the client side.

---

## Academic Foundation

This project is a practical implementation of the principles laid out in:

> **"Cloud Programming Simplified: A Berkeley View on Serverless Computing"**
> Hellerstein et al., UC Berkeley, 2019 — [arXiv:1902.03383](https://arxiv.org/abs/1902.03383)

The Berkeley paper identified the defining properties of serverless platforms:

| Berkeley Principle | Implementation in GoLambda |
|---|---|
| Decoupled computation and storage | Functions are stateless code units; no shared filesystem between invocations |
| Execution abstraction hides infrastructure | Callers invoke by name — container lifecycle is fully managed by the runtime |
| Millisecond billing granularity | Per-invocation execution duration tracked at millisecond precision |
| Automatic scaling | Worker pool goroutines are always-on; new invocations are queued and dispatched without manual intervention |
| No resource management by users | CPU, memory, and process isolation are enforced by the runtime; callers only submit code |

The paper also characterized **cold starts** as the primary latency cost in FaaS platforms. This runtime deliberately surfaces that cost: every invocation spawns a fresh Docker container, making the cold start penalty measurable, observable, and a first-class engineering concern rather than a hidden abstraction.

---

## Problem Statement

Serverless platforms abstract infrastructure away from developers — but at the cost of opacity. Most engineers use Lambda as a black box. This project answers the question: *what does the compute layer of a FaaS platform actually look like under the hood?*

Specifically, this runtime addresses three core engineering challenges:

1. **Isolation** — Every function invocation must execute in a fully isolated environment, with no shared state, filesystem, or process space between executions.
2. **Concurrency** — The platform must handle multiple simultaneous invocation requests without blocking, while maintaining fair scheduling and bounded resource consumption.
3. **Observability** — A production runtime must surface telemetry: invocation counts, failure rates, average execution latency, queue depth, and worker saturation — queryable in real time.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
│               (curl / HTTP client / Postman)                    │
└────────────────────────────┬────────────────────────────────────┘
                             │ REST API
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API GATEWAY LAYER                          │
│      POST /deploy    POST /invoke/:name    GET /metrics         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                   WORKER POOL SCHEDULER                         │
│                                                                 │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│   │ Worker 1 │  │ Worker 2 │  │ Worker 3 │  │ Worker N │        │
│   │goroutine │  │goroutine │  │goroutine │  │goroutine │        │
│   └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
│        └─────────────┴─────────────┴───────────-─┘              │
│                         Channel Queue                           │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                   RUNTIME ABSTRACTION LAYER                     │
│                  interface Runtime { Run(...) }                 │
│                                                                 │
│                  ┌─────────────────────┐                        │
│                  │   DockerRuntime     │                        │
│                  │  (concrete impl)    │                        │
│                  └─────────┬───────────┘                        │
└────────────────────────────┬────────────────────────────────────┘
                             │ docker run --rm
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                       DOCKER DAEMON                             │
│                                                                 │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│   │  Container   │  │  Container   │  │  Container   │          │
│   │  (fn: hello) │  │  (fn: proc)  │  │  (fn: calc)  │          │
│   │  ephemeral   │  │  ephemeral   │  │  ephemeral   │          │
│   └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### Design Decisions

| Concern | Approach | Rationale |
|---|---|---|
| Isolation | Docker container per invocation | Full process + filesystem isolation; matches Berkeley's stateless-execution model |
| Concurrency | Goroutine worker pool + channel queue | Leverages Go's CSP model; avoids thread-per-request overhead |
| Scheduling | Buffered channel as work queue | Natural backpressure; bounded parallelism without a third-party queue |
| Runtime abstraction | Interface-driven design | Pluggable backends — DockerRuntime today, alternative sandboxes tomorrow |
| Timeouts | `context.WithTimeout` per invocation | Prevents runaway containers from exhausting worker slots |
| Telemetry | In-memory atomic counters | Zero external dependencies; real-time metrics queryable via HTTP |

---

## Features

### Core Execution Engine
- **Concurrent Worker Pool** — Fixed pool of persistent goroutines. Workers pull invocation jobs from a shared channel queue, keeping scheduling fair and throughput bounded.
- **Runtime Abstraction Interface** — A clean `Runtime` interface decouples the scheduler from the underlying execution mechanism. `DockerRuntime` is the current implementation; alternative runtimes can be swapped in without touching the scheduler or API layer.
- **Docker-Sandboxed Execution** — Each function invocation spawns a fresh, ephemeral Docker container. Containers are automatically removed post-execution (`--rm`), leaving no residual state.
- **Execution Timeout Enforcement** — Every invocation runs under a `context.WithTimeout`. Containers that exceed the deadline are forcibly killed and the worker slot is immediately reclaimed.

### API Layer
- **Source Code Deployment** — Register functions by submitting Go source code directly via the deployment API. The runtime compiles and containerizes the function on the host. No pre-built images required.
- **Function Registry** — An in-memory registry maps function names to their compiled artifacts and runtime configuration.
- **Named Invocation API** — Trigger registered functions by name via `POST /invoke/:name`. The platform handles container lifecycle, stdout capture, and error propagation.

### Observability
- **Real-Time Metrics Endpoint** — `GET /metrics` returns a live JSON snapshot of runtime telemetry with no caching.
- **Tracked Signals**:
  - `total_invocations` — Cumulative invocation count across all functions
  - `failed_invocations` — Count of invocations that exited with a non-zero code or timed out
  - `active_workers` — Goroutines currently executing an invocation
  - `queue_depth` — Pending invocations waiting for an available worker
  - `average_execution_ms` — Rolling average of container execution duration in milliseconds
  - `average_queue_wait_ms` — Rolling average of time spent waiting in the work queue

### Infrastructure
- **Cloud-Deployed** — Running publicly on AWS EC2 (Linux, x86_64). Accessible without VPN or authentication.
- **Stateless API Layer** — No session affinity required; horizontally scalable.
- **Zero External Service Dependencies** — Runs with only Go and the Docker daemon. No managed queues, no databases, no sidecars.

---

## Runtime Workflow

```
1. DEPLOYMENT
   Client sends POST /deploy with {name, code}
   → Runtime receives raw Go source code
   → Compiles and registers the function in the FunctionRegistry
   → Returns 201 Created; function is immediately invocable

2. INVOCATION
   Client sends POST /invoke/:name
   → API handler looks up function in registry
   → Creates InvocationJob{fn, responseChan}
   → Job enqueued onto the channel-backed work queue
   → Queue wait timer starts

3. SCHEDULING
   Available Worker goroutine dequeues InvocationJob
   → Queue wait timer stops; delta recorded to average_queue_wait_ms
   → Worker acquires execution context with configured timeout

4. EXECUTION
   Worker calls Runtime.Run(ctx, fn)
   → DockerRuntime spawns: docker run --rm <compiled-fn-image>
   → Captures stdout; enforces timeout via context cancellation
   → On timeout: context cancelled → docker kill → worker reclaimed

5. TEARDOWN
   Container exits; Docker daemon removes it (--rm)
   → Worker updates total_invocations, average_execution_ms
   → Result written to responseChan
   → API handler reads result and returns HTTP 200

6. TELEMETRY
   GET /metrics returns live JSON counters and averages
   → All values are computed from in-memory atomic state
   → No scrape interval; every request reflects current state
```

---

## API Reference

### Deploy a Function

Submit Go source code to register a named function. The runtime handles compilation and registration.

```
POST /deploy
Content-Type: application/json
```

**Request**
```json
{
  "name": "hello",
  "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n fmt.Println(\"Hello from AWS runtime\")\n}"
}
```

**Response** `201 Created`
```json
{
  "status": "deployed",
  "function": "hello"
}
```

---

### Invoke a Function

Invoke a registered function by name. The platform schedules execution, manages the container lifecycle, and returns the output synchronously.

```
POST /invoke/:name
```

**Example**: `POST http://16.176.26.153:8080/invoke/hello`

**Response** `200 OK`
```json
{
  "output": "Hello from AWS runtime\n"
}
```

> Function name is passed as a URL path parameter. No request body is required for basic invocations.

---

### Query Metrics

Returns a live snapshot of runtime telemetry. All values reflect the current in-memory state at the moment of the request.

```
GET /metrics
```

**Response** `200 OK`
```json
{
  "total_invocations": 1,
  "failed_invocations": 0,
  "active_workers": 0,
  "queue_depth": 0,
  "average_execution_ms": 398.439524,
  "average_queue_wait_ms": 0.00384
}
```

> `average_execution_ms` of ~398ms reflects the Docker cold start cost for a minimal Go binary — consistent with the cold start latency range documented in the Berkeley serverless paper.

---

## Project Structure

```
serverless-runtime/
├── main.go                    # Entry point; wires server, pool, registry
├── api/
│   ├── server.go              # HTTP router and handler registration
│   ├── deploy.go              # POST /deploy — source intake, compile, register
│   └── invoke.go              # POST /invoke/:name — schedule and execute
├── runtime/
│   ├── runtime.go             # Runtime interface definition
│   └── docker.go              # DockerRuntime — container lifecycle management
├── scheduler/
│   ├── pool.go                # Worker pool; goroutine lifecycle management
│   ├── queue.go               # Channel-backed invocation queue
│   └── job.go                 # InvocationJob struct
├── registry/
│   └── registry.go            # Thread-safe FunctionRegistry (sync.RWMutex)
├── metrics/
│   └── metrics.go             # Atomic counters; rolling averages; JSON serializer
├── config/
│   └── config.go              # Runtime configuration (worker count, timeouts)
├── Dockerfile                 # Platform container image
├── docker-compose.yml         # Local development stack
└── README.md
```

---

## Local Development

### Prerequisites

- Go 1.21+
- Docker Engine (daemon must be running)
- `curl` or any HTTP client

### Setup

```bash
# Clone the repository
git clone https://github.com/iAdityaRaj/serverless-runtime-mini-aws-lambda.git
cd serverless-runtime-mini-aws-lambda

# Install dependencies
go mod download

# Verify Docker is accessible
docker info

# Start the platform (default: port 8080)
go run main.go
```

### Configuration via Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server listen port |
| `WORKER_COUNT` | `10` | Number of concurrent goroutine workers |
| `QUEUE_DEPTH` | `100` | Max buffered invocation queue depth |
| `DEFAULT_TIMEOUT` | `30s` | Per-invocation execution timeout |

### Quick Test

```bash
# 1. Deploy a Go function by submitting source code
curl -X POST http://localhost:8080/deploy \
  -H "Content-Type: application/json" \
  -d '{
    "name": "hello",
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n fmt.Println(\"Hello from AWS runtime\")\n}"
  }'

# 2. Invoke it
curl -X POST http://localhost:8080/invoke/hello

# 3. Query live metrics
curl http://localhost:8080/metrics
```

---

## AWS Deployment

The platform is deployed on AWS EC2 and publicly accessible at `http://16.176.26.153:8080`.

### EC2 Setup

```bash
# 1. Launch EC2 instance (t3.small or larger recommended)
#    Security group: allow inbound TCP on ports 22 and 8080

# 2. SSH into the instance
ssh -i your-key.pem ec2-user@16.176.26.153

# 3. Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 4. Install Docker
sudo yum update -y && sudo yum install docker -y
sudo systemctl start docker
sudo usermod -aG docker ec2-user
newgrp docker

# 5. Clone and build
git clone https://github.com/iAdityaRaj/serverless-runtime-mini-aws-lambda.git
cd serverless-runtime-mini-aws-lambda
go build -o golambda .

# 6. Run as persistent background process
nohup ./golambda > runtime.log 2>&1 &
```

### Docker Compose (Recommended)

```bash
docker-compose up -d
```

```yaml
# docker-compose.yml
version: "3.9"
services:
  runtime:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - WORKER_COUNT=10
      - QUEUE_DEPTH=100
      - DEFAULT_TIMEOUT=30s
    restart: unless-stopped
```

> **Note on Docker socket**: Mounting `/var/run/docker.sock` grants the runtime access to the host Docker daemon — the standard Docker-out-of-Docker (DooD) pattern for single-host deployments. In a production multi-tenant environment this would be replaced with a scoped Docker API proxy.

### Security Group Rules

| Type | Port | Source | Purpose |
|---|---|---|---|
| SSH | 22 | Your IP | Administration |
| Custom TCP | 8080 | 0.0.0.0/0 | API + Metrics |

---

## Cold Start Behavior

This runtime uses a **cold start per invocation** model — deliberately, to surface and measure the cost that Berkeley's serverless paper identifies as the primary latency challenge in FaaS platforms.

- Every invocation triggers: `docker run` → container init → process start → execution → `--rm` teardown.
- **Observed cold start overhead**: ~398ms average on EC2 for a minimal Go binary (visible in `/metrics`).
- **Workers are warm** — goroutines in the pool are always running and immediately available. The latency cost is entirely in container spin-up, not worker scheduling.

This is a deliberate architectural choice. The `average_execution_ms` metric directly quantifies this cost, making it observable rather than opaque — which is how it should be in a platform you want to understand and improve.

---

## Observability

`GET /metrics` returns live JSON telemetry with no external tooling required:

```json
{
  "total_invocations": 1,
  "failed_invocations": 0,
  "active_workers": 0,
  "queue_depth": 0,
  "average_execution_ms": 398.439524,
  "average_queue_wait_ms": 0.00384
}
```

### What Each Signal Tells You

| Signal | Interpretation |
|---|---|
| `average_execution_ms` rising | Cold start degradation; image pull latency; host I/O pressure |
| `queue_depth` > 0 sustained | Worker pool saturated; invocation backpressure building |
| `failed_invocations` climbing | Runtime error or timeout in function code; inspect container logs |
| `active_workers` = worker count | All workers busy; queue will grow until one completes |
| `average_queue_wait_ms` near zero | Workers draining the queue faster than requests arrive (healthy) |

---

## Engineering Highlights

These design choices reflect real production platform engineering concerns:

**Interface-Driven Runtime** — The `Runtime` interface is the primary extension point of the system. Swapping Docker for Firecracker MicroVMs, gVisor, or a Wasm interpreter requires only a new struct implementing a single interface — zero changes to the scheduler or API layer. This is the same boundary abstraction used in Kubernetes' Container Runtime Interface (CRI).

**Backpressure via Bounded Channels** — The work queue is a buffered Go channel. When the channel is full, the API layer immediately returns `429 Too Many Requests` rather than spawning unbounded goroutines. This is a deliberate systems decision: explicit, fail-fast backpressure is safer than implicit memory growth under load.

**Timeout Propagation via Context** — Execution contexts are created with `context.WithTimeout` and threaded into `Runtime.Run()`. The Docker executor propagates context cancellation to `docker kill`, ensuring that a deadline at the API layer reliably terminates the underlying OS process — not just the Go-side wait.

**In-Memory Telemetry Without a Framework** — Metrics use `sync/atomic` for counter operations and running averages. There is no external metrics library. The `/metrics` handler serializes live state directly into JSON — zero dependency, zero scrape lag, zero configuration.

**Worker Lifecycle Ownership** — The pool fully owns its goroutines. Workers loop on `for job := range queue`, blocking efficiently on the channel rather than spin-waiting. The pool exposes `Start()` and `Stop()` with graceful drain semantics — in-flight invocations complete before shutdown.

**Source-Level Function Deployment** — Functions are deployed as raw Go source code, not pre-built images. The platform compiles and containerizes on the host, mirroring how managed FaaS platforms accept source and return a handle — abstracting the build pipeline from the caller entirely.

***Warm Starts vs Cold Starts*** - This runtime currently uses a cold-start execution model: every invocation launches a fresh Docker container, executes the function, and destroys the container afterward.

The worker goroutines themselves are warm and persistent, meaning scheduling overhead is minimal, while container startup latency dominates execution time.



<div align="center">

*Built on the foundations of the Berkeley View on Serverless Computing.*

*Engineered in Go. Deployed on AWS. Observable by design.* | Aditya Raj (2025AIM1001)



</div>
