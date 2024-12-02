Привет) Приложение по REST принимает запрос вида
POST api/v1/wallet
{
valletId: UUID,
operationType: DEPOSIT or WITHDRAW,
amount: 1000
}
После выполнять логику по изменению счета в базе данных
Также есть возможность получить баланс кошелька
GET api/v1/wallets/{UUID}
Для создания кошелька используйте запрос в формате :
curl -X 'POST' \                                   
'http://localhost:8080/api/v1/wallets' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{"initialBalance": 1000}'

docker-compose запускается стандартной комантдой: docker-compose up --build

Swagger документация реализованна http://127.0.0.1:8080/swagger/index.html

Чтобы запустить без docker фаила поменяйте в .env поле DB_HOST значение на localhost
запустите командой go run ./cmd   