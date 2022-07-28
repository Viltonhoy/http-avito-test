[![Go Report Card](https://goreportcard.com/badge/github.com/Viltonhoy/http-avito-test)](https://goreportcard.com/report/github.com/Viltonhoy/http-avito-test)

Тестовое задание на позицию стажер-бекенд. Микросервис для работы с балансом пользователей.
=========================================

[Ссылка на тестовое задание](https://github.com/avito-tech/autumn-2021-intern-assignment)

Данный HTTP API микросервис предназначен для работы с балансом пользователей (зачисление и списание средств, перевод средств в любой мировой валюте от пользователя к пользователю, а также для получения информации о балансе пользователя и истории всех транзакций).

### Важные заметки:
- `.env` файл был оставлен в репозитории для удобста проверки тестового задания;

## План:

- [X] Выбрать базу данных из имеющихся для хранения и работы со счетом пользователя;
- [X] Подобрать и реализовать способ работы с банковским счетом пользователя (где и каким образом хранить средства);
- [X] Реализовать базовые методы зачисления, списания и перевода средств; 
- [X] Организовать работу с удаленным API калькулятором валют;
- [X] Реализовать метод вывода истории переводов пользователя, включая идентификатор пользователя, тип перевода, сумму, дату и время, потенциального получателя и описание;   
- [X] Изучить и использовать docker и docker-compose;
- [X] Написать unit/интеграционные тесты;
- [ ] Выполнить нагрузочное тестирование микросервиса;

Для работы с базой данных была выбрана реляционная СУБД PostgreSQL. Средсвта пользователя хранятся представлены целочисленным значением (в копейках). Изначально такой формат был выбран для работы с балансом пользователя с целью избежать округления. В дальнейшем, для работы с финансами, используется тип с фиксированной точкой Decimal, что позволяет избавиться от необходимости хранить значения в целочисленном формате. (?)

## Roll-up таблица:

Для решения проблемы пересчета баланса пользователей использовалось материализованное или обычное представление, которое можно обновлять в случае необходимости. Однако при большом количестве записей скорость и эффективность данного способа снижалась, из-за постоянного пересчета всех записей из основной таблицы. На замену была выбрана roll-up таблица, которая вставляет новую строку и обновляет существующие строки в случае конфликта ограничений. Данный вариант требует меньшее количество запросов, а так же не требует постоянного пересчета записей.

posting:

| tx_id | amount | account_id | type | addressee |
|------:|:--------:|:------:|:----------:|:----------|
| 1 | 1000 | 2 | deposit | null |
| 2 | -1000 | 0 | deposit | null |
| 3 | 1000 | 1 | deposit | null |
| 4 | -1000 | 0 | deposit | null |
| 5 | -100 | 2 | withdrawal | null |
| 6 | 100 | 0 | withdrawal | null |
| 7 | -100 | 1 | transfer | 3 |
| 8 | 100 | 3 | transfer | 1 |
| 9 | -100 | 3 | withdrawal | null |
| 10 | 100 | 0 | withdrawal | null |
| 11 | -100 | 1 | withdrawal | null |
| 11 | 100 | 0 | withdrawal | null |

balances с roll-up:
 
 перед выполнением операции tx_id = 3
 | account_id | balance | last_tx_id | 
 |----:|:----:|:----------|
  | 1 |  1000 -> 800 | 3 -> 11 | 
 | 2 | 900 | 5 | 
 | 3 | 100 | 8 | 

Итог: 2 обновления

 balances без roll-up:
 
 перед выполнением операции tx_id = 3
 | account_id | balance | 
  |----:|:----------| 
  | 1 |  1000 -> 900 -> 800 | 
 | 2 | 100 -> 900 |  
 | 3 | 100 -> 0 |  

Итог: 4 обновления

Для более подробного ознакомления - https://stefan-poeltl.medium.com/views-v-s-materialized-views-v-s-rollup-tables-with-postgresql-2b3824b45330

## Тип с фиксированной точкой:

Тип с фиксированной точкой Decimal позволяет работать с дяситичными и целочисленными значениями. Данный тип был вабран для работы с финансами, чтобы избежать окргуления и потери "лишней" копейки.

#### Функционал decimal:

- Нулевое значение равно 0, и его можно безопасно использовать без инициализации;
- Сложение, вычитание, умножение без потери точности;
- Деление с заданной точностью;
- Сериализация/десериализация базы данных/sql;

Ссылка на подробную документацию типа Decimal - https://pkg.go.dev/github.com/shopspring/decimal

## Двойная запись:

Для работы с банковским счетом пользователя был выбран способ двойной бухгалтерской записи. Особенность данного способа заключается в том, что в системе с "двоичной записью" каждое значение записывается дважды - как кредит и дебет (положительное и отрицательное значение). 

#### Данная запись имеет набор правил: 

 - Каждая запись в системе должна быть сбалансированной, т.е. сумма всех значений в рамках одной операции должна давать ноль;
 - Сумма всех значений во всей системе в любой момент времени должна давать ноль (правило т.н. "пробного баланса");
 - Уже занесенные в БД значения нельзя редактировать или удалять. При необходимости исправлений операция сперва должна быть отменена другой операцией с противоположным знаком, а затем повторена с правильным значением. Это позволяет реализовать надежный аудиторский след (полный лог всех транзакций, часто требуемый при проверках);

#### Преимущество такой записи над "единичной записью":

 - Отсутствие возможности редактирования и удаления записей, что позволяет контролировать историю записей, не боясь каких либо изменений извне; 
 - Возможность построить очень комплексные системы контроля ценностями;
 - (?)

Для более подробного ознакомления предоставляю ссылки на статьи:
- https://habr.com/ru/post/480394/
- https://www.balanced.software/double-entry-bookkeeping-for-programmers/

## Запуск локально с помощью Docker:
1. Убедитесь, что у вас самые последние образы контейнеров Docker:
```
docker-compose -f ./deployments/docker-compose.yaml pull
```
2. Запустите службу в локальном Docker:
```
docker-compose -f ./deployments/docker-compose.yaml up
```
## Справочник по выполнению запросов:
1. deposit:
  - ```
  {"User_id":1, "Amount":1000}
  ``` 
2. withdrawal:

## Список вопросов и проблем:
1. Получение баланса пользователя из таблицы с двойной записью;
  - Для получения баланса решено было использовать Roll-up таблицу;
2. Генерация объемного количества данных в таблице для тестирования, sql запрос не мог генерировать большое количество данных за раз, но возвращал положительный ответ;
  - Был написан отдельный генератор значений на go; 
3. Корректная работа с docker в wsl 2. Не мог запустить docker на системе windows 10 pro;
4. Первые коммиты не имели в себе нормального описания;
  - Было принято решение не трогать старые коммиты, так как их изменения могли привести к проблемам. Последующие коммиты имеют в себе полное описание изменений в проекте;

