This was vendored to make godep work again.

`github.com/StackExchange/wmi` has no buildable go source files on anything but windows.

This breaks running `godep save ./...` an anything but windows :(

The vendored version is https://github.com/StackExchange/wmi/commit/f3e2bae1e0cb5aef83e319133eabfee30013a4a5#diff-d5a622cd1a8b29bf64376f30672f8eec