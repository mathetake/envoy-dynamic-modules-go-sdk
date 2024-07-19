# EnvoyX Go SDK

This is the Go SDK for the EnvoyX modules. The modules are shared libraries that can be loaded by the EnvoyX proxy to extend HTTP filtering capabilities.

The shared library must be compiled with the same environment as EnvoyX, that means the programs must be compiled 
on amd64 Linux with the same version of glibc as the EnvoyX proxy.
