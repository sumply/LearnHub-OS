// Типы для учебных материалов
export type MaterialType = 'video' | 'pdf' | 'audio' | 'presentation' | 'document' | 'image' | 'archive' | 'link';

export interface Material {
  id: string;
  title: string;
  description: string;
  type: MaterialType;
  url: string;
  fileName?: string; // Имя файла для скачивания
  fileSize?: number; // Размер файла в байтах
  category: string; // id предмета
  grade?: string; // Класс (например, "11А", "10Б")
  uploadDate: Date;
  uploadedBy: string; // ID пользователя, загрузившего материал
  visible: boolean;
  tags?: string[]; // Теги для поиска
}

// Функция для определения типа файла по расширению
export function getMaterialType(fileName: string): MaterialType {
  const extension = fileName.toLowerCase().split('.').pop();
  switch (extension) {
    case 'pdf':
      return 'pdf';
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'wmv':
      return 'video';
    case 'mp3':
    case 'wav':
    case 'ogg':
      return 'audio';
    case 'ppt':
    case 'pptx':
      return 'presentation';
    case 'doc':
    case 'docx':
    case 'txt':
      return 'document';
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
      return 'image';
    case 'zip':
    case 'rar':
    case '7z':
      return 'archive';
    default:
      return 'document';
  }
}

// Функция для форматирования размера файла
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Б';
  const k = 1024;
  const sizes = ['Б', 'КБ', 'МБ', 'ГБ'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Тестовые материалы с файлами из папки materials
export const materials: Material[] = [
  {
    id: 'mat1',
    title: 'Создание запросов',
    description: 'Учебный материал по созданию SQL запросов',
    type: 'pdf',
    url: '/materials/#2 Создание запросов.pdf',
    fileName: '#2 Создание запросов.pdf',
    fileSize: 508 * 1024, // 508KB
    category: 'db',
    grade: '11А',
    uploadDate: new Date('2024-01-15'),
    uploadedBy: 'teacher1',
    visible: true,
    tags: ['SQL', 'запросы', 'база данных']
  },
  {
    id: 'mat2',
    title: 'Межтабличные связи',
    description: 'Изучение связей между таблицами в базах данных',
    type: 'pdf',
    url: '/materials/#3 Межтабличные связи.pdf',
    fileName: '#3 Межтабличные связи.pdf',
    fileSize: 480 * 1024,
    category: 'db',
    grade: '11А',
    uploadDate: new Date('2024-01-20'),
    uploadedBy: 'teacher1',
    visible: true,
    tags: ['связи', 'таблицы', 'нормализация']
  },
  {
    id: 'mat3',
    title: 'Основы работы в ОС QNX Neutrino',
    description: 'Лабораторная работа по операционной системе QNX',
    type: 'pdf',
    url: '/materials/1_Лабораторная_работа_Основы_работы_в_ОС_QNX_Neutrino_1.pdf',
    fileName: '1_Лабораторная_работа_Основы_работы_в_ОС_QNX_Neutrino_1.pdf',
    fileSize: 3.8 * 1024 * 1024, // 3.8MB
    category: 'os',
    grade: '11А',
    uploadDate: new Date('2024-02-01'),
    uploadedBy: 'teacher2',
    visible: true,
    tags: ['QNX', 'операционная система', 'лабораторная']
  },
  {
    id: 'mat4',
    title: 'Механизмы межпроцессного взаимодействия',
    description: 'Изучение IPC в операционных системах',
    type: 'pdf',
    url: '/materials/4_Лабораторная_работа_Механизмы_межпроцессного_взаимодействия_и.pdf',
    fileName: '4_Лабораторная_работа_Механизмы_межпроцессного_взаимодействия_и.pdf',
    fileSize: 1021 * 1024,
    category: 'os',
    grade: '11А',
    uploadDate: new Date('2024-02-05'),
    uploadedBy: 'teacher2',
    visible: true,
    tags: ['IPC', 'процессы', 'взаимодействие']
  },
  {
    id: 'mat5',
    title: 'Командные сценарии',
    description: 'Лабораторная работа по написанию скриптов',
    type: 'pdf',
    url: '/materials/2_Лабораторная_работа_Командные_сценарии_3.pdf',
    fileName: '2_Лабораторная_работа_Командные_сценарии_3.pdf',
    fileSize: 768 * 1024,
    category: 'os',
    grade: '11А',
    uploadDate: new Date('2024-02-10'),
    uploadedBy: 'teacher2',
    visible: true,
    tags: ['скрипты', 'команды', 'автоматизация']
  },
  {
    id: 'mat6',
    title: 'Управление процессами и потоками в ОС QNX',
    description: 'Продвинутая лабораторная работа по QNX',
    type: 'pdf',
    url: '/materials/3_Лабораторная_работа_Управление_процессами_и_потоками_в_ОС_QNX.pdf',
    fileName: '3_Лабораторная_работа_Управление_процессами_и_потоками_в_ОС_QNX.pdf',
    fileSize: 922 * 1024,
    category: 'os',
    grade: '11А',
    uploadDate: new Date('2024-02-15'),
    uploadedBy: 'teacher2',
    visible: true,
    tags: ['процессы', 'потоки', 'QNX', 'управление']
  },
  {
    id: 'mat7',
    title: 'Передача сообщений в ОС',
    description: 'Исследование механизмов передачи сообщений',
    type: 'pdf',
    url: '/materials/5_Лабораторная_работа_Исследование_работы_механизмов_передачи_сообщений.pdf',
    fileName: '5_Лабораторная_работа_Исследование_работы_механизмов_передачи_сообщений.pdf',
    fileSize: 968 * 1024,
    category: 'os',
    grade: '11А',
    uploadDate: new Date('2024-02-20'),
    uploadedBy: 'teacher2',
    visible: true,
    tags: ['сообщения', 'механизмы', 'коммуникация']
  },
  {
    id: 'mat8',
    title: 'Разработка базы данных в нотации IDEF1X',
    description: 'Методология проектирования баз данных',
    type: 'pdf',
    url: '/materials/#6_Разработка_базы_данных_в_нотации_IDEF1X.pdf',
    fileName: '#6_Разработка_базы_данных_в_нотации_IDEF1X.pdf',
    fileSize: 976 * 1024,
    category: 'db',
    grade: '11А',
    uploadDate: new Date('2024-02-25'),
    uploadedBy: 'teacher1',
    visible: true,
    tags: ['IDEF1X', 'проектирование', 'моделирование']
  },
  {
    id: 'mat9',
    title: 'Настройка сетевого оборудования',
    description: 'Конфигурация коммутаторов и маршрутизаторов',
    type: 'pdf',
    url: '/materials/02_Configuring_the_Initial_Settings_of_the_Switch_And_Implement.pdf',
    fileName: '02_Configuring_the_Initial_Settings_of_the_Switch_And_Implement.pdf',
    fileSize: 641 * 1024,
    category: 'networks',
    grade: '11А',
    uploadDate: new Date('2024-03-01'),
    uploadedBy: 'teacher3',
    visible: true,
    tags: ['сети', 'коммутаторы', 'конфигурация']
  },
  {
    id: 'mat10',
    title: 'Английский учебник',
    description: 'Учебное пособие по английскому языку',
    type: 'pdf',
    url: '/materials/angljskij_uchebnik.pdf',
    fileName: 'angljskij_uchebnik.pdf',
    fileSize: 9.0 * 1024 * 1024, // 9MB
    category: 'eng',
    grade: '11А',
    uploadDate: new Date('2024-03-05'),
    uploadedBy: 'teacher4',
    visible: true,
    tags: ['английский', 'учебник', 'грамматика']
  },
  {
    id: 'mat11',
    title: 'Книга игрока 2024',
    description: 'Справочник по настольным играм',
    type: 'pdf',
    url: '/materials/Kniga_Igroka_2024.pdf',
    fileName: 'Kniga_Igroka_2024.pdf',
    fileSize: 20 * 1024 * 1024, // 20MB
    category: 'games',
    grade: '11А',
    uploadDate: new Date('2024-03-10'),
    uploadedBy: 'teacher5',
    visible: true,
    tags: ['игры', 'справочник', 'развлечения']
  }
]; 