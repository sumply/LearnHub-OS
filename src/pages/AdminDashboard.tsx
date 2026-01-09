import { createSignal, For, Show, createResource, createMemo } from 'solid-js';
import Header from '../components/Header';
import { A } from '@solidjs/router';
import { getAllUsers, getAllUsersWithRoles } from '../services/userService';
import { getSubjects, getGroups, createGroup, addStudentsToGroup, createSubject, getGroupStudents } from '../services/subjectGroupService';
import { createUser } from '../utils/apiClient';
import * as apiClient from '../utils/apiClient';

const tableStyle = {
  width: '100%',
  borderCollapse: 'collapse',
  background: 'var(--bg-primary)',
  fontSize: '1em',
  color: 'var(--text-primary)',
};
const thStyle = {
  border: '1px solid var(--border-color)',
  padding: '0.5em 0.8em',
  textAlign: 'left' as const,
  background: 'var(--bg-secondary)',
  fontWeight: 600,
  color: 'var(--text-primary)',
};
const tdStyle = {
  border: '1px solid var(--border-color)',
  padding: '0.5em 0.8em',
  textAlign: 'left' as const,
  color: 'var(--text-primary)',
};

const AdminDashboard = () => {
  // Загружаем данные через API
  const [usersData] = createResource(getAllUsers);
  const [usersWithRolesData] = createResource(getAllUsersWithRoles);
  const [subjectsData] = createResource(getSubjects);
  const [groupsData] = createResource(getGroups);
  
  // Разделяем пользователей по ролям
  const students = createMemo(() => {
    const users = usersWithRolesData();
    if (!users) return [];
    const filtered = users.filter(u => u.role === apiClient.UserRole.STUDENT);
    if (import.meta.env.DEV) {
      console.log('[AdminDashboard] Студенты:', filtered.map(u => ({ id: u.id, name: `${u.last_name} ${u.first_name}`, role: u.role })));
    }
    return filtered;
  });
  
  const teachers = createMemo(() => {
    const users = usersWithRolesData();
    if (!users) return [];
    const filtered = users.filter(u => u.role === apiClient.UserRole.TEACHER);
    if (import.meta.env.DEV) {
      console.log('[AdminDashboard] Учителя:', filtered.map(u => ({ id: u.id, name: `${u.last_name} ${u.first_name}`, role: u.role })));
    }
    return filtered;
  });
  
  const admins = createMemo(() => {
    const users = usersWithRolesData();
    if (!users) return [];
    const filtered = users.filter(u => u.role === apiClient.UserRole.ADMIN || u.role === apiClient.UserRole.ROOT);
    if (import.meta.env.DEV) {
      console.log('[AdminDashboard] Админы:', filtered.map(u => ({ id: u.id, name: `${u.last_name} ${u.first_name}`, role: u.role })));
    }
    return filtered;
  });
  
  // Получаем все группы из API
  const groups = createMemo(() => {
    const groups = groupsData();
    if (!groups) return [];
    return groups;
  });
  
  // Загружаем студентов для каждой группы и создаем маппинг студент -> группа
  const [studentGroupMap] = createResource(async () => {
    const groupsList = groups();
    if (!groupsList || groupsList.length === 0) return new Map<number, apiClient.GroupResponse>();
    
    const map = new Map<number, apiClient.GroupResponse>(); // studentId -> GroupResponse
    
    // Для каждой группы получаем список студентов
    await Promise.all(groupsList.map(async (group) => {
      try {
        const students = await getGroupStudents(group.id);
        students.forEach(student => {
          map.set(student.id, group);
        });
      } catch (error) {
        // Если endpoint не существует, просто игнорируем ошибку
        console.warn(`Не удалось загрузить студентов группы ${group.id}:`, error);
      }
    }));
    
    return map;
  });
  
  // Функция для получения группы студента
  const getStudentGroup = (studentId: number): apiClient.GroupResponse | null => {
    return studentGroupMap()?.get(studentId) || null;
  };
  
  // Сигналы для модальных окон
  const [showAddModal, setShowAddModal] = createSignal(false);
  const [showAddGroupModal, setShowAddGroupModal] = createSignal(false);
  const [showAddStudentsModal, setShowAddStudentsModal] = createSignal(false);
  const [selectedGroup, setSelectedGroup] = createSignal<apiClient.GroupResponse | null>(null);
  
  // Форма добавления пользователя
  const [addRole, setAddRole] = createSignal('student');
  const [addFields, setAddFields] = createSignal({ name: '', surname: '', email: '', group: '', subjects: [] as string[] });
  
  // Форма создания группы
  const [groupName, setGroupName] = createSignal('');
  const [groupCuratorId, setGroupCuratorId] = createSignal<number | null>(null);
  const [groupStudentIds, setGroupStudentIds] = createSignal<number[]>([]);
  
  // Форма добавления студентов в группу
  const [selectedStudentIds, setSelectedStudentIds] = createSignal<number[]>([]);
  
  // Форма создания предмета
  const [showAddSubjectModal, setShowAddSubjectModal] = createSignal(false);
  const [subjectName, setSubjectName] = createSignal('');
  
  function handleAddField(field: string, value: any) {
    setAddFields({ ...addFields(), [field]: value });
  }
  
  function handleAddSubjects(subjId: string) {
    const current = addFields().subjects || [];
    if (current.includes(subjId)) {
      setAddFields({ ...addFields(), subjects: current.filter(s => s !== subjId) });
    } else {
      setAddFields({ ...addFields(), subjects: [...current, subjId] });
    }
  }
  
  async function saveNewUser() {
    try {
      const roleMap: Record<string, apiClient.UserRole> = {
        'student': apiClient.UserRole.STUDENT,
        'teacher': apiClient.UserRole.TEACHER,
        'admin': apiClient.UserRole.ADMIN,
      };
      
      await createUser({
        first_name: addFields().name,
        last_name: addFields().surname,
        email: addFields().email,
        role: roleMap[addRole()] || apiClient.UserRole.STUDENT,
      });
      
      setShowAddModal(false);
      setAddFields({ name: '', surname: '', email: '', group: '', subjects: [] });
      // Обновляем данные - очищаем кэш и перезагружаем
      import('../services/userService').then(module => module.clearUsersCache());
      window.location.reload();
    } catch (error) {
      console.error('Ошибка создания пользователя:', error);
      alert('Ошибка создания пользователя: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    }
  }
  
  function toggleGroupStudentSelection(studentId: number) {
    const current = groupStudentIds();
    if (current.includes(studentId)) {
      setGroupStudentIds(current.filter(id => id !== studentId));
    } else {
      setGroupStudentIds([...current, studentId]);
    }
  }
  
  async function saveNewGroup() {
    try {
      if (!groupName() || !groupCuratorId()) {
        alert('Пожалуйста, заполните название группы и выберите куратора');
        return;
      }
      
      // Создаем группу
      await createGroup(groupName(), groupCuratorId()!);
      
      // Получаем ID созданной группы (нужно будет получить из ответа или перезагрузить список)
      // Пока что перезагружаем список групп и находим созданную группу по имени
      import('../services/subjectGroupService').then(module => module.clearCache());
      
      // Если выбраны студенты, добавляем их в группу
      if (groupStudentIds().length > 0) {
        // Ждем немного, чтобы группа успела создаться
        await new Promise(resolve => setTimeout(resolve, 500));
        
        // Получаем обновленный список групп
        const updatedGroups = await getGroups();
        const createdGroup = updatedGroups.find(g => g.name === groupName());
        
        if (createdGroup) {
          await addStudentsToGroup(createdGroup.id, groupStudentIds());
        } else {
          console.warn('Не удалось найти созданную группу для добавления студентов');
        }
      }
      
      setShowAddGroupModal(false);
      setGroupName('');
      setGroupCuratorId(null);
      setGroupStudentIds([]);
      
      // Обновляем данные - очищаем кэш и перезагружаем
      import('../services/subjectGroupService').then(module => module.clearCache());
      import('../services/userService').then(module => module.clearUsersCache());
      window.location.reload();
    } catch (error) {
      console.error('Ошибка создания группы:', error);
      alert('Ошибка создания группы: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    }
  }
  
  async function saveNewSubject() {
    try {
      if (!subjectName()) {
        alert('Пожалуйста, введите название предмета');
        return;
      }
      
      await createSubject(subjectName());
      
      setShowAddSubjectModal(false);
      setSubjectName('');
      // Обновляем данные - очищаем кэш и перезагружаем
      import('../services/subjectGroupService').then(module => module.clearCache());
      window.location.reload();
    } catch (error) {
      console.error('Ошибка создания предмета:', error);
      alert('Ошибка создания предмета: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    }
  }
  
  function openAddStudentsModal(group: apiClient.GroupResponse) {
    setSelectedGroup(group);
    setSelectedStudentIds([]);
    setShowAddStudentsModal(true);
  }
  
  function toggleStudentSelection(studentId: number) {
    const current = selectedStudentIds();
    if (current.includes(studentId)) {
      setSelectedStudentIds(current.filter(id => id !== studentId));
    } else {
      setSelectedStudentIds([...current, studentId]);
    }
  }
  
  async function saveStudentsToGroup() {
    try {
      const group = selectedGroup();
      if (!group || selectedStudentIds().length === 0) {
        alert('Пожалуйста, выберите группу и студентов');
        return;
      }
      
      await addStudentsToGroup(group.id, selectedStudentIds());
      
      setShowAddStudentsModal(false);
      setSelectedGroup(null);
      setSelectedStudentIds([]);
      // Обновляем данные - очищаем кэш и перезагружаем
      import('../services/subjectGroupService').then(module => module.clearCache());
      import('../services/userService').then(module => module.clearUsersCache());
      window.location.reload();
    } catch (error) {
      console.error('Ошибка добавления студентов в группу:', error);
      alert('Ошибка добавления студентов в группу: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    }
  }

  return (
    <>
      <Header />
      <div class="main-shell">
        <h2 style={{ marginBottom: '2rem', color: '#2563eb', fontSize: '2rem', fontWeight: 700, letterSpacing: '0.01em' }}>Админ-панель</h2>
        
        {/* Кнопки действий */}
        <div style={{ display: 'flex', gap: '1em', marginBottom: '2rem', flexWrap: 'wrap' }}>
          <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.6em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddModal(true)}>Добавить пользователя</button>
          <button style={{ background: '#10b981', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.6em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddGroupModal(true)}>Создать группу</button>
          <button style={{ background: '#8b5cf6', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.6em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddSubjectModal(true)}>Создать предмет</button>
        </div>
        
        {/* Таблица студентов */}
        <div style={{ marginBottom: '3rem' }}>
          <h3 style={{ color: '#2563eb', marginBottom: '1rem', fontSize: '1.5rem', fontWeight: 600 }}>Ученики (роль 0)</h3>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={thStyle}>ID</th>
                <th style={thStyle}>ФИО</th>
                <th style={thStyle}>Группа</th>
                <th style={thStyle}>Действия</th>
              </tr>
            </thead>
            <tbody>
              <For each={students()}>{(student, idx) => {
                const studentGroup = getStudentGroup(student.id);
                const groupName = studentGroup?.name;
                const displayGroup = groupName && groupName !== '0' && groupName !== 'null' && groupName.trim() !== '' 
                  ? groupName 
                  : 'Ученик не состоит в группе';
                return (
                  <tr style={idx() % 2 === 1 ? { background: 'var(--bg-secondary)' } : {}}>
                    <td style={tdStyle}>{student.id}</td>
                    <td style={tdStyle}>{student.last_name} {student.first_name} {student.middle_name || ''}</td>
                    <td style={tdStyle}>{displayGroup}</td>
                    <td style={tdStyle}>
                      <button style={{ background: 'none', color: '#2563eb', border: '1px solid #2563eb', borderRadius: '3px', padding: '0.2em 0.7em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em', marginRight: '0.5em' }}>
                        <A href={`/admin/profile/${student.id}`} style={{ color: '#2563eb', textDecoration: 'none' }}>Профиль</A>
                      </button>
                    </td>
                  </tr>
                );
              }}</For>
              <Show when={students().length === 0}>
                <tr>
                  <td colSpan={4} style={{ ...tdStyle, textAlign: 'center', color: 'var(--text-secondary)' }}>Нет учеников</td>
                </tr>
              </Show>
            </tbody>
          </table>
        </div>
        
        {/* Таблица учителей */}
        <div style={{ marginBottom: '3rem' }}>
          <h3 style={{ color: '#2563eb', marginBottom: '1rem', fontSize: '1.5rem', fontWeight: 600 }}>Учителя (роль 1)</h3>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={thStyle}>ID</th>
                <th style={thStyle}>ФИО</th>
                <th style={thStyle}>Действия</th>
              </tr>
            </thead>
            <tbody>
              <For each={teachers()}>{(teacher, idx) => (
                <tr style={idx() % 2 === 1 ? { background: '#f9fafb' } : {}}>
                  <td style={tdStyle}>{teacher.id}</td>
                  <td style={tdStyle}>{teacher.last_name} {teacher.first_name} {teacher.middle_name || ''}</td>
                  <td style={tdStyle}>
                    <button style={{ background: 'none', color: '#2563eb', border: '1px solid #2563eb', borderRadius: '3px', padding: '0.2em 0.7em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em', marginRight: '0.5em' }}>
                      <A href={`/admin/profile/${teacher.id}`} style={{ color: '#2563eb', textDecoration: 'none' }}>Профиль</A>
                    </button>
                  </td>
                </tr>
              )}</For>
              <Show when={teachers().length === 0}>
                <tr>
                  <td colSpan={3} style={{ ...tdStyle, textAlign: 'center', color: '#6b7280' }}>Нет учителей</td>
                </tr>
              </Show>
            </tbody>
          </table>
        </div>
        
        {/* Таблица админов */}
        <div style={{ marginBottom: '3rem' }}>
          <h3 style={{ color: '#2563eb', marginBottom: '1rem', fontSize: '1.5rem', fontWeight: 600 }}>Администраторы (роль 2)</h3>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={thStyle}>ID</th>
                <th style={thStyle}>ФИО</th>
                <th style={thStyle}>Действия</th>
              </tr>
            </thead>
            <tbody>
              <For each={admins()}>{(admin, idx) => (
                <tr style={idx() % 2 === 1 ? { background: 'var(--bg-secondary)' } : {}}>
                  <td style={tdStyle}>{admin.id}</td>
                  <td style={tdStyle}>{admin.last_name} {admin.first_name} {admin.middle_name || ''}</td>
                  <td style={tdStyle}>
                    <button style={{ background: 'none', color: '#2563eb', border: '1px solid #2563eb', borderRadius: '3px', padding: '0.2em 0.7em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em', marginRight: '0.5em' }}>
                      <A href={`/admin/profile/${admin.id}`} style={{ color: '#2563eb', textDecoration: 'none' }}>Профиль</A>
                    </button>
                  </td>
                </tr>
              )}</For>
              <Show when={admins().length === 0}>
                <tr>
                  <td colSpan={3} style={{ ...tdStyle, textAlign: 'center', color: '#6b7280' }}>Нет администраторов</td>
                </tr>
              </Show>
            </tbody>
          </table>
        </div>
        
        {/* Таблица предметов */}
        <div style={{ marginBottom: '3rem' }}>
          <h3 style={{ color: '#2563eb', marginBottom: '1rem', fontSize: '1.5rem', fontWeight: 600 }}>Предметы</h3>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={thStyle}>ID</th>
                <th style={thStyle}>Название</th>
              </tr>
            </thead>
            <tbody>
              <For each={subjectsData() || []}>{(subject, idx) => (
                <tr style={idx() % 2 === 1 ? { background: 'var(--bg-secondary)' } : {}}>
                  <td style={tdStyle}>{subject.id}</td>
                  <td style={tdStyle}>{subject.name}</td>
                </tr>
              )}</For>
              <Show when={!subjectsData() || subjectsData()!.length === 0}>
                <tr>
                  <td colSpan={2} style={{ ...tdStyle, textAlign: 'center', color: '#6b7280' }}>Нет предметов</td>
                </tr>
              </Show>
            </tbody>
          </table>
        </div>
        
        {/* Список групп */}
        <div style={{ marginBottom: '3rem' }}>
          <h3 style={{ color: '#2563eb', marginBottom: '1rem', fontSize: '1.5rem', fontWeight: 600 }}>Группы (классы)</h3>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={thStyle}>ID</th>
                <th style={thStyle}>Название</th>
                <th style={thStyle}>Куратор</th>
                <th style={thStyle}>Действия</th>
              </tr>
            </thead>
            <tbody>
              <For each={groups()}>{(group, idx) => (
                <tr style={idx() % 2 === 1 ? { background: 'var(--bg-secondary)' } : {}}>
                  <td style={tdStyle}>{group.id}</td>
                  <td style={tdStyle}>{group.name}</td>
                  <td style={tdStyle}>{group.curator.short_name}</td>
                  <td style={tdStyle}>
                    <button style={{ background: '#10b981', color: '#fff', border: 'none', borderRadius: '3px', padding: '0.3em 0.8em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em', marginRight: '0.5em' }} onClick={() => openAddStudentsModal(group)}>
                      Добавить студентов
                    </button>
                  </td>
                </tr>
              )}</For>
              <Show when={groups().length === 0}>
                <tr>
                  <td colSpan={4} style={{ ...tdStyle, textAlign: 'center', color: '#6b7280' }}>Нет групп</td>
                </tr>
              </Show>
            </tbody>
          </table>
        </div>
        
        {/* Модальное окно добавления пользователя */}
        <Show when={showAddModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setShowAddModal(false)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '24px', padding: '2.5rem 2.5rem 2.5rem 2.5rem', width: '600px', maxWidth: '90vw', maxHeight: '85vh', overflowY: 'auto', overflowX: 'hidden', boxShadow: '0 12px 48px rgba(37,99,235,0.25)', position: 'relative' }} class="modal-scroll" onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.5rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s', zIndex: 10 }} onClick={() => setShowAddModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.5em', fontSize: '1.5em', fontWeight: 700 }}>Добавить пользователя</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Роль:</label><br />
                <select value={addRole()} onInput={e => setAddRole(e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }}>
                  <option value="student">Ученик</option>
                  <option value="teacher">Учитель</option>
                  <option value="admin">Администратор</option>
                </select>
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Фамилия:</label><br />
                <input type="text" value={addFields().surname} onInput={e => handleAddField('surname', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Имя:</label><br />
                <input type="text" value={addFields().name} onInput={e => handleAddField('name', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Email:</label><br />
                <input type="email" value={addFields().email} onInput={e => handleAddField('email', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} required />
              </div>
              <div style={{ marginBottom: '0.7em', padding: '0.7em', background: '#e3eafc', borderRadius: '6px', fontSize: '0.9em', color: '#2563eb' }}>
                <strong>Примечание:</strong> Логин и пароль будут автоматически сгенерированы сервером и отправлены на указанный email.
              </div>
              <div style={{ display: 'flex', gap: '1em', marginTop: '1.5em' }}>
                <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={saveNewUser}>Сохранить</button>
                <button style={{ background: '#f3f4f6', color: '#213547', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddModal(false)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
        
        {/* Модальное окно создания группы */}
        <Show when={showAddGroupModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setShowAddGroupModal(false)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '24px', padding: '2.5rem 2.5rem 2.5rem 2.5rem', width: '600px', maxWidth: '90vw', maxHeight: '85vh', overflowY: 'auto', overflowX: 'hidden', boxShadow: '0 12px 48px rgba(37,99,235,0.25)', position: 'relative' }} class="modal-scroll" onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.5rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s', zIndex: 10 }} onClick={() => setShowAddGroupModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.5em', fontSize: '1.5em', fontWeight: 700 }}>Создать группу</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Название группы:</label><br />
                <input type="text" value={groupName()} onInput={e => setGroupName(e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} placeholder="Например: 10А" />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Куратор (учитель):</label><br />
                <select value={groupCuratorId()?.toString() || ''} onInput={e => setGroupCuratorId(e.currentTarget.value ? parseInt(e.currentTarget.value) : null)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }}>
                  <option value="">Выберите учителя</option>
                  <For each={teachers()}>{teacher => (
                    <option value={teacher.id.toString()}>{teacher.last_name} {teacher.first_name} {teacher.middle_name || ''}</option>
                  )}</For>
                </select>
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Студенты (необязательно):</label><br />
                <div style={{ maxHeight: '200px', overflowY: 'auto', border: '1px solid #d1d5db', borderRadius: '6px', padding: '0.5em', marginTop: '0.5em' }}>
                  <For each={students()}>{student => (
                    <label style={{ display: 'block', padding: '0.3em', cursor: 'pointer', borderRadius: '4px', background: groupStudentIds().includes(student.id) ? '#e3eafc' : 'transparent', marginBottom: '0.3em' }}>
                      <input 
                        type="checkbox" 
                        checked={groupStudentIds().includes(student.id)} 
                        onChange={() => toggleGroupStudentSelection(student.id)} 
                        style={{ marginRight: '0.5em' }} 
                      />
                      {student.last_name} {student.first_name} {student.middle_name || ''} (ID: {student.id})
                    </label>
                  )}</For>
                  <Show when={students().length === 0}>
                    <div style={{ textAlign: 'center', color: '#6b7280', padding: '1em' }}>Нет доступных студентов</div>
                  </Show>
                </div>
              </div>
              <div style={{ display: 'flex', gap: '1em', marginTop: '1.5em' }}>
                <button style={{ background: '#10b981', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={saveNewGroup}>Создать</button>
                <button style={{ background: '#f3f4f6', color: '#213547', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => { setShowAddGroupModal(false); setGroupStudentIds([]); }}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
        
        {/* Модальное окно добавления студентов в группу */}
        <Show when={showAddStudentsModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setShowAddStudentsModal(false)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '24px', padding: '2.5rem 2.5rem 2.5rem 2.5rem', width: '700px', maxWidth: '90vw', maxHeight: '85vh', overflowY: 'auto', overflowX: 'hidden', boxShadow: '0 12px 48px rgba(37,99,235,0.25)', position: 'relative' }} class="modal-scroll" onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.5rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s', zIndex: 10 }} onClick={() => setShowAddStudentsModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.5em', fontSize: '1.5em', fontWeight: 700 }}>Добавить студентов в группу: {selectedGroup()?.name}</h3>
              <div style={{ marginBottom: '1em', maxHeight: '400px', overflowY: 'auto' }}>
                <For each={students()}>{student => (
                  <label style={{ display: 'block', padding: '0.5em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid var(--border-color)', cursor: 'pointer', background: selectedStudentIds().includes(student.id) ? 'var(--bg-tertiary)' : 'var(--bg-secondary)', color: 'var(--text-primary)' }}>
                    <input type="checkbox" checked={selectedStudentIds().includes(student.id)} onChange={() => toggleStudentSelection(student.id)} style={{ marginRight: '0.5em' }} />
                    {student.last_name} {student.first_name} {student.middle_name || ''} (ID: {student.id})
                  </label>
                )}</For>
                <Show when={students().length === 0}>
                  <div style={{ textAlign: 'center', color: '#6b7280', padding: '2em' }}>Нет доступных студентов</div>
                </Show>
              </div>
              <div style={{ display: 'flex', gap: '1em', marginTop: '1.5em' }}>
                <button style={{ background: '#10b981', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={saveStudentsToGroup} disabled={selectedStudentIds().length === 0}>Добавить выбранных</button>
                <button style={{ background: '#f3f4f6', color: '#213547', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddStudentsModal(false)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
        
        {/* Модальное окно создания предмета */}
        <Show when={showAddSubjectModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setShowAddSubjectModal(false)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '24px', padding: '2.5rem 2.5rem 2.5rem 2.5rem', width: '600px', maxWidth: '90vw', maxHeight: '85vh', overflowY: 'auto', overflowX: 'hidden', boxShadow: '0 12px 48px rgba(37,99,235,0.25)', position: 'relative' }} class="modal-scroll" onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.5rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s', zIndex: 10 }} onClick={() => setShowAddSubjectModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.5em', fontSize: '1.5em', fontWeight: 700 }}>Создать предмет</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Название предмета:</label><br />
                <input type="text" value={subjectName()} onInput={e => setSubjectName(e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} placeholder="Например: Математика" />
              </div>
              <div style={{ display: 'flex', gap: '1em', marginTop: '1.5em' }}>
                <button style={{ background: '#8b5cf6', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={saveNewSubject}>Создать</button>
                <button style={{ background: '#f3f4f6', color: '#213547', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddSubjectModal(false)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
      </div>
      {/* Стили скроллбара для модальных окон */}
      <style>{`
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

export default AdminDashboard;
