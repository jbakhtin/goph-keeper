## Анализ и технические требования

### Сервер должен реализовывать следующую бизнес-логику:
#### 1. Регистрация, аутентификация и авторизация пользователей
##### Регистрация
Регистрация пользователя проходит по почте и паролю с подтверждением пароля;

##### Аутентификация
По запросу к серверу на аутентификацию:
1. генерируется Access Token и Refresh Token;

##### Авторизация 
По запросу к серверу на авторизацию:

```
Если Access Token не истек, то:
   Клиент успешно прошел авторизацию;
Иначе, если Refresh Token не истек, то:
   Сервер обновляет Access Token и Refresh Token и высылает клиенту;
Иначе:
   Клиент должен пройти этап Аутентификации повторно;
```

#### 2. Хранение приватных данных

Так как идея приложения заключается в том, что бы предоставить пользователю возможность обеспечить
для своих данных надежную защиту, то необходимо предоставить ему возможность шифровать добавляемые
данные вне зависимости от типа данных. Также любая метаинформация прикрепленная к основному объекту
хранения должна быть зашифрована с целью невозможности скомпрометировать пользователя, в случае, 
если злоумышленники получат доступ к базе данных. 

Необходимо не допустить попадание пользовательской информации в журнал логирования.

Передача данных на сервер будет происходить по защищенному соединению. При сохранении в базу данных
или в файловое хранилище данные будет также шифроваться. 

#### 3. Синхронизация данных между несколькими авторизованными клиентами одного владельца;
Синхронизация данных будет происходить онлайн. 
Пользователь авторизовавшийся с другого клиента сможет наблюдать свои данные на сервере.  

#### 4. Передача приватных данных владельцу по запросу.
    
  По запросу пользователя определенной информации, пользователю будет предоставлен расшифрованный результат в зависимости от типа информации.
- Если пользователь запрашивает key/value информацию, то информацию отобразится в консоли.
- Если пользователь запрашивает файл, то будет скачен файл готовый для чтения пользователем. 

### Клиент должен реализовывать следующую бизнес-логику:
#### 1. Регистрация, аутентификация и авторизация пользователей на удалённом сервере
##### Регистрация
Регистрация по почте и паролю:
  ```
  % client register
    Enter email (z-aA-Z1-9): ...
    Enter password (z-aA-Z1-9): ...
    Confirm password (z-aA-Z1-9): ...
  
  You are successful registered! But need to confirm the emil addres for sinc clients and remote servier.
 ```
Для авторизации нужно подтвердить логин.

##### Аутентификация
Успешная аутентификация:
  ```
  % client login
    Enter email (z-aA-Z1-9): ...
    Enter password (z-aA-Z1-9): ...
  
  You are successful autentificated!
 ```

Неуспешная аутентификация:
  ```
  % client login
    Enter email (z-aA-Z1-9): ...
    Enter password (z-aA-Z1-9): ...
  
  You are not autentificated! Need to confirm your email.
 ```

Успешно авторизованный клиент хранит Access Token и Refresh Token в переменных окружения.

3. Пользователь может отозвать токены авторизации.
  ```
  % client logout
  
  You are successful logout!
  ```

#### 2. Добавление и получение к приватным данным по запросу.

Данные операции может сделать только авторизованный клиент.

##### Добавление данных 

1. Данные типа ключ/значение

Cигнатура команды:
  ```
  % client kv put <key> <value> [arguments]
  ```
Пример использования:
  ```
  % client kv put john@gmail.com 1234 -description="My secret credentials to ThisWebSite.com" -expire="1d"
  
    ========= Key/Value =========
    key                value
    ---                -----
    created_at         2023-10-20 17:19:00
    expire_at          2023-10-21 17:19:00
    description        "My secret credentials to ThisWebSite.com"
  ```

2. Данные типа банковские карты:
Cигнатура команды:
  ```
  % client bc put <number> [arguments]
  
  <number> - card number
  
  arguments:
  
  -code - CSV code on back side of card 
  -vfd - validity from date
  -ed - expire date
  -chn - card holder name
  -description - Description of the bank card 
  -sc - sort code
  -an - account number
  -pnl - payment network logo
  -bl - banking logo
  -expire - date the card was removed from the database
  
  ```
Пример использования:
  ```
  % client bc put "1111 1111 1111 1111" -code=123 -description="My most important card"
  
    ========= Bank Card =========
    Number             number
    ---                -----
    created_at         2023-10-20 17:19:00
    expire             nil
    description        My most important card
    code               123
    vfd                <nil>
    ed                 <nil>
    chn                <nil>
    sc                 <nil>
    an                 <nil>
    pnl                <nil>
    bl                 <nil>
  ```

3. Данные типа файл:
   Cигнатура команды:
  ```
  % client file put <file path> [arguments]
  
  <file path> - path to file on your local machine
  
  arguments:
  
  -description - Description of the bank card 
  -expire - date the card was removed from the database
  
  ```
Пример использования:
  ```
  % client file put /Users/root/.ssh/pets/id_rsa -description="My ssh "
  
    ========= File =========
    File               fileName
    ---                -----
    created_at         2023-10-20 17:19:00
    expire             nil
    description        My most important card
  ```

##### Получение данных

1. Данные типа ключ/значение
// TODO: скопировано из предыдущего пункта - переделать
Cигнатура команды:
  ```
  % client kv put <key> <value> [arguments]
  ```
Пример использования:
  ```
  % client kv put john@gmail.com 1234 -description="My secret credentials to ThisWebSite.com" -expire="1d"
  
    ========= Key/Value =========
    key                value
    ---                -----
    created_at         2023-10-20 17:19:00
    expire_at          2023-10-21 17:19:00
    description        "My secret credentials to ThisWebSite.com"
  ```

2. Данные типа банковские карты:
   Cигнатура команды:
  ```
  % client bc put <number> [arguments]
  
  <number> - card number
  
  arguments:
  
  -code - CSV code on back side of card 
  -vfd - validity from date
  -ed - expire date
  -chn - card holder name
  -description - Description of the bank card 
  -sc - sort code
  -an - account number
  -pnl - payment network logo
  -bl - banking logo
  -expire - date the card was removed from the database
  
  ```
Пример использования:
  ```
  % client bc put "1111 1111 1111 1111" -code=123 -description="My most important card"
  
    ========= Bank Card =========
    Number             number
    ---                -----
    created_at         2023-10-20 17:19:00
    expire             nil
    description        My most important card
    code               123
    vfd                <nil>
    ed                 <nil>
    chn                <nil>
    sc                 <nil>
    an                 <nil>
    pnl                <nil>
    bl                 <nil>
  ```

3. Данные типа файл:
   Cигнатура команды:
  ```
  % client file put <file path> [arguments]
  
  <file path> - path to file on your local machine
  
  arguments:
  
  -description - Description of the bank card 
  -expire - date the card was removed from the database
  
  ```
Пример использования:
  ```
  % client file put /Users/root/.ssh/pets/id_rsa -description="My ssh "
  
    ========= File =========
    File               fileName
    ---                -----
    created_at         2023-10-20 17:19:00
    expire             nil
    description        My most important card
  ```

### Функции, реализация которых остаётся на усмотрение исполнителя:
- создание, редактирование и удаление данных на стороне сервера или клиента;
- формат регистрации нового пользователя;
- выбор хранилища и формат хранения данных;
- обеспечение безопасности передачи и хранения данных;
- протокол взаимодействия клиента и сервера;
- механизмы аутентификации пользователя и авторизации доступа к информации.

### Дополнительные требования:
- клиент должен распространяться в виде CLI-приложения с возможностью запуска на платформах Windows, Linux и Mac OS;
- клиент должен давать пользователю возможность получить информацию о версии и дате сборки бинарного файла клиента;

- придерживаться принципов чистой архитектуры;
- реализовать Graceful Shutdown для клиента и сервера;
- README.md должен описывать инструкции для локального запуска проекта, тестирования и запуска линтера;
- схема БД должна быть описана с помощью механизма миграций;
- добавить CodeCov Report в PR;
- В CI должен быть добавлен линтер https://github.com/golangci/golangci-lint;
- Нужно валидировать входные параметры в запросах. Если используются proto файлы, то с помощью https://github.com/bufbuild/protovalidate-go;
- Добавить docker-compose для сервиса.


Нужно учитывать, что текстовые и бинарные данные могут быть любого объема.

На основе ТЗ выделены основные типы данных:
- текстовые данные:
    - string
    - file:
        - .txt
        - .json
        - .csv
        - .xml
        - .html
        - .sh
        - ...
- бинарные данные:
    - .exe
    - .bin
    - .jpg
    - .png
    - .zip
    - ...