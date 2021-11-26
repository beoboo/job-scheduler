# TODO

* ~~GRPC~~
* ~~MTLS~~
* ~~Stream output~~
* Multithreading
    * ~~Protect scheduler jobs~~
    * ~~Protect jobs data~~
    * Improve performance on copying data

# Ideas

* Use a sync.Map to sync scheduler data
* check race conditions and deadlocks