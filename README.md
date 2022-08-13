1. Протокол взаимодействия - HTTP
2. Хранение данных на стороне сервера: sqlite
3. Авторизация клиента - basic + https
4. сохраняемые данные - пароли (`Account`), текстовые записи (`Note`),
   данные платёжных карт (`Card`) и бинарные данные (`Binary`).
5. Шифрование данных: на стороне клиента с помощью мастер-ключа
   алгоритмом AES. 

1. Модули:
   * `internal/store/`: код, работающий с БД
   * `internal/server/`: код, работающий в http-сервере
   * `internal/client/`: код, работающий в http-клиенте
   * `cmd/server/`: код для запуска сервера
   * `cmd/client/`: код для запуска клиента
   * `cmd/server/internal/`, cmd/client/internal` - внутренние модули команд
     сервера и клиента.

Командная строка клиента:
```
go run cmd/client/main.go MODE -a ACTION flags
```
гдеs
* `MODE` - один из `user`, `cache`, `acc`, `note`, `card` или `bin`
* `ACTION`
  * для режима `user` один из `register`, `verify` или `password`
  * для режима `cache` один из `clean` или `sync`
  * для режимов `acc`, `note` или `card` - один из
    `list`, `store`, `get`, `update` или `delete`
* `flags`:
  * `-h` - получить справку по флагам
  * для режима `acc`:
    ```
    -i int
    	record ID
    -n string
    	account name
    -l string
    	account URL
    -p string
    	account password
    -u string
    	account user name
    -m string
    	account metainfo
    ```
  * для режима `note`:
    ```
    -i int
    	record ID
    -n string
    	note name
    -t string
    	note text
    -m string
    	note metainfo
    ```
  * для режима `card`:
    ```
    -i int
    	record ID
    -n string
    	card name
    -ch string
    	card holder
    -num string
    	card number
    -c string
    	card CVC code
    -em int
    	card expiry month
    -ey int
    	card expiry year
    -m string
    	card metainfo
    ```
  * для режима `bin`:
    ```
    -i int
    	binary record ID
    -n string
    	binary record name
    -f string
    	file name
    -m string
    	bin record metainfo
    ```

