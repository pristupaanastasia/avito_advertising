# API Avito advertising
## Запуск
### Создание таблицы

sudo docker-compose run db bash

psql --host=db --username=postgres

pass: 1805

\connect avito_db

### Запуск докера

docker-compose build

docker-compose up
## Описание методов ##

### IndexHandler ###

#### Метод получения списка объявлений ####
Реализовано:
+ Пагинация: на одной странице 10 объявлений;
+ Cортировки: по цене (возрастание/убывание) и по дате создания (возрастание/убывание). По умолчанию стоит
  сортировка по возрастанию цены. Для выбора других сортировок в параметры передать "sort=update_desc",
  "sort=price", "sort=price_desc".
+ Поля в ответе: название объявления, ссылка на главное фото, цена.
  Ответ выводится в виде JSON.
### FindHandler ###

#### Метод получения конкретного объявления ####

+ Поля в ответе: название объявления, цена, ссылка на главное фото;
+ Передав ключ в параметры url "fields": описание, ссылки на все фото.
### CreateHandler ###

#### Метод создания объявления ####
+ Принимает значения: название, описание, ссылки на фото (не более 3), цена;
+ Возвращает ID созданного объявления и код результата (ошибка или успех).

## Примеры

curl -X GET "http://localhost:9000/ad?page=1&sort=price"

[{"id":8,"price":12,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":1,"price":23,"name":"синий","image":["https://images.app.goo.gl/yoP6Yc7iPZQsGA858"],"update":"2021-02-03"},{"id":9,"price":88,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":10,"price":99,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":11,"price":100,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":7,"price":234,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":5,"price":234,"name":"","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":4,"price":234,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":3,"price":234,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"},{"id":2,"price":234,"name":"wolf","image":["https://images.app.goo.gl/UKvzedV5obN8y5tC7"],"update":"2021-02-04"}]%

###find
curl -X GET "http://localhost:9000/find?id=1&fields"

{"id":1,"price":23,"name":"синий","description":"красивый синий","image":["https://images.app.goo.gl/yoP6Yc7iPZQsGA858","https://images.app.goo.gl/BVk7Nho7LeRDiJcV8"],"update":"2021-02-03"}

curl -X GET "http://localhost:9000/find?id=1"

{"id":1,"price":23,"name":"синий","image":["https://images.app.goo.gl/yoP6Yc7iPZQsGA858"],"update":"2021-02-03"}

curl -X GET "http://localhost:9000/find?id=2"
{}
###create
curl -X GET "http://localhost:9000/create?price=234&name=wolf&description=nice_wolf&image=https://images.app.goo.gl/UKvzedV5obN8y5tC7&image=https://images.app.goo.gl/mcXR7BVDdccmFVWS6"

{"id":2,"status":201}

curl -X GET "http://localhost:9000/create?price=234&name=wolf&description=nice_wolf&image=https://images.app.goo.gl/UKvzedV5obN8y5tC7&image=https://images.app.goo.gl/mcXR7BVDdccmFVWS6&image=https://images.app.goo.gl/UKvzedV5obN8y5tC7&image=https://images.app.goo.gl/UKvzedV5obN8y5tC7"

{"Err":{},"Code":400,"Name":"Количество ссылок не должно превышать 3"}

