# Gendiff

Утилита командной строки на Go для сравнения конфигурационных файлов и отображения различий в различных форматах.

## Возможности

- **Множественные форматы файлов**: Поддержка JSON и YAML файлов
- **Множественные форматы вывода**: 
  - `stylish` (по умолчанию) - Человекочитаемый diff с индикаторами +/-
  - `plain` - Простые текстовые описания изменений
  - `json` - Структурированный JSON вывод
- **Рекурсивное сравнение**: Обрабатывает вложенные объекты и массивы
- **Кроссплатформенность**: Работает на Windows, macOS и Linux

## Установка

### Требования
- Go 1.21 или выше

### Сборка из исходников
```bash
git clone <repository-url>
cd go-test-project-2
make build
```

## Использование

### Базовое использование
```bash
./bin/gendiff file1.json file2.json
```

### Указание формата вывода
```bash
./bin/gendiff -f plain file1.json file2.json
./bin/gendiff --format json file1.yml file2.yml
```

### Справка
```bash
./bin/gendiff --help
```

## Примеры

### Входные файлы

**file1.json:**
```json
{
  "host": "hexlet.io",
  "timeout": 50,
  "proxy": "123.234.53.22",
  "follow": false
}
```

**file2.json:**
```json
{
  "timeout": 20,
  "verbose": true,
  "host": "hexlet.io"
}
```

### Форматы вывода

#### Stylish (по умолчанию)
```bash
./bin/gendiff file1.json file2.json
```
Вывод:
```
{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}
```

#### Plain
```bash
./bin/gendiff -f plain file1.json file2.json
```
Вывод:
```
Property 'follow' was removed
Property 'proxy' was removed
Property 'timeout' was updated. From 50 to 20
Property 'verbose' was added with value: true
```

#### JSON
```bash
./bin/gendiff -f json file1.json file2.json
```
Вывод:
```json
{
  "type": "root",
  "children": [
    {
      "type": "removed",
      "key": "follow",
      "oldValue": false
    },
    {
      "type": "unchanged",
      "key": "host",
      "value": "hexlet.io"
    },
    {
      "type": "removed",
      "key": "proxy",
      "oldValue": "123.234.53.22"
    },
    {
      "type": "updated",
      "key": "timeout",
      "oldValue": 50,
      "newValue": 20
    },
    {
      "type": "added",
      "key": "verbose",
      "newValue": true
    }
  ]
}
```

## Разработка

### Структура проекта
```
.
├── cmd/gendiff/          # CLI приложение
├── testdata/             # Тестовые фикстуры
├── gendiff.go            # Основной код библиотеки
├── gendiff_test.go       # Тесты
├── Makefile              # Команды сборки
└── go.mod                # Файл Go модуля
```

### Доступные make команды
```bash
make build              # Собрать бинарный файл
make test               # Запустить тесты
make lint               # Запустить линтер
make clean              # Очистить артефакты сборки
make setup              # Установить зависимости
```

### Запуск тестов
```bash
go test -v              # Запустить все тесты
go test -cover          # Запустить тесты с отчетом о покрытии
go test -race           # Запустить тесты с обнаружением гонок
```

## Лицензия

Этот проект распространяется под лицензией MIT.