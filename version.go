package main

//needs to be a var (no const)
//so that we van overwrite during linking with -X main.GITCOMMIT ...
var (
	Version   = "0.1.0-dev"
	GITCOMMIT = "HEAD"
)
