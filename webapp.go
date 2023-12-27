package main

import (
	"embed"
)

//go:embed webapp/*
var webapp embed.FS
