# Basic Case

This example shows the most basic case that seems useful for using this library. The library keeps attmpting to get the lock every 1000 miliseconds, once it gets a lock it checks it has the lock every second for 10 seconds, before releasing the lock.

It uses debug level logging to show everything that goes on.
