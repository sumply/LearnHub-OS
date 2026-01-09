package main

import (
	"server/internal/dto"
)

func createStudents() []*dto.UserCreateReq {
	return []*dto.UserCreateReq{
		{
			FirstName:  "Игорь",
			LastName:   "Лебедев",
			MiddleName: "Сергеевич",
			Email:      "igor6@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Наталья",
			LastName:   "Михайлова",
			MiddleName: "Владимировна",
			Email:      "natalya7@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Константин",
			LastName:   "Григорьев",
			MiddleName: "Иванович",
			Email:      "konstantin8@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Светлана",
			LastName:   "Орлова",
			MiddleName: "Павловна",
			Email:      "svetlana9@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Даниил",
			LastName:   "Никитин",
			MiddleName: "Александрович",
			Email:      "danil10@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Ксения",
			LastName:   "Фролова",
			MiddleName: "Михайловна",
			Email:      "kseniya11@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Максим",
			LastName:   "Соколов",
			MiddleName: "Игоревич",
			Email:      "maxim12@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Александра",
			LastName:   "Воронова",
			MiddleName: "Сергеевна",
			Email:      "alexandra13@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Евгений",
			LastName:   "Тихонов",
			MiddleName: "Павлович",
			Email:      "evgeny14@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Мария",
			LastName:   "Крылова",
			MiddleName: "Александровна",
			Email:      "maria15@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Вася",
			LastName:   "Пупкин",
			MiddleName: "Пупкинович",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Петя",
			LastName:   "Иванов",
			MiddleName: "Иванович",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Маша",
			LastName:   "Сидорова",
			MiddleName: "Сидоровна",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Дима",
			LastName:   "Кузнецов",
			MiddleName: "Кузнецович",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       0,
		},
		{
			FirstName:  "Аня",
			LastName:   "Фёдорова",
			MiddleName: "Фёдоровна",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       0,
		},
	}
}

func createTeachers() []*dto.UserCreateReq {
	return []*dto.UserCreateReq{
		{
			FirstName:  "Сергей",
			LastName:   "Петров",
			MiddleName: "Алексеевич",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       1,
		},
		{
			FirstName:  "Елена",
			LastName:   "Васильева",
			MiddleName: "Игоревна",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       1,
		},
		{
			FirstName:  "Иван",
			LastName:   "Смирнов",
			MiddleName: "Павлович",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       1,
		},
		{
			FirstName:  "Ольга",
			LastName:   "Кузьмина",
			MiddleName: "Николаевна",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       1,
		},
		{
			FirstName:  "Александр",
			LastName:   "Морозов",
			MiddleName: "Владимирович",
			Email:      "yanushkevichvladislav@mail.ru",
			Role:       1,
		},
	}
}

func createGroups() []*dto.SubjectCreateReq {
	return []*dto.SubjectCreateReq{
		{
			Name: "Русский язык",
		},
		{
			Name: "Математика",
		},
		{
			Name: "Физика",
		},
		{
			Name: "История",
		},
		{
			Name: "Биология",
		},
		{
			Name: "Информатика",
		},
	}
}
