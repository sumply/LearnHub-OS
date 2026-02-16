import { type Component, createSignal, For, Show, createResource, Index } from 'solid-js';
import { getCurrentUser } from '../utils/api';
import { getAllActivities, type TaskUnion, addActivity, updateActivity, removeActivity } from '../utils/activitiesService';
import Header from '../components/Header';
import { getSubjects, getGroups } from '../services/subjectGroupService';
import { getAllUsers } from '../services/userService';
import { getAuthToken } from '../utils/api';
import { createQuiz } from '../services/quizService';
import { getQuizzes } from '../services/quizService';
import * as apiClient from '../utils/apiClient';
import type { MaterialAttachment } from '../utils/api';
import QuizStatistics from '../components/QuizStatistics';

const Tasks: Component = () => {
  const [selectedTask, setSelectedTask] = createSignal<TaskUnion | null>(null);
  const [subjectFilter, setSubjectFilter] = createSignal('all');
  const [typeFilter, setTypeFilter] = createSignal('all');
  const [teacherFilter, setTeacherFilter] = createSignal('all');
  // Загружаем все квизы через GET /quizzes
  const [quizzesData, { refetch: refetchQuizzes }] = createResource(getQuizzes);
  
  // Объединяем квизы с локальными карточками
  const allTasks = (): TaskUnion[] => {
    const localTasks = getAllActivities().filter(t => t.type === 'flashcard'); // Только карточки из локального хранилища
    const quizzes = quizzesData() || [];
    
    // Преобразуем квизы в формат TaskUnion
    const apiQuizzes: TaskUnion[] = quizzes.map(quiz => ({
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
  const [showStatisticsModal, setShowStatisticsModal] = createSignal(false);
  const [statisticsQuizId, setStatisticsQuizId] = createSignal<string>('');
  const [statisticsQuizTitle, setStatisticsQuizTitle] = createSignal<string>('');
  const [editTask, setEditTask] = createSignal<TaskUnion | null>(null);
  const [deleteTask, setDeleteTask] = createSignal<TaskUnion | null>(null);

  // Управляемые поля формы
  const [formType, setFormType] = createSignal<'quiz' | 'flashcard'>('quiz');
  const [formCategory, setFormCategory] = createSignal('math');
  const [formSubjectId, setFormSubjectId] = createSignal<string | null>(null);
  const [formGroupIds, setFormGroupIds] = createSignal<string[]>([]);
  const [formTitle, setFormTitle] = createSignal('');
  const [formSummary, setFormSummary] = createSignal('');
  const [formAnswer, setFormAnswer] = createSignal('');
  const [formDeadline, setFormDeadline] = createSignal('');
  const [formMaxAttempts, setFormMaxAttempts] = createSignal<number | undefined>(undefined);
  const [formQuestions, setFormQuestions] = createSignal([
    { question: '', options: ['', '', '', ''], correct: 0, type: 'single' as 'single' | 'multiple' | 'numeric', score: 1 },
  ]);
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
    setFormDeadline('');
    setFormMaxAttempts(undefined);
    setFormQuestions([{ question: '', options: ['', '', '', ''], correct: 0, type: 'single', score: 1 }]);
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
      setIsCreatingQuiz(true);
      try {
        // Получаем owner_id из auth данных
        const authData = apiClient.getAuthData();
        if (!authData) {
          alert('Ошибка: пользователь не авторизован');
          return;
        }
        const ownerId = authData.user_id;

        // Преобразуем вопросы в формат API
        const apiQuestions: apiClient.QuizQuestion[] = formQuestions()
          .filter(q => q.question.trim() !== '') // Фильтруем пустые вопросы
          .map(q => {
            const questionText = q.question.trim();
            const score = q.score || 1;
            const type = q.type || 'single';
            
            // Обрабатываем варианты ответов в зависимости от типа вопроса
            const allOptions = q.options || [];
            const validOptions = allOptions.filter(opt => opt.trim() !== '');
            
            if (type === 'numeric') {
              // Для числового вопроса
              const numericValue = validOptions.length > 0 && !isNaN(parseFloat(validOptions[0])) 
                ? parseFloat(validOptions[0]) 
                : 0;
              return {
                text: questionText,
                score,
                type: 'numeric' as const,
                details: {
                  correct: numericValue
                } as apiClient.QuizQuestionDetailsNumeric
              };
            } else if (type === 'multiple') {
              // Для множественного выбора
              const originalCorrectIndices = Array.isArray(q.correct) ? q.correct : [q.correct];
              const correctAnswers = originalCorrectIndices
                .map(idx => allOptions[idx])
                .filter(opt => opt && opt.trim() !== '')
                .map(opt => opt.trim());
              
              return {
                text: questionText,
                score,
                type: 'multiple' as const,
                details: {
                  options: validOptions.map(opt => opt.trim()),
                  correct: correctAnswers
                } as apiClient.QuizQuestionDetailsMultiple
              };
            } else {
              // Для одиночного выбора (single) - по умолчанию
              const originalCorrectIndex = typeof q.correct === 'number' ? q.correct : 0;
              const originalCorrectOption = allOptions[originalCorrectIndex];
              const correctAnswer = originalCorrectOption && originalCorrectOption.trim() !== '' 
                ? originalCorrectOption.trim() 
                : (validOptions.length > 0 ? validOptions[0].trim() : '');
              
              return {
                text: questionText,
                score,
                type: 'single' as const,
                details: {
                  options: validOptions.map(opt => opt.trim()),
                  correct: correctAnswer
                } as apiClient.QuizQuestionDetailsSingle
              };
            }
          });

        // Формируем deadline в формате ISO, если указан
        let deadline: string | undefined = undefined;
        if (formDeadline()) {
          const deadlineDate = new Date(formDeadline());
          if (!isNaN(deadlineDate.getTime())) {
            deadline = deadlineDate.toISOString();
          }
        }

        // Определяем subject_id: если не выбран, используем "все темы" (возможно, нужно специальное значение или пустая строка)
        // Пока отправляем пустую строку или null, если не выбран
        const subjectId = formSubjectId() || 'все темы';
        
        // Определяем group_ids: если не выбраны, отправляем ["общий"]
        const groupIds = formGroupIds().length > 0 ? formGroupIds() : ['общий'];

        const quizData: apiClient.QuizCreateRequest = {
          title: formTitle(),
          owner_id: ownerId,
          summary: formSummary() || formTitle(),
          subject_id: subjectId, // UUID предмета или "все темы"
          group_ids: groupIds, // UUID[] групп или ["общий"]
          deadline: deadline,
          max_attempts: formMaxAttempts(),
          questions: apiQuestions
        };

        await createQuiz(quizData);
        
        // Обновляем список квизов через прогресс
        await refetchQuizzes();
        
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
        teacher: user.name,
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
          <ul class="tasks-grid" style={{ 
            'list-style': 'none', 
            padding: 0, 
            display: 'grid', 
            gap: '2rem', 
            'grid-template-columns': 'repeat(3, 1fr)'
          }}>
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
                  'min-height': '140px',
                  'display': 'flex',
                  'flex-direction': 'column',
                  'justify-content': 'space-between',
                  position: 'relative',
                  overflow: 'hidden',
                }}
                onClick={() => {
                  if (task.type === 'quiz' && user && user.role === 'student') {
                    // Для студентов - переходим на страницу информации о квизе
                    const quizId = task.id.replace('quiz_', '');
                    window.location.href = `/quiz/${quizId}`;
                  } else {
                    // Для остальных - открываем модальное окно
                    setSelectedTask(task);
                  }
                }}
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
                  <div style={{ 
                    display: 'flex', 
                    gap: '0.4em', 
                    marginTop: '0.7em',
                    flexWrap: 'wrap',
                    alignItems: 'center'
                  }}>
                    <button 
                      class="task-action-btn"
                      title={task.hidden ? 'Показать задание' : 'Скрыть задание'} 
                      style={{ 
                        background: '#e3eafc', 
                        color: task.hidden ? '#888' : '#2563eb', 
                        border: 'none', 
                        borderRadius: '6px', 
                        padding: '0',
                        cursor: 'pointer', 
                        fontWeight: 600,
                        fontSize: '0.9em',
                        width: '28px',
                        height: '28px',
                        minWidth: '28px',
                        maxWidth: '28px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0,
                        boxSizing: 'border-box'
                      }} 
                      onClick={e => { e.stopPropagation(); handleToggleHidden(task); }}
                    >
                      {task.hidden ? '👁‍🗨' : '👁'}
                    </button>
                    <button 
                      class="task-action-btn"
                      style={{ 
                        background: '#e3eafc', 
                        color: '#2563eb', 
                        border: 'none', 
                        borderRadius: '6px', 
                        padding: '0',
                        cursor: 'pointer', 
                        fontWeight: 600,
                        fontSize: '0.9em',
                        width: '28px',
                        height: '28px',
                        minWidth: '28px',
                        maxWidth: '28px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0,
                        boxSizing: 'border-box'
                      }} 
                      onClick={e => { e.stopPropagation(); openEditModal(task); }}
                      title="Редактировать"
                    >
                      ✎
                    </button>
                    <button 
                      class="task-action-btn"
                      style={{ 
                        background: '#fff0f0', 
                        color: '#e76f51', 
                        border: 'none', 
                        borderRadius: '6px', 
                        padding: '0',
                        cursor: 'pointer', 
                        fontWeight: 600,
                        fontSize: '0.9em',
                        width: '28px',
                        height: '28px',
                        minWidth: '28px',
                        maxWidth: '28px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0,
                        boxSizing: 'border-box'
                      }} 
                      onClick={e => { e.stopPropagation(); openDeleteModal(task); }}
                      title="Удалить"
                    >
                      🗑
                    </button>
                    {task.type === 'quiz' && (user.role === 'admin' || user.role === 'teacher') && (
                      <button 
                        class="task-action-btn"
                        style={{ 
                          background: '#e8f5e9', 
                          color: '#2e7d32', 
                          border: 'none', 
                          borderRadius: '6px', 
                          padding: '0',
                          cursor: 'pointer', 
                          fontWeight: 600,
                          fontSize: '0.9em',
                          width: '28px',
                          height: '28px',
                          minWidth: '28px',
                          maxWidth: '28px',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          flexShrink: 0,
                          boxSizing: 'border-box'
                        }} 
                        onClick={e => { 
                          e.stopPropagation(); 
                          const quizId = task.id.replace('quiz_', '');
                          setStatisticsQuizId(quizId);
                          setStatisticsQuizTitle(task.title);
                          setShowStatisticsModal(true);
                        }}
                        title="Статистика квиза"
                      >
                        📊
                      </button>
                    )}
                    <button 
                      class="task-action-btn"
                      style={{ 
                        background: '#e3eafc', 
                        color: '#2563eb', 
                        border: 'none', 
                        borderRadius: '6px', 
                        padding: '0',
                        cursor: 'pointer', 
                        fontWeight: 600,
                        fontSize: '0.9em',
                        width: '28px',
                        height: '28px',
                        minWidth: '28px',
                        maxWidth: '28px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0,
                        boxSizing: 'border-box'
                      }} 
                      onClick={e => { e.stopPropagation(); setSelectedTask(task); }}
                      title="Перейти"
                    >
                      →
                    </button>
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
                  <div style={{ color: '#888', 'margin-bottom': '1.2em' }}>
                    Преподаватель: {(selectedTask() as any).teacher}
                  </div>
                  <Show when={user && user.role === 'student'}>
                    <div style={{ 'margin-top': '2em', 'margin-bottom': '1em' }}>
                      <A 
                        href={`/quiz/take/${(selectedTask() as any).id.replace('quiz_', '')}`}
                        style={{
                          display: 'inline-block',
                          background: '#2563eb',
                          color: '#fff',
                          border: 'none',
                          borderRadius: '8px',
                          padding: '0.75em 2em',
                          fontWeight: 600,
                          cursor: 'pointer',
                          fontSize: '1.1em',
                          textDecoration: 'none',
                          transition: 'background 0.2s'
                        }}
                        onMouseOver={(e) => e.currentTarget.style.background = '#1e4ed8'}
                        onMouseOut={(e) => e.currentTarget.style.background = '#2563eb'}
                      >
                        Начать попытку
                      </A>
                    </div>
                  </Show>
                  <Show when={user && (user.role === 'admin' || user.role === 'teacher')}>
                    <div style={{ color: '#888', 'margin-bottom': '1.2em', fontSize: '0.9em' }}>
                      Для просмотра вопросов используйте кнопку "Статистика"
                    </div>
                  </Show>
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
                  <label style={{ 'font-weight': 500 }}>Предмет (необязательно):</label><br />
                  <select value={formSubjectId() || ''} onInput={e => {
                    const value = e.currentTarget.value;
                    setFormSubjectId(value || null);
                  }} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
                    <option value="">Все темы</option>
                    <For each={teacherSubjects()}>
                      {(subj) => (
                        <option value={typeof subj.id === 'string' ? subj.id : subj.id.toString()}>{subj.name}</option>
                      )}
                    </For>
                  </select>
                </div>
                {formType() === 'quiz' && (
                  <>
                    <div style={{ 'margin-bottom': '1em' }}>
                      <label style={{ 'font-weight': 500 }}>Группы (необязательно):</label><br />
                      <div style={{ border: '1px solid #e3eafc', borderRadius: '8px', padding: '0.5em', maxHeight: '200px', overflowY: 'auto' }}>
                        <Show when={groupsData()} fallback={<div style={{ color: '#888', padding: '0.5em' }}>Загрузка...</div>}>
                          <For each={groupsData() || []}>
                            {(group) => {
                              const groupId = typeof group.id === 'string' ? group.id : group.id.toString();
                              const isSelected = formGroupIds().includes(groupId);
                              return (
                                <label style={{ display: 'flex', alignItems: 'center', padding: '0.4em', cursor: 'pointer', borderRadius: '4px', marginBottom: '0.2em', background: isSelected ? '#e3eafc' : 'transparent' }}>
                                  <input 
                                    type="checkbox" 
                                    checked={isSelected}
                                    onChange={() => {
                                      const currentIds = formGroupIds();
                                      if (isSelected) {
                                        setFormGroupIds(currentIds.filter(id => id !== groupId));
                                      } else {
                                        setFormGroupIds([...currentIds, groupId]);
                                      }
                                    }}
                                    style={{ marginRight: '0.5em', cursor: 'pointer' }}
                                  />
                                  <span style={{ color: '#213547', fontSize: '0.95em' }}>{group.name}</span>
                                </label>
                              );
                            }}
                          </For>
                        </Show>
                        {(!groupsData() || groupsData()!.length === 0) && (
                          <div style={{ color: '#888', padding: '0.5em', fontSize: '0.9em' }}>Нет доступных групп</div>
                        )}
                      </div>
                    </div>
                    <div style={{ 'margin-bottom': '1em' }}>
                      <label style={{ 'font-weight': 500 }}>Дедлайн (необязательно):</label><br />
                      <input 
                        type="datetime-local" 
                        value={formDeadline()} 
                        max="2027-12-31T23:59"
                        onInput={e => setFormDeadline(e.currentTarget.value)} 
                        style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} 
                      />
                    </div>
                    <div style={{ 'margin-bottom': '1em' }}>
                      <label style={{ 'font-weight': 500 }}>Максимальное количество попыток (необязательно):</label><br />
                      <input 
                        type="number" 
                        min="1" 
                        value={formMaxAttempts() || ''} 
                        onInput={e => {
                          const value = e.currentTarget.value;
                          setFormMaxAttempts(value ? parseInt(value, 10) : undefined);
                        }} 
                        style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} 
                        placeholder="Например: 3"
                      />
                    </div>
                  </>
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
                    <Index each={formQuestions()}>
                      {(q, idx) => (
                      <div style={{ 'margin-bottom': '0.7em', 'border': '1px solid #e3eafc', 'border-radius': '8px', padding: '0.7em' }}>
                        <div style={{ display: 'flex', gap: '0.5em', marginBottom: '0.5em' }}>
                          <input 
                            type="text" 
                            value={q().question} 
                            onInput={(e) => {
                              const arr = formQuestions().map((item, i) => 
                                i === idx ? { ...item, question: e.currentTarget.value } : item
                              );
                              setFormQuestions(arr);
                            }} 
                            placeholder={`Вопрос ${idx + 1}`} 
                            style={{ flex: 1, padding: '0.4em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                          />
                          <select 
                            value={q().type} 
                            onInput={(e) => {
                              const arr = formQuestions().map((item, i) => 
                                i === idx ? { ...item, type: e.currentTarget.value as 'single' | 'multiple' | 'numeric' } : item
                              );
                              setFormQuestions(arr);
                            }} 
                            style={{ padding: '0.4em', border: '1px solid #e3eafc', borderRadius: '6px' }}
                          >
                            <option value="single">Одиночный выбор</option>
                            <option value="multiple">Множественный выбор</option>
                            <option value="numeric">Числовой ответ</option>
                          </select>
                          <input 
                            type="number" 
                            min="1" 
                            value={q().score || 1} 
                            onInput={(e) => {
                              const arr = formQuestions().map((item, i) => 
                                i === idx ? { ...item, score: parseInt(e.currentTarget.value, 10) || 1 } : item
                              );
                              setFormQuestions(arr);
                            }} 
                            placeholder="Баллы" 
                            style={{ width: '80px', padding: '0.4em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                          />
                        </div>
                        {q().type === 'numeric' ? (
                          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.3em' }}>
                            <input 
                              type="number" 
                              step="any"
                              value={q().options && q().options[0] ? q().options[0] : ''} 
                              onInput={(e) => {
                                const arr = formQuestions().map((item, i) => {
                                  if (i === idx) {
                                    const newOptions = [e.currentTarget.value, '', '', ''];
                                    return { ...item, options: newOptions };
                                  }
                                  return item;
                                });
                                setFormQuestions(arr);
                              }} 
                              placeholder="Правильный числовой ответ" 
                              style={{ flex: 1, padding: '0.3em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                            />
                          </div>
                        ) : (
                          <Index each={q().options}>
                            {(opt, oidx) => (
                            <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.3em' }}>
                              <input 
                                type="text" 
                                value={opt()} 
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
                              {q().type === 'single' ? (
                                <>
                                  <input 
                                    type="radio" 
                                    name={`correct${idx}`} 
                                    checked={q().correct === oidx} 
                                    onChange={() => {
                                      const arr = formQuestions().map((item, i) => 
                                        i === idx ? { ...item, correct: oidx } : item
                                      );
                                      setFormQuestions(arr);
                                    }} 
                                    style={{ marginLeft: '0.7em' }} 
                                  />
                                  <span style={{ marginLeft: '0.3em', color: '#2563eb', fontSize: '0.95em' }}>Правильный</span>
                                </>
                              ) : (
                                <>
                                  <input 
                                    type="checkbox" 
                                    checked={Array.isArray(q().correct) ? q().correct.includes(oidx) : false} 
                                    onChange={() => {
                                      const arr = formQuestions().map((item, i) => {
                                        if (i === idx) {
                                          const currentCorrect = Array.isArray(item.correct) ? item.correct : [];
                                          const newCorrect = currentCorrect.includes(oidx)
                                            ? currentCorrect.filter(c => c !== oidx)
                                            : [...currentCorrect, oidx];
                                          return { ...item, correct: newCorrect };
                                        }
                                        return item;
                                      });
                                      setFormQuestions(arr);
                                    }} 
                                    style={{ marginLeft: '0.7em' }} 
                                  />
                                  <span style={{ marginLeft: '0.3em', color: '#2563eb', fontSize: '0.95em' }}>Правильный</span>
                                </>
                              )}
                            </div>
                            )}
                          </Index>
                        )}
                        <button type="button" onClick={() => {
                          const arr = [...formQuestions()];
                          arr.splice(idx, 1);
                          setFormQuestions(arr);
                        }} style={{ color: '#e76f51', background: 'none', border: 'none', cursor: 'pointer', marginTop: '0.3em' }}>Удалить вопрос</button>
                      </div>
                      )}
                    </Index>
                    <button type="button" onClick={() => setFormQuestions([...formQuestions(), { question: '', options: ['', '', '', ''], correct: 0, type: 'single', score: 1 }])} style={{ background: '#2563eb', color: '#fff', fontWeight: 600, padding: '0.4em 1em', borderRadius: '8px', border: 'none', cursor: 'pointer', marginTop: '0.5em' }}>+ Добавить вопрос</button>
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
                  <label style={{ 'font-weight': 500 }}>Предмет (необязательно):</label><br />
                  <select value={formSubjectId() || ''} onInput={e => {
                    const value = e.currentTarget.value;
                    setFormSubjectId(value || null);
                  }} style={{ padding: '0.4em 1em', 'border-radius': '8px', border: '1.5px solid #e3eafc', color: '#2563eb', 'font-weight': 600 }}>
                    <option value="">Все темы</option>
                    <For each={teacherSubjects()}>
                      {(subj) => (
                        <option value={typeof subj.id === 'string' ? subj.id : subj.id.toString()}>{subj.name}</option>
                      )}
                    </For>
                  </select>
                </div>
                <div style={{ 'margin-bottom': '1em' }}>
                  <label style={{ 'font-weight': 500 }}>{formType() === 'quiz' ? 'Название теста' : 'Вопрос'}:</label><br />
                  <input type="text" value={formTitle()} onInput={e => setFormTitle(e.currentTarget.value)} style={{ width: '100%', padding: '0.5em', 'border-radius': '8px', border: '1.5px solid #e3eafc' }} required />
                </div>
                {formType() === 'quiz' ? (
                  <div style={{ 'margin-bottom': '1em' }}>
                    <label style={{ 'font-weight': 500 }}>Вопросы теста:</label>
                    <Index each={formQuestions()}>
                      {(q, idx) => (
                      <div style={{ 'margin-bottom': '0.7em', 'border': '1px solid #e3eafc', 'border-radius': '8px', padding: '0.7em' }}>
                        <input 
                          type="text" 
                          value={q().question} 
                          onInput={(e) => {
                            const arr = formQuestions().map((item, i) => 
                              i === idx ? { ...item, question: e.currentTarget.value } : item
                            );
                            setFormQuestions(arr);
                          }} 
                          placeholder={`Вопрос ${idx + 1}`} 
                          style={{ width: '100%', marginBottom: '0.5em', padding: '0.4em', border: '1px solid #e3eafc', borderRadius: '6px' }} 
                        />
                        <Index each={q().options}>
                          {(opt, oidx) => (
                          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.3em' }}>
                            <input 
                              type="text" 
                              value={opt()} 
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
                              checked={q().correct === oidx} 
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
                          )}
                        </Index>
                        <button type="button" onClick={() => {
                          const arr = [...formQuestions()];
                          arr.splice(idx, 1);
                          setFormQuestions(arr);
                        }} style={{ color: '#e76f51', background: 'none', border: 'none', cursor: 'pointer', marginTop: '0.3em' }}>Удалить вопрос</button>
                      </div>
                      )}
                    </Index>
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

      {/* Модальное окно статистики квиза */}
      <Show when={showStatisticsModal()}>
        <QuizStatistics
          quizId={statisticsQuizId()}
          quizTitle={statisticsQuizTitle()}
          onClose={() => {
            setShowStatisticsModal(false);
            setStatisticsQuizId('');
            setStatisticsQuizTitle('');
          }}
        />
      </Show>
    </>
  );
};

export default Tasks; 