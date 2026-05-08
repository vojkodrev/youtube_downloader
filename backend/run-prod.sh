#!/bin/bash
go build -o backend . && GIN_MODE=release ./backend
