import { createSignal, For, Show, createResource, createMemo } from 'solid-js';
import Header from '../components/Header';
import { A } from '@solidjs/router';
import { getAllUsers, getAllUsersWithRoles } from '../services/userService';
import { getSubjects, getGroups, createGroup, addStudentsToGroup, createSubject } from '../services/subjectGroupService';
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
    const filtered = users.filter(u => u.role === 'student' || u.role === 0);
    if (import.meta.env.DEV) {
      console.log('[AdminDashboard] Студенты:', filtered.map(u => ({ id: u.id, name: `${u.last_name} ${u.first_name}`, role: u.role })));
    }
    return filtered;
  });
  
  const teachers = createMemo(() => {
    const users = usersWithRolesData();
    if (!users) return [];
    const filtered = users.filter(u => u.role === 'teacher' || u.role === 1);
    if (import.meta.env.DEV) {
      console.log('[AdminDashboard] Учителя:', filtered.map(u => ({ id: u.id, name: `${u.last_name} ${u.first_name}`, role: u.role })));
    }
    return filtered;
  });
  
  const admins = createMemo(() => {
    const users = usersWithRolesData();
    if (!users) return [];
    const filtered = users.filter(u => u.role === 'admin' || u.role === 'root' || u.role === 2 || u.role === 3);
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
  
  // Создаем маппинг студент -> группа из данных групп (students уже включены в GroupResponse)
  const studentGroupMap = createMemo(() => {
    const groupsList = groups();
    if (!groupsList || groupsList.length === 0) return new Map<number | string, apiClient.GroupResponse>();
    
    const map = new Map<number | string, apiClient.GroupResponse>(); // studentId -> GroupResponse
    
    // Для каждой группы используем список студентов из group.students
    groupsList.forEach((group) => {
      if (group.students && group.students.length > 0) {
        group.students.forEach(student => {
          map.set(student.id, group);
        });
      }
    });
    
    return map;
  });
  
  // Функция для получения группы студента
  const getStudentGroup = (studentId: number | string): apiClient.GroupResponse | null => {
    return studentGroupMap()?.get(studentId) || null;
  };
  
  // Сигналы для модальных окон
  const [showAddModal, setShowAddModal] = createSignal(false);
  const [showAddGroupModal, setShowAddGroupModal] = createSignal(false);
  const [showAddStudentsModal, setShowAddStudentsModal] = createSignal(false);
  const [selectedGroup, setSelectedGroup] = createSignal<apiClient.GroupResponse | null>(null);
  
  // Форма добавления пользователя
  const [addRole, setAddRole] = createSignal('student');
  const [addFields, setAddFields] = createSignal({ 
    name: '', 
    surname: '', 
    email: '', 
    group: null as string | number | null, 
    subjects: [] as number[],
    groups: [] as (string | number)[] // Для учителей - группы, в которых они преподают
  });
  
  // Форма создания группы
  const [groupName, setGroupName] = createSignal('');
  const [groupCuratorId, setGroupCuratorId] = createSignal<number | null>(null);
  const [groupStudentIds, setGroupStudentIds] = createSignal<number[]>([]);
  
  // Форма добавления студентов в группу
  const [selectedStudentIds, setSelectedStudentIds] = createSignal<number[]>([]);
  
  // Форма создания предмета
  const [showAddSubjectModal, setShowAddSubjectModal] = createSignal(false);
  const [subjectName, setSubjectName] = createSignal('');
  
  function handleAddField(field: string, value: string | number | null) {
    setAddFields({ ...addFields(), [field]: value });
  }
  
  function toggleSubjectSelection(subjId: number) {
    const current = addFields().subjects || [];
    if (current.includes(subjId)) {
      setAddFields({ ...addFields(), subjects: current.filter(s => s !== subjId) });
    } else {
      setAddFields({ ...addFields(), subjects: [...current, subjId] });
    }
  }
  
  function toggleGroupSelection(groupId: number) {
    const current = addFields().groups || [];
    if (current.includes(groupId)) {
      setAddFields({ ...addFields(), groups: current.filter(g => g !== groupId) });
    } else {
      setAddFields({ ...addFields(), groups: [...current, groupId] });
    }
  }
  
  async function saveNewUser() {
    try {
      // Подготавливаем данные для создания пользователя
      const userData: apiClient.UserCreateRequest = {
        first_name: addFields().name,
        last_name: addFields().surname,
        email: addFields().email,
        role: addRole() || 'student',
      };
      
      // Если это учитель, добавляем предметы и группы
      if (addRole() === 'teacher') {
        if (addFields().subjects && addFields().subjects.length > 0) {
          userData.subject_ids = addFields().subjects;
        }
        if (addFields().groups && addFields().groups.length > 0) {
          userData.group_ids = addFields().groups;
        }
      }
      
      // Создаем пользователя
      await createUser(userData);
      
      // Получаем ID созданного пользователя из ответа или из списка пользователей
      // Сначала обновляем список пользователей
      import('../services/userService').then(module => module.clearUsersCache());
      const updatedUsers = await getAllUsers();
      const createdUser = updatedUsers.find(u => {
        const fullName = `${addFields().surname} ${addFields().name}`;
        return u.short_name === fullName || (u.short_name && addFields().name && u.short_name.includes(addFields().name));
      });
      
      if (createdUser) {
        const userId = createdUser.id;
        
        // Если это студент и выбрана группа - добавляем в группу
        if (addRole() === 'student' && addFields().group !== null) {
          try {
            await addStudentsToGroup(addFields().group!, [userId]);
          } catch (error) {
            console.warn('Не удалось добавить студента в группу:', error);
          }
        }
        
        // Если это учитель - добавляем в группы и предметы
        if (addRole() === 'teacher') {
          // Добавляем в группы (если выбраны)
          if (addFields().groups && addFields().groups.length > 0) {
            for (const groupId of addFields().groups) {
              try {
                // Предполагаем, что есть эндпоинт для добавления учителя в группу
                // Если нет, можно использовать существующий механизм
                await addStudentsToGroup(groupId, { student_ids: [] }); // Пока пропускаем, нужно уточнить API
              } catch (error) {
                console.warn(`Не удалось добавить учителя в группу ${groupId}:`, error);
              }
            }
          }
          // Предметы для учителя - нужно уточнить API для этого
        }
      }
      
      setShowAddModal(false);
      setAddFields({ name: '', surname: '', email: '', group: null, subjects: [], groups: [] });
      // Обновляем данные - очищаем кэш и перезагружаем
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
      if (!groupName()) {
        alert('Пожалуйста, введите название группы');
        return;
      }
      
      // Подготавливаем данные для создания группы
      const groupData: apiClient.GroupCreateRequest = {
        name: groupName(),
      };
      
      // Добавляем куратора, если выбран
      if (groupCuratorId() !== null) {
        // Нужно получить UUID учителя по его ID
        const teachersList = teachers();
        const selectedTeacher = teachersList.find(t => {
          const teacherId = typeof t.id === 'string' ? t.id : String(t.id);
          const curatorId = String(groupCuratorId()!);
          return teacherId === curatorId || String(t.id) === curatorId;
        });
        if (selectedTeacher) {
          // ID может быть как number, так и string (UUID) - приводим к string
          groupData.curator_id = String(selectedTeacher.id);
        }
      }
      
      // Добавляем студентов, если выбраны
      if (groupStudentIds().length > 0) {
        // Нужно получить UUID студентов по их ID
        const studentsList = students();
        const selectedStudents = studentsList.filter(s => groupStudentIds().includes(s.id));
        groupData.student_ids = selectedStudents.map(s => String(s.id));
      }
      
      // Создаем группу со всеми данными сразу
      await createGroup(groupData);
      
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
                  <td style={tdStyle}>{group.curator?.short_name || 'Не назначен'}</td>
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
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(0,0,0,0.3)', zIndex: 1000 }} onClick={() => setShowAddModal(false)}>
            <div style={{ 
              background: 'var(--bg-primary)', 
              borderRadius: '32px', 
              border: '6px solid #2563eb',
              padding: '2.5rem 2.5rem 2.5rem 2.5rem', 
              width: '450px', 
              maxHeight: '90vh', 
              overflowY: 'auto', 
              overflowX: 'hidden', 
              boxShadow: '0 20px 60px rgba(37,99,235,0.35)', 
              position: 'fixed',
              top: 'calc(50% - 100px)',
              left: '2rem',
              transform: 'translateY(-50%)',
              zIndex: 1001,
              animation: 'slideInFromLeft 0.3s ease-out'
            } as any} class="modal-scroll modal-window" onClick={e => e.stopPropagation()}>
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
              
              {/* Поля для студентов - выбор группы */}
              <Show when={addRole() === 'student'}>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Группа (необязательно):</label><br />
                  <select 
                    value={addFields().group?.toString() || ''} 
                    onInput={e => handleAddField('group', e.currentTarget.value ? parseInt(e.currentTarget.value) : null)} 
                    style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }}
                  >
                    <option value="">Не выбрано</option>
                    <For each={groups()}>
                      {(group) => (
                        <option value={group.id}>{group.name}</option>
                      )}
                    </For>
                  </select>
                </div>
              </Show>
              
              {/* Поля для учителей - выбор групп и предметов */}
              <Show when={addRole() === 'teacher'}>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Группы (необязательно):</label><br />
                  <div style={{ maxHeight: '150px', overflowY: 'auto', border: '1px solid #d1d5db', borderRadius: '6px', padding: '0.5em' }}>
                    <For each={groups()}>
                      {(group) => (
                        <label style={{ display: 'block', marginBottom: '0.3em', cursor: 'pointer' }}>
                          <input 
                            type="checkbox" 
                            checked={addFields().groups?.includes(group.id) || false}
                            onChange={() => toggleGroupSelection(group.id)}
                            style={{ marginRight: '0.5em' }}
                          />
                          {group.name}
                        </label>
                      )}
                    </For>
                  </div>
                </div>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Предметы (необязательно):</label><br />
                  <div style={{ maxHeight: '150px', overflowY: 'auto', border: '1px solid #d1d5db', borderRadius: '6px', padding: '0.5em' }}>
                    <For each={subjectsData() || []}>
                      {(subject) => (
                        <label style={{ display: 'block', marginBottom: '0.3em', cursor: 'pointer' }}>
                          <input 
                            type="checkbox" 
                            checked={addFields().subjects?.includes(subject.id) || false}
                            onChange={() => toggleSubjectSelection(subject.id)}
                            style={{ marginRight: '0.5em' }}
                          />
                          {subject.name}
                        </label>
                      )}
                    </For>
                  </div>
                </div>
              </Show>
              
              <div style={{ marginBottom: '0.7em', padding: '0.7em', background: '#e3eafc', borderRadius: '6px', fontSize: '0.9em', color: '#2563eb' }}>
                <strong>Примечание:</strong> Логин и пароль будут автоматически сгенерированы сервером и отправлены на указанный email.
              </div>
              <div style={{ display: 'flex', gap: '0.75em', marginTop: '3.5em', justifyContent: 'flex-end', flexWrap: 'wrap' }}>
                <button style={{ 
                  background: '#2563eb', 
                  color: '#fff', 
                  border: '2px solid #1e4ed8', 
                  borderRadius: '12px', 
                  padding: '0.75em 2em', 
                  cursor: 'pointer', 
                  fontWeight: 600,
                  fontSize: '0.95em',
                  transition: 'all 0.2s ease',
                  boxShadow: '0 2px 8px rgba(37,99,235,0.3)'
                } as any} 
                onMouseEnter={(e) => e.currentTarget.style.background = '#1e4ed8'}
                onMouseLeave={(e) => e.currentTarget.style.background = '#2563eb'}
                onClick={saveNewUser}>Сохранить</button>
                <button style={{ 
                  background: '#f3f4f6', 
                  color: '#213547', 
                  border: '2px solid #d1d5db', 
                  borderRadius: '12px', 
                  padding: '0.75em 2em', 
                  cursor: 'pointer', 
                  fontWeight: 600,
                  fontSize: '0.95em',
                  transition: 'all 0.2s ease',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                } as any} 
                onMouseEnter={(e) => { e.currentTarget.style.background = '#e5e7eb'; e.currentTarget.style.borderColor = '#9ca3af'; }}
                onMouseLeave={(e) => { e.currentTarget.style.background = '#f3f4f6'; e.currentTarget.style.borderColor = '#d1d5db'; }}
                onClick={() => setShowAddModal(false)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
        
        {/* Модальное окно создания группы */}
        <Show when={showAddGroupModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(0,0,0,0.3)', zIndex: 1000 }} onClick={() => setShowAddGroupModal(false)}>
            <div style={{ 
              background: 'var(--bg-primary)', 
              borderRadius: '32px', 
              border: '6px solid #10b981',
              padding: '2.5rem 2.5rem 2.5rem 2.5rem', 
              width: '450px', 
              maxHeight: '90vh', 
              overflowY: 'auto', 
              overflowX: 'hidden', 
              boxShadow: '0 20px 60px rgba(16,185,129,0.35)', 
              position: 'fixed',
              top: 'calc(50% - 100px)',
              left: '2rem',
              transform: 'translateY(-50%)',
              zIndex: 1001,
              animation: 'slideInFromLeft 0.3s ease-out'
            } as any} class="modal-scroll modal-window" onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.5rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#10b981', transition: 'color 0.2s', zIndex: 10 }} onClick={() => setShowAddGroupModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.5em', fontSize: '1.5em', fontWeight: 700 }}>Создать группу</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Название группы:</label><br />
                <input type="text" value={groupName()} onInput={e => setGroupName(e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} placeholder="Например: 10А" />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Куратор (учитель, необязательно):</label><br />
                <select value={groupCuratorId()?.toString() || ''} onInput={e => setGroupCuratorId(e.currentTarget.value ? parseInt(e.currentTarget.value) : null)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }}>
                  <option value="">Не выбрано</option>
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
              <div style={{ display: 'flex', gap: '0.75em', marginTop: '3.5em', justifyContent: 'flex-end', flexWrap: 'wrap' }}>
                <button style={{ 
                  background: '#10b981', 
                  color: '#fff', 
                  border: '2px solid #059669', 
                  borderRadius: '12px', 
                  padding: '0.75em 2em', 
                  cursor: 'pointer', 
                  fontWeight: 600,
                  fontSize: '0.95em',
                  transition: 'all 0.2s ease',
                  boxShadow: '0 2px 8px rgba(16,185,129,0.3)'
                } as any} 
                onMouseEnter={(e) => e.currentTarget.style.background = '#059669'}
                onMouseLeave={(e) => e.currentTarget.style.background = '#10b981'}
                onClick={saveNewGroup}>Создать</button>
                <button style={{ 
                  background: '#f3f4f6', 
                  color: '#213547', 
                  border: '2px solid #d1d5db', 
                  borderRadius: '12px', 
                  padding: '0.75em 2em', 
                  cursor: 'pointer', 
                  fontWeight: 600,
                  fontSize: '0.95em',
                  transition: 'all 0.2s ease',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                } as any} 
                onMouseEnter={(e) => { e.currentTarget.style.background = '#e5e7eb'; e.currentTarget.style.borderColor = '#9ca3af'; }}
                onMouseLeave={(e) => { e.currentTarget.style.background = '#f3f4f6'; e.currentTarget.style.borderColor = '#d1d5db'; }}
                onClick={() => { setShowAddGroupModal(false); setGroupStudentIds([]); }}>Закрыть</button>
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
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(0,0,0,0.3)', zIndex: 1000 }} onClick={() => setShowAddSubjectModal(false)}>
            <div style={{ 
              background: 'var(--bg-primary)', 
              borderRadius: '32px', 
              border: '6px solid #8b5cf6',
              padding: '2.5rem 2.5rem 2.5rem 2.5rem', 
              width: '450px', 
              maxHeight: '90vh', 
              overflowY: 'auto', 
              overflowX: 'hidden', 
              boxShadow: '0 20px 60px rgba(139,92,246,0.35)', 
              position: 'fixed',
              top: 'calc(50% - 100px)',
              left: '2rem',
              transform: 'translateY(-50%)',
              zIndex: 1001,
              animation: 'slideInFromLeft 0.3s ease-out'
            } as any} class="modal-scroll modal-window" onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.5rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#8b5cf6', transition: 'color 0.2s', zIndex: 10 }} onClick={() => setShowAddSubjectModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.5em', fontSize: '1.5em', fontWeight: 700 }}>Создать предмет</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Название предмета:</label><br />
                <input type="text" value={subjectName()} onInput={e => setSubjectName(e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em', borderRadius: '6px', border: '1px solid #d1d5db' }} placeholder="Например: Математика" />
              </div>
              <div style={{ display: 'flex', gap: '0.75em', marginTop: '3.5em', justifyContent: 'flex-end', flexWrap: 'wrap' }}>
                <button style={{ 
                  background: '#8b5cf6', 
                  color: '#fff', 
                  border: '2px solid #7c3aed', 
                  borderRadius: '12px', 
                  padding: '0.75em 2em', 
                  cursor: 'pointer', 
                  fontWeight: 600,
                  fontSize: '0.95em',
                  transition: 'all 0.2s ease',
                  boxShadow: '0 2px 8px rgba(139,92,246,0.3)'
                } as any} 
                onMouseEnter={(e) => e.currentTarget.style.background = '#7c3aed'}
                onMouseLeave={(e) => e.currentTarget.style.background = '#8b5cf6'}
                onClick={saveNewSubject}>Создать</button>
                <button style={{ 
                  background: '#f3f4f6', 
                  color: '#213547', 
                  border: '2px solid #d1d5db', 
                  borderRadius: '12px', 
                  padding: '0.75em 2em', 
                  cursor: 'pointer', 
                  fontWeight: 600,
                  fontSize: '0.95em',
                  transition: 'all 0.2s ease',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                } as any} 
                onMouseEnter={(e) => { e.currentTarget.style.background = '#e5e7eb'; e.currentTarget.style.borderColor = '#9ca3af'; }}
                onMouseLeave={(e) => { e.currentTarget.style.background = '#f3f4f6'; e.currentTarget.style.borderColor = '#d1d5db'; }}
                onClick={() => setShowAddSubjectModal(false)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
      </div>
      {/* Стили скроллбара и анимации для модальных окон */}
      <style>{`
        @keyframes slideInFromLeft {
          from {
            opacity: 0;
            transform: translateY(-50%) translateX(-100%);
          }
          to {
            opacity: 1;
            transform: translateY(-50%) translateX(0);
          }
        }
        
        .modal-window {
          transition: transform 0.3s ease-out, opacity 0.3s ease-out;
        }
        
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
        
        @media (max-width: 768px) {
          .modal-window {
            width: calc(100vw - 2rem) !important;
            left: 1rem !important;
            max-width: 450px;
          }
        }
      `}</style>
    </>
  );
};

export default AdminDashboard;
