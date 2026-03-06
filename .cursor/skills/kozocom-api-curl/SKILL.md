---
name: kozocom-api-curl
description: Guide for making API requests to kozocom portal using curl with authentication token from kozocom-tui. Use when testing API endpoints, debugging API calls, or making manual requests to the kozocom portal API.
---

# Kozocom API Curl Guide

Hướng dẫn sử dụng curl với token từ kozocom-tui để gọi API.

## Token Location

Token được lưu tại đường dẫn:

```
$HOME/.config/kozocom-tui/token
```

Hoặc trên macOS:

```
$HOME/Library/Application Support/kozocom-tui/token
```

Trích xuất token để dùng với curl:

```bash
TOKEN=$(cat ~/.config/kozocom-tui/token)
```

## Base URL

```
http://localhost:8000/api/v1
```

## Common curl commands

### GET request

```bash
TOKEN=$(cat ~/.config/kozocom-tui/token)
curl -H "Authorization: Bearer $TOKEN" \
     -H "Accept: application/json" \
     http://localhost:8000/api/v1/user/profile
```

### POST request with JSON body

```bash
TOKEN=$(cat ~/.config/kozocom-tui/token)
curl -X POST \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d '{"key": "value"}' \
     http://localhost:8000/api/v1/endpoint
```

### POST with data from file

```bash
TOKEN=$(cat ~/.config/kozocom-tui/token)
curl -X POST \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d @request.json \
     http://localhost:8000/api/v1/endpoint
```

### Pretty print JSON response

```bash
TOKEN=$(cat ~/.config/kozocom-tui/token)
curl -s -H "Authorization: Bearer $TOKEN" \
     http://localhost:8000/api/v1/endpoint | jq .
```

## Quick alias

Thêm vào `.bashrc` hoặc `.zshrc`:

```bash
alias kozocurl='curl -H "Authorization: Bearer $(cat ~/.config/kozocom-tui/token)" -H "Accept: application/json"'
```

Sau đó dùng:

```bash
kozocurl http://localhost:8000/api/v1/user/profile
```
