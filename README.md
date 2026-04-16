# Gendiff

[![CI](https://github.com/krenar-rm/go-test-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/krenar-rm/go-test-project-244/actions/workflows/ci.yml)
![Coverage](.github/badges/coverage.svg)

[![asciicast](https://asciinema.org/a/lBR12XeBBYBa7MZy.svg)](https://asciinema.org/a/lBR12XeBBYBa7MZy)

Утилита для сравнения двух конфигурационных файлов.

Поддерживает JSON и YAML на вход. Форматы вывода: `stylish`, `plain`, `json`.

## Сборка

```bash
make build
```

## Использование

```bash
./bin/gendiff file1.json file2.yml
./bin/gendiff -f plain file1.json file2.json
./bin/gendiff --format json file1.yml file2.yml
```

## Разработка

```bash
make test
make lint
```
