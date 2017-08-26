# Rest API for Go-interface with access to tarantool

* **GET**
    * /users -- получение списка пользователей
        * curl -X GET -i http://localhost:9090/users
    * /spaces -- получение списка пространств пользователя
        * curl -X GET -H "Authorization: Bearer token" -i http://localhost:9090/spaces
    * /history -- получение истории пользователя
        * curl -X GET -H "Authorization: Bearer token" -i http://localhost:9090/history
    * /space/{name_space}/history -- получение истории работы с пространством
        * curl -X GET -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_spaces}/history
    * /spaces/{name_space}/permissions -- получение списка пользователей, имеющих доступ к данному пространству
        * curl -X GET -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}/permissions
    * /spaces/{name_space}/tuples -- получение списка кортежей в данном пространстве
        * curl -X GET -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}/tuples
    * /spaces/{name_space}/tuples/{id_tuple} -- получение конкретного кортежа
        * curl -X GET -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}/tuples/{id_tuple}

* **POST**
    * /users -- регистрация нового пользователя и/или выдача ему токена для входа
        * curl -X POST -i http://localhost:9090/users -d "user=test&password=12345"
    * /users/{name}/spaces/{name_space}/permissions -- добавление прав на чтение и редактирование пространства
        * curl -X POST -H "Authorization: Bearer token" -i http://localhost:9090/users/{name}/spaces/{name_space}/permissions
    * /spaces/{name_space} -- добавление нового пространства
        * curl -X POST -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}
    * /spaces/{name_space}/tuples/{id_tuple} -- добавление нового кортежа
        * curl -X POST -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}/tuples/{id_tuple}

* **PUT**
    * /spaces/{name_space}/tuples/{id_tuple} -- обновление кортежа
        * curl -X PUT -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}/tuples/id_tuple}
        
* **DELETE**
    * /spaces/{name_space}/tuples/{id_tuple} -- удаление кортежа
        * curl -X DELETE -H "Authorization: Bearer token" -i http://localhost:9090/spaces/{name_space}/tuples/{id_tuple}