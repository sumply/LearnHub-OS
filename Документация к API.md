## Перечисления:

#### Роль пользователя (USER_ROLE):

1. Студент
2. Учитель
3. Админ
4. Суперпользователь

#### Статус прохождения квиза (PROGRESS_STATUS):

0. К прохождению не приступили
1. В процессе прохождения
2. Ожидает проверки
3. Завершен

#### Статус ответа на вопрос при прохождении (ANSWER_STATUS):

0. Ответ не дан
1. Ожидает проверки
2. Неверный
3. Верный

## Эндпоинты и форматы ответа:

***

### POST /login

**Авторизация пользователя.**

Запрос:

```json
{
    "login": "<your_login>",
    "password": "<your_password>"
}
```

Ответ:

```json
{
    "refresh_token": "<jwt_token>",
    "access_token": "<jwt_token>"
}
```


***

### POST /users

**Создание пользователя (логин и пароль генерирует сервер).**

Запрос:

```json
{
    "first_name": "<first_name>",
    "last_name": "<last_name>",
    "middle_name": "<middle_name> - необязательно",
    "email": "email@email.ru",
    "role": USER_ROLE
}
```

***

### GET /users

**Получение списка пользователей.**

Ответ:

```json
[
{
    "id": 0,
    "short_name": "<last_name F.M.>"
}
]
```

***

### GET /users/me

**Получение информации о себе (авторизация происходит через токен).**

Ответ:

```json
{
    "id": 0,
    "first_name": "<first_name>",
    "last_name": "<last_name>",
    "middle_name": "<middle_name>"
}
```

***

### GET /users/{user_id}

**Получение расширенной информации о пользователе.**

Ответ:

```json
{
    "id": 0,
    "first_name": "<first_name>",
    "last_name": "<last_name>",
    "middle_name": "<middle_name>"
}
```

***

## POST /groups

**Создание нового класса**

Запрос:

```json
{
    "name": "<group_name>",
    "curator_id": 0,
}
```

***

### GET /groups

**Получение списка всех классов**

Ответ:

```json
[
{
    "id": 0,
    "name": "<group_name>",
    curator: {
        "id": 0,
        "short_name": "<last_name F.M.>"
    }
}
]
```

***

### POST /groups/{group_id}/students

**Добавление к классу учеников**

Запрос:

```json
{
    "student_ids": [0, 1, 2]
}
```

***

### POST /subject

**Создание нового предмета**

Запрос:

```json
{
    "name": "<subject_name>"
}
```

***

### GET /subject

**Получение списка всех предметов**

Ответ:

```json
[
{
    "id": 0,
    "name": "<subject_name>"
}
]
```

***

### POST /quizzes

**Создание квиза**

Запрос:

```json
{
    "title": "<quiz_title>",
    "summary": "<quiz_summary>",
    "subject_id": 0,
    "group_ids": [0, 1],
    "questions": [
        {
            "text": "<question_text>",
            "options": [
                {
                    "text": "<option_text>",
                    "is_correct": true
                },
                {
                    "text": "<option_text>",
                    "is_correct": false
                }
            ]
        }
    ]
}
```

*Уточнение 1: если массив "options" будет пустым, правильность вопроса проверяет создатель квиза.*

*Уточнение 2: пользователь, под которым создается этот ресурс, автоматически становится его владельцем.*

***

### GET /quizzes

**Получение списка квизов.**

Ответ:

```json
[
    {
        "id": 0,
        "title": "<quiz_title>",
        "summary": "<quiz_summary>",
        "owner": {
            "id": 0,
            "short_name": "<last_name F.M.>"
        },
        "subject": {
            "id": 0,
            "name": "<subject_name>"
        },
        "group": [
            {
                "id": 0,
                "name": "<group_name>"
            }
        ]
    }
]
```

***

### DELETE /quizzes/{quiz_id}

***

### GET /progress

Ответ:

```json
[
    {
        "id": 0,
        "quiz": {
            "id": 0,
            "title": "<quiz_title>",
            "summary": "<quiz_summary>",
            "owner": {
                "id": 0,
                "short_name": "<last_name F.M.>"
            },
            "subject": {
                "id": 0,
                "name": "<subject_name>"
            }
        },
        "user": {
            "id": 0,
            "short_name": "<last_name F.M.>"
        },
        "status": PROGRESS_STATUS,
        "score": 0,
        "start_date": "2025-01-01T10:00:00Z - необязательно",
        "completed_date": "2025-01-01T10:30:00Z - необязательно"
    }
]
```

*Уточнение 1: прогресс по квизу генерируется автоматически при создания квиза.*

*Уточнение 2: поле "user" - аккаунт ученика, к которому приклеплен прогресс.*

***

### POST /progress/{progress_id}/start

**Начало прохождения квиза**

*Уточнение: начать проходить квиз может только пользователь, к которому этот прогресс прикреплен*

***

### POST /progress/{progress_id}/finish

**Завершение прохождения квиза**

*Уточнение 1: закончить проходить квиз может только пользователь, к которому этот прогресс прикреплен.*

*Уточнение 2: если квиз подразумевает ручную проверку ответа, статус квиза становится (Ожидает проверки).*

*Уточнение 3: при статусе (Ожидает проверки), окончательно завершить прохождение квиза может его создатель*

***

### PATCH /progress/{progress_id}/answer/{answer_id}

**Записывание ответа на вопрос квиза*

Запрос:

```json
{
    "text": "<answer_text>"
}
```

*Уточнение 1: поле "text" должно быть равным одному из вариантов ответа*

*Уточнение 2: если вариантов ответа нет, можно вписывать, что угодно*

***

### POST /progress/{progress_id}/answer/{answer_id}/correct

**Установить статус ответа на вопрос как "верный"**

*Уточнение 1: проверить возможно только тот вопрос, который ожидает проверки*

*Уточенние 2: проверяет ответ создатель квиза*

***

### POST /progress/{progress_id}/answer/{answer_id}/incorrect

**Установить статус ответа на вопрос как "неверный"**

*Уточнение 1: проверить возможно только тот вопрос, который ожидает проверки*

*Уточенние 2: проверяет ответ создатель квиза*