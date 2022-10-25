package webui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed dist/*
