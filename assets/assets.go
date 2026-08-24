package assets

import "embed"

//go:embed  images/*
//go:embed  images/*/*
//go:embed  images/*/*/*/*
var images embed.FS

//go:embed maps/*
var maps embed.FS

//go:embed sfx/*
var sfx embed.FS

var Assets = struct {
	Images embed.FS
	Maps   embed.FS
	Sfx    embed.FS
}{
	Images: images,
	Maps:   maps,
	Sfx:    sfx,
}
