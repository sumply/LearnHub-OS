import { type Component, createSignal, For, Show, createResource } from 'solid-js';
import { getCurrentUser } from '../utils/api';
import { getAllActivities, type TaskUnion, addActivity, updateActivity, removeActivity } from '../utils/activitiesService';
import Header from '../components/Header';
import { getSubjects, getGroups } from '../services/subjectGroupService';
import { getAllUsers } from '../services/userService';
import { getAuthToken } from '../utils/api';
import { createQuiz } from '../services/quizService';
import { getProgress } from '../services/progressService';
import * as apiClient from '../utils/apiClient';
import type { MaterialAttachment } from '../utils/api';

const Tasks: Component = () => {
  const [selectedTask, setSelectedTask] = createSignal<TaskUnion | null>(null);
  const [subjectFilter, setSubjectFilter] = createSignal('all');
  const [typeFilter, setTypeFilter] = createSignal('all');
  const [teacherFilter, setTeacherFilter] = createSignal('all');
  // Загружаем квизы через GET /progress
  const [progressData, { refetch: refetchProgress }] = createResource(getProgress);
  
  // Объединяем квизы из прогресса с локальными карточками
  const allTasks = (): TaskUnion[] => {
    const localTasks = getAllActivities().filter(t => t.type === 'flashcard'); // Только карточки из локального хранилища
    const progress = progressData() || [];
    
    // Извлекаем уникальные квизы из прогресса
    const uniqueQuizzes = new Map<number, apiClient.QuizShortResponse>();
    progress.forEach(p => {
      if (p.quiz && !uniqueQuizzes.has(p.quiz.id)) {
        // Преобразуем QuizProgressResponse.quiz в QuizShortResponse формат
        uniqueQuizzes.set(p.quiz.id, {
          id: p.quiz.id,
          title: p.quiz.title,
          summary: p.quiz.summary,
          total_score: 0, // Не доступно в progress
          owner: {
            id: p.quiz.owner.id,
            short_name: p.quiz.owner.short_name,
            role: apiClient.UserRole.STUDENT, // Не доступно в progress, используем дефолт
          },
          subject: {
            id: p.quiz.subject.id,
            name: p.quiz.subject.name,
          },
          group: [], // Не доступно в progress
        });
      }
    });
    
    // Преобразуем квизы из прогресса в формат TaskUnion
    const apiQuizzes: TaskUnion[] = Array.from(uniqueQuizzes.values()).map(quiz => ({
      id: `quiz_${quiz.id}`,
      type: 'quiz' as const,
      title: quiz.title,
      category: quiz.subject?.name || 'unknown',
      teacher: quiz.owner?.short_name || 'Неизвестно',
      questions: [], // Вопросы загружаются отдельно при необходимости
      hidden: false,
      requiresConfirmation: false,
      requiresTeacherCheck: false,
    }));
    
    return [...apiQuizzes, ...localTasks];
  };
  const [showAddModal, setShowAddModal] = createSignal(false);
  const [showEditModal, setShowEditModal] = createSignal(false);
  const [showDeleteModal, setShowDeleteModal] = createSignal(false);
  const [editTask, setEditTask] = createSignal<TaskUnion | null>(null);
  const [deleteTask, setDeleteTask] = createSignal<TaskUnion | null>(null);

  // Управляемые поля формы
  const [formType, setFormType] = createSignal<'quiz' | 'flashcard'>('quiz');
  const [formCategory, setFormCategory] = createSignal('math');
  const [formSubjectId, setFormSubjectId] = createSignal<number | null>(null);
  const [formGroupIds, setFormGroupIds] = createSignal<number[]>([]);
  const [formTitle, setFormTitle] = createSignal('');
  const [formSummary, setFormSummary] = createSignal('');
  const [formAnswer, setFormAnswer] = createSignal('');
  const [formQuestions, setFormQuestions] = createSignal([
    { question: '', options: ['', '', '', ''], correct: 0 },
  ]);
  const [formTeacher, setFormTeacher] = createSignal('');
  const [formAttachments, setFormAttachments] = createSignal<MaterialAttachment[]>([]);
  const [newAttachmentType, setNewAttachmentType] = createSignal('pdf');
  const [newAttachmentUrl, setNewAttachmentUrl] = createSignal('');
  const [newAttachmentName, setNewAttachmentName] = createSignal('');
  const [formRequiresConfirmation, setFormRequiresConfirmation] = createSignal(false);
  const [isCreatingQuiz, setIsCreatingQuiz] = createSignal(false);

  // Проверка авторизации
  const user = getCurrentUser();
  if (!user) {
    return <div style={{ 'text-align': 'center', 'margin-top': '2rem', color: '#e76f51' }}>Нет доступа. Пожалуйста, войдите в систему.</div>;
  }

  // Фильтрация по предмету, типу и учителю
  const filteredTasks = () =>
    allTasks().filter(task =>
      (subjectFilter() === 'all' || task.category === subjectFilter()) &&
      (typeFilter() === 'all' || task.type === typeFilter()) &&
      (teacherFilter() === 'all' || task.teacher === teacherFilter()) &&
      // Ученики не видят скрытые задания
      (!(user.role === 'student') || !task.hidden)
    );

  // Скрыть/показать задание (только для локальных карточек)
  function handleToggleHidden(task: TaskUnion) {
    if (task.type === 'flashcard') {
      updateActivity({ ...task, hidden: !task.hidden });
      // Локальные карточки обновляются автоматически через allTasks()
    }
  }

  // Получаем список уникальных учителей
  const allTeachers = Array.from(new Set(allTasks().map(t => t.teacher)));

  // Загружаем предметы через API только если пользователь авторизован
  const [subjectsData] = createResource(() => {
    if (!user || !getAuthToken()) return Promise.resolve([]);
    return getSubjects();
  });
  const [usersData] = createResource(() => {
    if (!user || !getAuthToken()) return Promise.resolve([]);
    return getAllUsers();
  });
  const [groupsData] = createResource(() => {
    if (!user || !getAuthToken()) return Promise.resolve([]);
    return getGroups();
  });

  // Получить предметы для учителя
  const teacherSubjects = () => {
    const subjects = subjectsData();
    if (!subjects) return [];
    if (user && user.role === 'teacher' && 'subjects' in user) {
      const teacher = user as { subjects: string[] };
      return subjects.filter(s => teacher.subjects.includes(s.id.toString()));
    }
    return subjects;
  };

  // Получить список учителей для выбора (для админа)
  const teacherNames = () => {
    const users = usersData();
    if (!users) return [];
    // API /users возвращает short_name; пока не фильтруем по роли, чтобы не ломать UX
    return users.map(u => u.short_name);
  };

  // Сброс формы
  function resetForm() {
    setFormType('quiz');
    setFormCategory('math');
    setFormSubjectId(null);
    setFormGroupIds([]);
    setFormTitle('');
    setFormSummary('');
    setFormAnswer('');
    setFormQuestions([{ question: '', options: ['', '', '', ''], correct: 0 }]);
    setFormTeacher(user.name);
    setFormAttachments([]);
    setNewAttachmentType('pdf');
    setNewAttachmentUrl('');
    setNewAttachmentName('');
    setFormRequiresConfirmation(false);
  }

  // Добавление задания
  async function handleAddTask(e: Event) {
    e.preventDefault();
    
    if (formType() === 'quiz') {
      // Создаем квиз через API
      if (!formSubjectId()) {
        alert('Пожалуйста, выберите предмет');
        return;
      }
      
      setIsCreatingQuiz(true);
      try {
        // Преобразуем вопросы в формат API
        const apiQuestions: apiClient.QuizQuestion[] = formQuestions()
          .filter(q => q.question.trim() !== '') // Фильтруем пустые вопросы
          .map(q => {
            // Если есть варианты ответов, создаем options
            const allOptions = q.options || [];
            const validOptions = allOptions.filter(opt => opt.trim() !== '');
            
            if (validOptions.length > 0) {
              // Находим правильный ответ среди валидных опций
              const originalCorrectIndex = q.correct;
              const originalCorrectOption = allOptions[originalCorrectIndex];
              const correctIndexInValid = validOptions.findIndex(opt => opt === originalCorrectOption);
              
              return {
                text: q.question.trim(),
                options: validOptions.map((opt, idx) => ({
                  text: opt.trim(),
                  is_correct: idx === correctIndexInValid && correctIndexInValid >= 0
                }))
              };
            } else {
              // Если вариантов нет - открытый вопрос (пустой массив options)
              return {
                text: q.question.trim(),
                options: []
              };
            }
          });

        const quizData: apiClient.QuizCreateRequest = {
          title: formTitle(),
          summary: formSummary() || formTitle(), // Используем summary или title как fallback
          subject_id: formSubjectId()!,
          group_ids: formGroupIds(), // Всегда массив, может быть пустым
          questions: apiQuestions
        };

        await createQuiz(quizData);
        
        // Обновляем список квизов через прогресс
        await refetchProgress();
        
        setShowAddModal(false);
        resetForm();
        alert('Квиз успешно создан!');
      } catch (error) {
        console.error('Ошибка создания квиза:', error);
        alert('Ошибка создания квиза: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
      } finally {
        setIsCreatingQuiz(false);
      }
    } else {
      // Карточки пока сохраняем локально (API для них нет)
      const id = `flashcard_${Date.now()}`;
      addActivity({
        id,
        type: 'flashcard',
        question: formTitle(),
        answer: formAnswer(),
        category: formCategory(),
        teacher: user.role === 'admin' ? formTeacher() : user.name,
        attachments: formAttachments(),
        requiresConfirmation: formRequiresConfirmation(),
      });
      // Локальные карточки обновляются автоматически через allTasks()
      setShowAddModal(false);
      resetForm();
    }
  }

  // Открыть модалку редактирования
  function openEditModal(task: TaskUnion) {
    setEditTask(task);
    setFormType(task.type);
    setFormCategory(task.category);
    setFormTitle(task.type === 'quiz' ? task.title : task.question);
    setFormAnswer(task.type === 'flashcard' ? task.answer : '');
    setFormQuestions(task.type === 'quiz' ? [...task.questions] : [{ question: '', options: ['', '', '', ''], correct: 0 }]);
    setFormTeacher(task.teacher);
    setFormAttachments(task.attachments ? [...task.attachments] : []);
    setFormRequiresConfirmation(!!task.requiresConfirmation);
  }

  // Сохранить изменения задания
  function handleEditTask(e: Event) {
    e.preventDefault();
    if (!editTask()) return;
    const id = editTask()!.id;
    if (formType() === 'quiz') {
      // Для квизов из API редактирование через API (пока не реализовано)
      // Для локальных квизов (если они есть)
      if (id.startsWith('quiz_')) {
        alert('Редактирование квизов через API пока не реализовано');
        return;
      }
      updateActivity({
        id,
        type: 'quiz',
        title: formTitle(),
        category: formCategory(),
        teacher: user && user.role === 'admin' ? formTeacher() : user.name,
        questions: formQuestions(),
        attachments: formAttachments(),
        requiresConfirmation: formRequiresConfirmation(),
      });
    } else {
      updateActivity({
        id,
        type: 'flashcard',
        question: formTitle(),
        answer: formAnswer(),
        category: formCategory(),
        teacher: user && user.role === 'admin' ? formTeacher() : user.name,
        attachments: formAttachments(),
        requiresConfirmation: formRequiresConfirmation(),
      });
    }
    // Локальные карточки обновляются автоматически через allTasks()
    setShowEditModal(false);
    setEditTask(null);
    resetForm();
  }

  // Открыть модалку удаления
  function openDeleteModal(task: TaskUnion) {
    setDeleteTask(task);
    setShowDeleteModal(true);
  }

  // Удалить задание
  function handleDeleteTask() {
    if (!deleteTask()) return;
    const id = deleteTask()!.id;
    if (id.startsWith('quiz_')) {
      // Удаление квиза через API (пока не реализовано)
      alert('Удаление квизов через API пока не реализовано');
      setShowDeleteModal(false);
      return;
    }
    removeActivity(id);
    // Локальные карточки обновляются автоматически через allTasks()
    setShowDeleteModal(false);
    setDeleteTask(null);
  }

  // Проверка прав на редактирование/удаление
  function canEditOrDelete(task: TaskUnion) {
    return user && (user.role === 'admin' || (user.role === 'teacher' && task.teacher === user.name));
  }

  // Проверка роли для отображения кнопки добавления
  const canAddTask = user.role === 'admin' || user.role === 'teacher';

  return (
    <>
      <Header />
      <div class="main-shell">
        {/* основной контент Tasks */}
        <h2 style={{ 'margin-bottom': '2rem', color: '#2563eb', 'font-size': '2rem', 'font-weight': 700, 'letter-spacing': '0.01em' }}>Интерактивные задания</h2>
        {/* Кнопка добавить */}
        {canAddTask && (
          <button
            style={{
              'background': 'linear-gradient(90deg, #2563eb 60%, #e76f51 100%)',
              color: '#fff',
              'font-weight': 600,
              'font-size': '1.1em',
              padding: '0.7em 1.7em',
              'border-radius': '10px',
              border: 'none',
              cursor: 'pointer',
              'margin-bottom': '1.5rem',
              'box-shadow': '0 2px 8px rgba(37,99,235,0.10)',
              transition: 'background 0.2s',
            }}
            onClick={() => setShowAddModal(true)}
          >
            + Добавить задание
          </button>
        )}
        {/* Фильтры */}
        <div style={{ display: 'flex', gap: '1.5rem', 'margin-bottom': '2rem', 'flex-wrap': 'wrap' }}>
          <div>
            <label style={{ 'font-weight': 500, color: '#213547', 'margin-right': '0.5em' }}>Предмет:</label>
            <select value={subjectFilter()} onInput={e => setSubjectFilter(e.currentTarget.value)} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
              <option value="all">Все</option>
              <Show when={subjectsData()} fallback={<option disabled>Загрузка...</option>}>
                <For each={subjectsData() || []}>
                  {(subj) => (
                    <option value={subj.id.toString()}>{subj.name}</option>
                  )}
                </For>
              </Show>
            </select>
          </div>
          <div>
            <label style={{ 'font-weight': 500, color: '#213547', 'margin-right': '0.5em' }}>Тип задания:</label>
            <select value={typeFilter()} onInput={e => setTypeFilter(e.currentTarget.value)} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
              <option value="all">Все</option>
              <option value="quiz">Тест</option>
              <option value="flashcard">Карточка</option>
            </select>
          </div>
          <div>
            <label style={{ 'font-weight': 500, color: '#213547', 'margin-right': '0.5em' }}>Преподаватель:</label>
            <select value={teacherFilter()} onInput={e => setTeacherFilter(e.currentTarget.value)} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
              <option value="all">Все</option>
              {allTeachers.map(teacher => (
                <option value={teacher}>{teacher}</option>
              ))}
            </select>
          </div>
        </div>
        <Show when={filteredTasks().length > 0} fallback={<div style={{ color: '#888', 'font-size': '1.1em' }}>Нет доступных заданий.</div>}>
          <ul style={{ 'list-style': 'none', padding: 0, display: 'grid', gap: '3rem', 'grid-template-columns': 'repeat(auto-fit, minmax(260px, 1fr))' }}>
            <For each={filteredTasks()}>{task => (
              <li
                style={{
                  background: 'var(--bg-primary)',
                  color: 'var(--text-primary)',
                  'border-radius': '14px',
                  padding: '1.5rem 1.2rem',
                  'box-shadow': '0 2px 12px rgba(37,99,235,0.09)',
                  cursor: 'pointer',
                  transition: 'box-shadow 0.2s, transform 0.2s',
                  border: '1.5px solid #e3eafc',
                  'min-height': '120px',
                  'display': 'flex',
                  'flex-direction': 'column',
                  'justify-content': 'center',
                  position: 'relative',
                  overflow: 'hidden',
                }}
                onClick={() => setSelectedTask(task)}
                onMouseOver={e => (e.currentTarget.style.boxShadow = '0 4px 16px rgba(37,99,235,0.13)')}
                onMouseOut={e => (e.currentTarget.style.boxShadow = '0 2px 12px rgba(37,99,235,0.09)')}
              >
                {/* Цветная полоска слева */}
                <span style={{
                  position: 'absolute',
                  left: 0,
                  top: 0,
                  height: '100%',
                  width: '7px',
                  background: task.type === 'quiz' ? 'linear-gradient(180deg, #2563eb 60%, #e76f51 100%)' : 'linear-gradient(180deg, #e76f51 60%, #2563eb 100%)',
                  'border-top-left-radius': '14px',
                  'border-bottom-left-radius': '14px',
                }}></span>
                <span style={{ 'font-size': '1.1em', 'font-weight': 600, color: '#213547', 'margin-bottom': '0.5em', 'margin-left': '12px' }}>{task.type === 'quiz' ? 'Тест' : 'Карточка'}</span>
                <span style={{ 'font-size': '1.15em', color: '#2563eb', 'font-weight': 700, 'margin-left': '12px' }}>{task.type === 'quiz' ? task.title : task.question}</span>
                <span style={{ color: '#888', 'font-size': '0.98em', 'margin-top': '0.5em', 'margin-left': '12px' }}>Категория: {task.category}</span>
                <span style={{ color: '#888', 'font-size': '0.98em', 'margin-top': '0.2em', 'margin-left': '12px' }}>Преподаватель: {task.teacher}</span>
                {/* Кнопки редактировать/удалить */}
                {canEditOrDelete(task) && (
                  <div style={{ display: 'flex', gap: '0.5em', marginTop: '0.7em' }}>
                    <button title={task.hidden ? 'Показать задание' : 'Скрыть задание'} style={{ background: '#e3eafc', color: task.hidden ? '#888' : '#2563eb', border: 'none', borderRadius: '7px', padding: '0.3em 0.9em', cursor: 'pointer', fontWeight: 600 }} onClick={e => { e.stopPropagation(); handleToggleHidden(task); }}>
                      {task.hidden ? '👁‍🗨' : '👁'}
                    </button>
                    <button style={{ background: '#e3eafc', color: '#2563eb', border: 'none', borderRadius: '7px', padding: '0.3em 0.9em', cursor: 'pointer', fontWeight: 600 }} onClick={e => { e.stopPropagation(); openEditModal(task); }}>✎</button>
                    <button style={{ background: '#fff0f0', color: '#e76f51', border: 'none', borderRadius: '7px', padding: '0.3em 0.9em', cursor: 'pointer', fontWeight: 600 }} onClick={e => { e.stopPropagation(); openDeleteModal(task); }}>🗑</button>
                    <button style={{ background: '#e3eafc', color: '#2563eb', border: 'none', borderRadius: '7px', padding: '0.3em 0.9em', cursor: 'pointer', fontWeight: 600 }} onClick={e => { e.stopPropagation(); setSelectedTask(task); }}>Перейти</button>
                  </div>
                )}
              </li>
            )}</For>
          </ul>
        </Show>
        {/* Модальное окно задания */}
        <Show when={!!selectedTask()}>
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              width: '100vw',
              height: '100vh',
              background: 'rgba(37,99,235,0.10)',
              'z-index': 1000,
              display: 'flex',
              'align-items': 'center',
              'justify-content': 'center',
              transition: 'background 0.2s',
              animation: 'fadeInBg 0.2s',
            }}
            onClick={() => setSelectedTask(null)}
          >
            <div
              style={{
                background: 'var(--bg-primary)',
                'border-radius': '16px',
                padding: '2.2rem 2rem 2rem 2rem',
                'min-width': '320px',
                'max-width': '95vw',
                'box-shadow': '0 8px 32px rgba(37,99,235,0.18)',
                position: 'relative',
                animation: 'fadeInModal 0.25s',
              }}
              onClick={e => e.stopPropagation()}
            >
              <button
                style={{
                  position: 'absolute',
                  top: '1.2rem',
                  right: '1.5rem',
                  'font-size': '1.7em',
                  background: 'none',
                  border: 'none',
                  cursor: 'pointer',
                  color: '#2563eb',
                  transition: 'color 0.2s',
                }}
                onMouseOver={e => (e.currentTarget.style.color = '#e76f51')}
                onMouseOut={e => (e.currentTarget.style.color = '#2563eb')}
                onClick={() => setSelectedTask(null)}
                aria-label="Закрыть"
              >
                &times;
              </button>
              {selectedTask()?.type === 'quiz' ? (
                <>
                  <h3 style={{ color: '#2563eb', 'margin-bottom': '0.7em' }}>{(selectedTask() as any).title}</h3>
                  <div style={{ color: '#888', 'margin-bottom': '1.2em' }}>
                    Тест по категории: {(selectedTask() as any).category}
                  </div>
                  <ol>
                    {((selectedTask() as any)?.questions || []).map((q: any, idx: number) => (
                      <li style={{ 'margin-bottom': '1em' }}>
                        <div style={{ 'font-weight': 500 }}>{q.question}</div>
                        <ul style={{ 'padding-left': '1.2em', 'margin-top': '0.5em' }}>
                          {(q.options || []).map((opt: string, i: number) => (
                            <li>{opt}</li>
                          ))}
                        </ul>
                      </li>
                    ))}
                  </ol>
                </>
              ) : (
                <>
                  <h3 style={{ color: '#2563eb', 'margin-bottom': '0.7em' }}>Карточка</h3>
                  <div style={{ 'font-size': '1.1em', 'margin-bottom': '0.7em' }}>
                    <b>Вопрос:</b> {(selectedTask() as any).question}
                  </div>
                  <div style={{ 'margin-top': '1em', color: '#213547', 'font-size': '1.08em' }}>
                    <b>Ответ:</b> {(selectedTask() as any).answer}
                  </div>
                </>
              )}
              {selectedTask()?.attachments && selectedTask()?.attachments.length > 0 && (
                <div style={{ marginTop: '1.5em', marginBottom: '1em' }}>
                  <b style={{ color: '#2563eb', fontSize: '1.08em' }}>Материалы:</b>
                  <ul style={{ paddingLeft: '1.2em', marginTop: '0.7em' }}>
                    {selectedTask()?.attachments.map(att => (
                      <li style={{ marginBottom: '0.5em' }}>
                        <span style={{ color: '#2563eb', fontWeight: 600, marginRight: '0.5em' }}>{att.type.toUpperCase()}:</span>
                        {att.type === 'video' ? (
                          <a href={att.url} target="_blank" rel="noopener noreferrer" style={{ color: '#e76f51', textDecoration: 'underline' }}>Смотреть видео</a>
                        ) : att.type === 'pdf' ? (
                          <a href={att.url} target="_blank" rel="noopener noreferrer" style={{ color: '#2563eb', textDecoration: 'underline' }}>Открыть PDF</a>
                        ) : att.type === 'link' ? (
                          <a href={att.url} target="_blank" rel="noopener noreferrer" style={{ color: '#2563eb', textDecoration: 'underline' }}>Перейти по ссылке</a>
                        ) : (
                          <span>{att.url}</span>
                        )}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
              <Show when={selectedTask() && user.role === 'student'}>
                <div style={{ marginTop: '1.5em' }}>
                  <Show when={!selectedTask().requiresConfirmation}>
                    <button style={{ background: '#2a9d8f', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.7em 2em', fontWeight: 600, fontSize: '1.1em', cursor: 'pointer', boxShadow: '0 2px 8px rgba(42,157,143,0.10)' }}
                      onClick={() => {
                        // Добавить запись в activityHistory
                        const all = getAllUsers();
                        const idx = all.findIndex(u => u.id === user.id);
                        if (idx !== -1) {
                          const history = all[idx].activityHistory || [];
                          history.push({
                            activityId: selectedTask().id,
                            type: selectedTask().type,
                            title: selectedTask().type === 'quiz' ? selectedTask().title : selectedTask().question,
                            date: new Date().toISOString(),
                            status: 'passed',
                          });
                          all[idx].activityHistory = history;
                          localStorage.setItem('users', JSON.stringify(all));
                          alert('Урок отмечен как пройден!');
                        }
                      }}>
                      Я прошёл урок
                    </button>
                  </Show>
                  <Show when={selectedTask().requiresConfirmation}>
                    <form onSubmit={e => {
                      e.preventDefault();
                      const form = e.target as HTMLFormElement;
                      const files = (form.elements.namedItem('confirmationFiles') as HTMLInputElement).files;
                      const fileList: string[] = [];
                      if (files) {
                        for (let i = 0; i < files.length; i++) {
                          fileList.push(files[i].name);
                        }
                      }
                      // Добавить запись в activityHistory со статусом in_progress
                      const all = getAllUsers();
                      const idx = all.findIndex(u => u.id === user.id);
                      if (idx !== -1) {
                        const history = all[idx].activityHistory || [];
                        history.push({
                          activityId: selectedTask().id,
                          type: selectedTask().type,
                          title: selectedTask().type === 'quiz' ? selectedTask().title : selectedTask().question,
                          date: new Date().toISOString(),
                          status: 'in_progress',
                          files: fileList,
                        });
                        all[idx].activityHistory = history;
                        localStorage.setItem('users', JSON.stringify(all));
                        alert('Файлы отправлены! Ожидайте проверки.');
                      }
                    }} style={{ marginTop: '1em', display: 'flex', flexDirection: 'column', gap: '1em' }}>
                      <input type="file" name="confirmationFiles" multiple style={{ fontSize: '1em' }} />
                      <button type="submit" style={{ background: '#2a9d8f', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.7em 2em', fontWeight: 600, fontSize: '1.1em', cursor: 'pointer', boxShadow: '0 2px 8px rgba(42,157,143,0.10)' }}>
                        Отправить подтверждение
                      </button>
                    </form>
                  </Show>
                </div>
              </Show>
            </div>
          </div>
        </Show>
        {/* Модальное окно добавления задания */}
        <Show when={showAddModal()}>
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              width: '100vw',
              height: '100vh',
              background: 'rgba(37,99,235,0.10)',
              'z-index': 1000,
              display: 'flex',
              'align-items': 'center',
              'justify-content': 'center',
              transition: 'background 0.2s',
              animation: 'fadeInBg 0.2s',
            }}
            onClick={() => { setShowAddModal(false); resetForm(); }}
          >
            <div
              style={{
                background: 'var(--bg-primary)',
                'border-radius': '24px',
                padding: '2.5rem 2.5rem 2.5rem 2.5rem',
                width: '700px',
                height: '700px',
                'max-width': '90vw',
                'max-height': '90vh',
                'box-shadow': '0 12px 48px rgba(37,99,235,0.25)',
                position: 'relative',
                animation: 'fadeInModal 0.25s',
                display: 'flex',
                'flex-direction': 'column',
                'overflow-y': 'auto',
                'overflow-x': 'hidden',
              }}
              class="modal-scroll"
              onClick={e => e.stopPropagation()}
            >
              <div style={{ position: 'sticky', top: 0, background: 'var(--bg-primary)', 'z-index': 5, 'padding-bottom': '1em', 'margin-bottom': '1em' }}>
                <button
                  style={{
                    position: 'absolute',
                    top: '0',
                    right: '0',
                    'font-size': '1.7em',
                    background: 'none',
                    border: 'none',
                    cursor: 'pointer',
                    color: '#2563eb',
                    transition: 'color 0.2s',
                    'z-index': 10,
                  }}
                  onMouseOver={e => (e.currentTarget.style.color = '#e76f51')}
                  onMouseOut={e => (e.currentTarget.style.color = '#2563eb')}
                  onClick={() => { setShowAddModal(false); resetForm(); }}
                  aria-label="Закрыть"
                >
                  &times;
                </button>
                <h3 style={{ color: '#2563eb', 'margin-bottom': '0', 'font-size': '1.5em', 'font-weight': 700, 'padding-right': '2.5rem' }}>Добавить задание</h3>
              </div>
              <form onSubmit={handleAddTask} style={{ 'flex': 1, 'overflow-y': 'auto', 'overflow-x': 'hidden', 'padding-right': '0.5rem' }}>
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>Тип задания:</label><br />
                  <select value={formType()} onInput={e => setFormType(e.currentTarget.value as 'quiz' | 'flashcard')} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
                    <option value="quiz">Тест</option>
                    <option value="flashcard">Карточка</option>
                  </select>
                </div>
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>Предмет:</label><br />
                  <select value={formSubjectId()?.toString() || ''} onInput={e => {
                    const value = e.currentTarget.value;
                    setFormSubjectId(value ? parseInt(value, 10) : null);
                  }} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }} required>
                    <option value="">Выберите предмет</option>
                    <For each={teacherSubjects()}>
                      {(subj) => (
                        <option value={subj.id.toString()}>{subj.name}</option>
                      )}
                    </For>
                  </select>
                </div>
                {formType() === 'quiz' && (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Группы (необязательно):</label><br />
                    <select 
                      multiple 
                      value={formGroupIds().map(id => id.toString())} 
                      onInput={e => {
                        const selected = Array.from(e.currentTarget.selectedOptions, opt => parseInt(opt.value, 10));
                        setFormGroupIds(selected);
                      }} 
                      style={{ width: '100%', padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600, minHeight: '100px' }}
                    >
                      <Show when={groupsData()} fallback={<option disabled>Загрузка...</option>}>
                        <For each={groupsData() || []}>
                          {(group) => (
                            <option value={group.id.toString()}>{group.name}</option>
                          )}
                        </For>
                      </Show>
                    </select>
                    <div style={{ fontSize: '0.85em', color: '#888', marginTop: '0.3em' }}>Удерживайте Ctrl (Cmd на Mac) для выбора нескольких групп</div>
                  </div>
                )}
                {user && user.role === 'admin' && (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Учитель:</label><br />
                    <select value={formTeacher()} onInput={e => setFormTeacher(e.currentTarget.value)} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
                      {teacherNames().map(name => (
                        <option value={name}>{name}</option>
                      ))}
                    </select>
                  </div>
                )}
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>{formType() === 'quiz' ? 'Название теста' : 'Вопрос'}:</label><br />
                  <input type="text" value={formTitle()} onInput={e => setFormTitle(e.currentTarget.value)} style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} required />
                </div>
                {formType() === 'quiz' && (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Описание (summary):</label><br />
                    <textarea value={formSummary()} onInput={e => setFormSummary(e.currentTarget.value)} style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc', minHeight: '80px', resize: 'vertical' }} placeholder="Краткое описание квиза" />
                  </div>
                )}
                {formType() === 'quiz' ? (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Вопросы теста:</label>
                    {formQuestions().map((q, idx) => (
                      <div style={{ 'margin-bottom': '0.7em', 'border': '1px solid #e3eafc', 'border-radius': '8px', padding: '0.7em' }}>
                        <input 
                          type="text" 
                          value={q.question} 
                          onInput={(e) => {
                            const arr = formQuestions().map((item, i) => 
                              i === idx ? { ...item, question: e.currentTarget.value } : item
                            );
                            setFormQuestions(arr);
                          }} 
                          placeholder={`Вопрос ${idx + 1}`} 
                          style={{ width: '100%', marginBottom: '0.5em', padding: '0.4em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                        />
                        {q.options.map((opt, oidx) => (
                          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.3em' }}>
                            <input 
                              type="text" 
                              value={opt} 
                              onInput={(e) => {
                                const arr = formQuestions().map((item, i) => {
                                  if (i === idx) {
                                    const newOptions = [...item.options];
                                    newOptions[oidx] = e.currentTarget.value;
                                    return { ...item, options: newOptions };
                                  }
                                  return item;
                                });
                                setFormQuestions(arr);
                              }} 
                              placeholder={`Вариант ${oidx + 1}`} 
                              style={{ flex: 1, padding: '0.3em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                            />
                            <input 
                              type="radio" 
                              name={`correct${idx}`} 
                              checked={q.correct === oidx} 
                              onChange={() => {
                                const arr = formQuestions().map((item, i) => 
                                  i === idx ? { ...item, correct: oidx } : item
                                );
                                setFormQuestions(arr);
                              }} 
                              style={{ marginLeft: '0.7em' }} 
                            />
                            <span style={{ marginLeft: '0.3em', color: '#2563eb', fontSize: '0.95em' }}>Правильный</span>
                          </div>
                        ))}
                        <button type="button" onClick={() => {
                          const arr = [...formQuestions()];
                          arr.splice(idx, 1);
                          setFormQuestions(arr);
                        }} style={{ color: '#e76f51', background: 'none', border: 'none', cursor: 'pointer', marginTop: '0.3em' }}>Удалить вопрос</button>
                      </div>
                    ))}
                    <button type="button" onClick={() => setFormQuestions([...formQuestions(), { question: '', options: ['', '', '', ''], correct: 0 }])} style={{ background: '#2563eb', color: '#fff', fontWeight: 600, padding: '0.4em 1em', borderRadius: '8px', border: 'none', cursor: 'pointer', marginTop: '0.5em' }}>+ Добавить вопрос</button>
                  </div>
                ) : (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Ответ:</label><br />
                    <input type="text" value={formAnswer()} onInput={e => setFormAnswer(e.currentTarget.value)} style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} required />
                  </div>
                )}
                <div style={{ 'margin-bottom': '1em', 'border': '1px solid #e3eafc', 'border-radius': '8px', padding: '0.7em' }}>
                  <label style={{ 'font-weight': 500, display: 'block', 'margin-bottom': '0.5em' }}>Прикреплённые материалы:</label>
                  <For each={formAttachments()}>{(att, idx) => (
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.7em', marginBottom: '0.4em' }}>
                      <span style={{ color: '#2563eb', fontWeight: 600 }}>{att.type.toUpperCase()}</span>
                      <span>{att.url}</span>
                      <button type="button" style={{ color: '#e76f51', background: 'none', border: 'none', cursor: 'pointer' }} onClick={() => {
                        const arr = [...formAttachments()];
                        arr.splice(idx(), 1);
                        setFormAttachments(arr);
                      }}>Удалить</button>
                    </div>
                  )}</For>
                  <div style={{ display: 'flex', gap: '0.7em', marginTop: '0.5em' }}>
                    <select value={newAttachmentType()} onInput={e => setNewAttachmentType(e.currentTarget.value)} style={{ padding: '0.3em', borderRadius: '6px', border: '1px solid #e3eafc' }}>
                      <option value="pdf">PDF</option>
                      <option value="video">Видео</option>
                      <option value="link">Ссылка</option>
                    </select>
                    <input type="text" placeholder="Ссылка или путь к файлу" value={newAttachmentUrl()} onInput={e => setNewAttachmentUrl(e.currentTarget.value)} style={{ flex: 1, padding: '0.3em', borderRadius: '6px', border: '1px solid #e3eafc' }} />
                    <button type="button" style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.3em 1em', fontWeight: 600, cursor: 'pointer' }} onClick={() => {
                      if (!newAttachmentUrl()) return;
                      setFormAttachments([...formAttachments(), { type: newAttachmentType(), url: newAttachmentUrl() }]);
                      setNewAttachmentUrl('');
                    }}>Добавить</button>
                  </div>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5em', marginTop: '1em', fontWeight: 500 }}>
                  <input type="checkbox" checked={formRequiresConfirmation()} onInput={e => setFormRequiresConfirmation(e.currentTarget.checked)} />
                  Требует подтверждения (загрузка файлов)
                </div>
                <button type="submit" style={{ background: '#2563eb', color: '#fff', 'font-weight': 600, padding: '0.6em 1.5em', 'border-radius': '8px', border: 'none', cursor: 'pointer' }}>
                  Сохранить
                </button>
              </form>
            </div>
          </div>
        </Show>
        {/* Модалка редактирования */}
        <Show when={showEditModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', 'z-index': 1000, display: 'flex', 'align-items': 'center', 'justify-content': 'center', transition: 'background 0.2s', animation: 'fadeInBg 0.2s' }} onClick={() => { setShowEditModal(false); setEditTask(null); resetForm(); }}>
            <div style={{ background: 'var(--bg-primary)', 'border-radius': '16px', padding: '2.2rem 2rem 2rem 2rem', 'min-width': '340px', 'max-width': '95vw', 'box-shadow': '0 8px 32px rgba(37,99,235,0.18)', position: 'relative', animation: 'fadeInModal 0.25s' }} onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.2rem', right: '1.5rem', 'font-size': '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s' }} onMouseOver={e => (e.currentTarget.style.color = '#e76f51')} onMouseOut={e => (e.currentTarget.style.color = '#2563eb')} onClick={() => { setShowEditModal(false); setEditTask(null); resetForm(); }} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', 'margin-bottom': '1.2em' }}>Редактировать задание</h3>
              <form onSubmit={handleEditTask}>
                {/* Форма аналогична добавлению, но с предзаполнением */}
                {/* Тип задания и id менять нельзя */}
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>Тип задания:</label><br />
                  <input type="text" value={formType() === 'quiz' ? 'Тест' : 'Карточка'} disabled style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600, background: '#f3f6fa' }} />
                </div>
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>Предмет:</label><br />
                  <select value={formSubjectId()?.toString() || ''} onInput={e => {
                    const value = e.currentTarget.value;
                    setFormSubjectId(value ? parseInt(value, 10) : null);
                  }} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }} required>
                    <option value="">Выберите предмет</option>
                    <For each={teacherSubjects()}>
                      {(subj) => (
                        <option value={subj.id.toString()}>{subj.name}</option>
                      )}
                    </For>
                  </select>
                </div>
                {user && user.role === 'admin' && (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Учитель:</label><br />
                    <select value={formTeacher()} onInput={e => setFormTeacher(e.currentTarget.value)} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
                      {teacherNames().map(name => (
                        <option value={name}>{name}</option>
                      ))}
                    </select>
                  </div>
                )}
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>{formType() === 'quiz' ? 'Название теста' : 'Вопрос'}:</label><br />
                  <input type="text" value={formTitle()} onInput={e => setFormTitle(e.currentTarget.value)} style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} required />
                </div>
                {formType() === 'quiz' ? (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Вопросы теста:</label>
                    {formQuestions().map((q, idx) => (
                      <div style={{ 'margin-bottom': '0.7em', 'border': '1px solid #e3eafc', 'border-radius': '8px', padding: '0.7em' }}>
                        <input 
                          type="text" 
                          value={q.question} 
                          onInput={(e) => {
                            const arr = formQuestions().map((item, i) => 
                              i === idx ? { ...item, question: e.currentTarget.value } : item
                            );
                            setFormQuestions(arr);
                          }} 
                          placeholder={`Вопрос ${idx + 1}`} 
                          style={{ width: '100%', marginBottom: '0.5em', padding: '0.4em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                        />
                        {q.options.map((opt, oidx) => (
                          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.3em' }}>
                            <input 
                              type="text" 
                              value={opt} 
                              onInput={(e) => {
                                const arr = formQuestions().map((item, i) => {
                                  if (i === idx) {
                                    const newOptions = [...item.options];
                                    newOptions[oidx] = e.currentTarget.value;
                                    return { ...item, options: newOptions };
                                  }
                                  return item;
                                });
                                setFormQuestions(arr);
                              }} 
                              placeholder={`Вариант ${oidx + 1}`} 
                              style={{ flex: 1, padding: '0.3em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                            />
                            <input 
                              type="radio" 
                              name={`correct${idx}`} 
                              checked={q.correct === oidx} 
                              onChange={() => {
                                const arr = formQuestions().map((item, i) => 
                                  i === idx ? { ...item, correct: oidx } : item
                                );
                                setFormQuestions(arr);
                              }} 
                              style={{ marginLeft: '0.7em' }} 
                            />
                            <span style={{ marginLeft: '0.3em', color: '#2563eb', fontSize: '0.95em' }}>Правильный</span>
                          </div>
                        ))}
                        <button type="button" onClick={() => {
                          const arr = [...formQuestions()];
                          arr.splice(idx, 1);
                          setFormQuestions(arr);
                        }} style={{ color: '#e76f51', background: 'none', border: 'none', cursor: 'pointer', marginTop: '0.3em' }}>Удалить вопрос</button>
                      </div>
                    ))}
                    <button type="button" onClick={() => setFormQuestions([...formQuestions(), { question: '', options: ['', '', '', ''], correct: 0 }])} style={{ background: '#2563eb', color: '#fff', fontWeight: 600, padding: '0.4em 1em', borderRadius: '8px', border: 'none', cursor: 'pointer', marginTop: '0.5em' }}>+ Добавить вопрос</button>
                  </div>
                ) : (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Ответ:</label><br />
                    <input type="text" value={formAnswer()} onInput={e => setFormAnswer(e.currentTarget.value)} style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} required />
                  </div>
                )}
                <div style={{ 'margin-bottom': '1em', 'border': '1px solid #e3eafc', 'border-radius': '8px', padding: '0.7em' }}>
                  <label style={{ 'font-weight': 500, display: 'block', 'margin-bottom': '0.5em' }}>Прикреплённые материалы:</label>
                  <For each={formAttachments()}>{(att, idx) => (
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.7em', marginBottom: '0.4em' }}>
                      <span style={{ color: '#2563eb', fontWeight: 600 }}>{att.type.toUpperCase()}</span>
                      <span>{att.url}</span>
                      <button type="button" style={{ color: '#e76f51', background: 'none', border: 'none', cursor: 'pointer' }} onClick={() => {
                        const arr = [...formAttachments()];
                        arr.splice(idx(), 1);
                        setFormAttachments(arr);
                      }}>Удалить</button>
                    </div>
                  )}</For>
                  <div style={{ display: 'flex', gap: '0.7em', marginTop: '0.5em' }}>
                    <select value={newAttachmentType()} onInput={e => setNewAttachmentType(e.currentTarget.value)} style={{ padding: '0.3em', borderRadius: '6px', border: '1px solid #e3eafc' }}>
                      <option value="pdf">PDF</option>
                      <option value="video">Видео</option>
                      <option value="link">Ссылка</option>
                    </select>
                    <input type="text" placeholder="Ссылка или путь к файлу" value={newAttachmentUrl()} onInput={e => setNewAttachmentUrl(e.currentTarget.value)} style={{ flex: 1, padding: '0.3em', borderRadius: '6px', border: '1px solid #e3eafc' }} />
                    <button type="button" style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.3em 1em', fontWeight: 600, cursor: 'pointer' }} onClick={() => {
                      if (!newAttachmentUrl()) return;
                      setFormAttachments([...formAttachments(), { type: newAttachmentType(), url: newAttachmentUrl() }]);
                      setNewAttachmentUrl('');
                    }}>Добавить</button>
                  </div>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5em', marginTop: '1em', fontWeight: 500 }}>
                  <input type="checkbox" checked={formRequiresConfirmation()} onInput={e => setFormRequiresConfirmation(e.currentTarget.checked)} />
                  Требует подтверждения (загрузка файлов)
                </div>
                <button type="submit" style={{ background: '#2563eb', color: '#fff', 'font-weight': 600, padding: '0.6em 1.5em', 'border-radius': '8px', border: 'none', cursor: 'pointer' }}>
                  Сохранить
                </button>
              </form>
            </div>
          </div>
        </Show>
        {/* Модалка подтверждения удаления */}
        <Show when={showDeleteModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', 'z-index': 1000, display: 'flex', 'align-items': 'center', 'justify-content': 'center', transition: 'background 0.2s', animation: 'fadeInBg 0.2s' }} onClick={() => { setShowDeleteModal(false); setDeleteTask(null); }}>
            <div style={{ background: 'var(--bg-primary)', 'border-radius': '16px', padding: '2.2rem 2rem 2rem 2rem', 'min-width': '320px', 'max-width': '95vw', 'box-shadow': '0 8px 32px rgba(37,99,235,0.18)', position: 'relative', animation: 'fadeInModal 0.25s' }} onClick={e => e.stopPropagation()}>
              <h3 style={{ color: '#e76f51', 'margin-bottom': '1.2em' }}>Удалить задание?</h3>
              <div style={{ marginBottom: '1.2em', color: '#213547' }}>Вы уверены, что хотите удалить это задание?</div>
              <div style={{ display: 'flex', gap: '1em', 'justify-content': 'flex-end' }}>
                <button style={{ background: '#e3eafc', color: '#2563eb', border: 'none', 'border-radius': '7px', padding: '0.5em 1.2em', cursor: 'pointer', fontWeight: 600 }} onClick={() => { setShowDeleteModal(false); setDeleteTask(null); }}>Отмена</button>
                <button style={{ background: '#fff0f0', color: '#e76f51', border: 'none', 'border-radius': '7px', padding: '0.5em 1.2em', cursor: 'pointer', fontWeight: 600 }} onClick={handleDeleteTask}>Удалить</button>
              </div>
            </div>
          </div>
        </Show>
      </div>
      {/* Анимации для модалки и стили скроллбара */}
      <style>{`
        @keyframes fadeInBg {
          from { background: rgba(37,99,235,0); }
          to { background: rgba(37,99,235,0.10); }
        }
        @keyframes fadeInModal {
          from { opacity: 0; transform: translateY(30px); }
          to { opacity: 1; transform: none; }
        }
        /* Стили скроллбара для модальных окон */
        .modal-scroll::-webkit-scrollbar {
          width: 10px;
        }
        .modal-scroll::-webkit-scrollbar-track {
          background: #f1f1f1;
          border-radius: 10px;
        }
        .modal-scroll::-webkit-scrollbar-thumb {
          background: #2563eb;
          border-radius: 10px;
        }
        .modal-scroll::-webkit-scrollbar-thumb:hover {
          background: #1e4ed8;
        }
      `}</style>
    </>
  );
};

export default Tasks; 