import { createSignal, For, createResource } from "solid-js";
import { getSubjects } from '../services/subjectGroupService';
import type { MaterialType } from '../config/materials';

const materialTypes: { value: MaterialType, label: string }[] = [
  { value: 'pdf', label: 'PDF' },
  { value: 'video', label: 'Видео' },
  { value: 'audio', label: 'Аудио' },
  { value: 'presentation', label: 'Презентация' },
  { value: 'document', label: 'Документ' },
  { value: 'image', label: 'Изображение' },
  { value: 'archive', label: 'Архив' },
  { value: 'link', label: 'Ссылка' },
];

const grades = ["11А", "11Б", "10А", "10Б", "9А", "9Б"];

interface AddMaterialFormProps {
  onMaterialAdded?: (data: {
    title: string;
    description: string;
    type: MaterialType;
    url: string;
    fileName: string;
    fileSize: number;
    category: string;
    grade?: string;
    visible: boolean;
    tags?: string[];
  }) => void;
}

const AddMaterialForm = (props: AddMaterialFormProps) => {
  // Загружаем предметы через API
  const [subjectsData] = createResource(getSubjects);
  
  const [title, setTitle] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [file, setFile] = createSignal<File | null>(null);
  const [category, setCategory] = createSignal("");
  const [grade, setGrade] = createSignal(grades[0]);
  const [type, setType] = createSignal<MaterialType>('pdf');
  const [visible, setVisible] = createSignal(true);
  const [tags, setTags] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  
  // Устанавливаем первый предмет по умолчанию когда данные загрузятся
  createResource(() => {
    const subjects = subjectsData();
    if (subjects && subjects.length > 0 && !category()) {
      setCategory(subjects[0].id.toString());
    }
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    if (!title() || !file()) {
      setError("Пожалуйста, заполните все обязательные поля и выберите файл.");
      return;
    }
    setLoading(true);
    try {
      const url = URL.createObjectURL(file()!);
      props.onMaterialAdded && props.onMaterialAdded({
        title: title(),
        description: description(),
        type: type(),
        url,
        fileName: file()!.name,
        fileSize: file()!.size,
        category: category(),
        grade: grade(),
        visible: visible(),
        tags: tags().split(",").map(t => t.trim()).filter(Boolean)
      });
      setTitle("");
      setDescription("");
      setFile(null);
      setTags("");
    } catch (e) {
      setError("Ошибка при добавлении материала. Попробуйте еще раз.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ "margin-bottom": "24px" }}>
      <h3>Добавить учебный материал</h3>
      <div>
        <label>
          Название*:<br />
          <input
            type="text"
            value={title()}
            onInput={e => setTitle((e.target as HTMLInputElement).value)}
            required
            style={{ width: "100%", "margin-bottom": "8px" }}
          />
        </label>
      </div>
      <div>
        <label>
          Описание:<br />
          <textarea
            value={description()}
            onInput={e => setDescription((e.target as HTMLTextAreaElement).value)}
            style={{ width: "100%", "margin-bottom": "8px" }}
          />
        </label>
      </div>
      <div>
        <label>
          Предмет*:<br />
          <select value={category()} onInput={e => setCategory((e.target as HTMLSelectElement).value)} style={{ width: "100%", "margin-bottom": "8px" }}>
            <For each={subjectsData() || []}>{subject => <option value={subject.id.toString()}>{subject.name}</option>}</For>
          </select>
        </label>
      </div>
      <div>
        <label>
          Класс:<br />
          <select value={grade()} onInput={e => setGrade((e.target as HTMLSelectElement).value)} style={{ width: "100%", "margin-bottom": "8px" }}>
            <For each={grades}>{g => <option value={g}>{g}</option>}</For>
          </select>
        </label>
      </div>
      <div>
        <label>
          Тип материала:<br />
          <select value={type()} onInput={e => setType((e.target as HTMLSelectElement).value as MaterialType)} style={{ width: "100%", "margin-bottom": "8px" }}>
            <For each={materialTypes}>{t => <option value={t.value}>{t.label}</option>}</For>
          </select>
        </label>
      </div>
      <div>
        <label>
          Файл*:<br />
          <input
            type="file"
            onInput={e => setFile((e.target as HTMLInputElement).files?.[0] || null)}
            required
            style={{ "margin-bottom": "8px" }}
          />
        </label>
      </div>
      <div>
        <label>
          Теги (через запятую):<br />
          <input
            type="text"
            value={tags()}
            onInput={e => setTags((e.target as HTMLInputElement).value)}
            style={{ width: "100%", "margin-bottom": "8px" }}
          />
        </label>
      </div>
      <div>
        <label>
          <input type="checkbox" checked={visible()} onInput={e => setVisible((e.target as HTMLInputElement).checked)} /> Видимый для всех
        </label>
      </div>
      {error() && <div style={{ color: "red", "margin-bottom": "8px" }}>{error()}</div>}
      <button type="submit" disabled={loading()}>
        {loading() ? "Загрузка..." : "Добавить материал"}
      </button>
    </form>
  );
};

export default AddMaterialForm; 