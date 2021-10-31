# GRPCHub - A Server-side GRPC Load Balancing Example

[![GitHub license](https://img.shields.io/github/license/ccwutw/grpchub)](https://github.com/ccwutw/grpchub/blob/main/LICENSE) 
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/ccwutw/grpchub/blob/main/CONTRIBUTING.md)

## Overview
gRPC is a modern RPC protocol which integrates service discovery, name resolver, load balancer, tracing and others.
This project will focus on server-side load balancing (a.k.a proxy load balancing).
In server-side load balancing, gRPC clients connect to a gRPC Load Balancer (LB) proxy instead of a gRPC server. 
The LB distributes incoming RPC calls from clients to some of the available gRPC servers that implement the actual logic for serving the calls. 
The LB keeps track of load on each server and distributes load via different algorithms. The clients themselves do not know about the backend servers.
To compare with the client-side load balancing, the proxy architecture is typically used for user facing services where clients from open internet can connect to servers in a data center.

## Installation

## Contributing to GRPCHub
