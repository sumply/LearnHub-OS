import { createSignal, For, Show } from 'solid-js';
import { users } from '../config/users';
import { subjects as allSubjects } from '../config/subjectsGroups';
import { addUser, getAllUsers, findUserByEmail, isEmailTaken } from '../config/users';
import Header from '../components/Header';
import { A } from '@solidjs/router';

const tableStyle = {
  width: '100%',
  borderCollapse: 'collapse',
  background: '#fff',
  fontSize: '1em',
};
const thStyle = {
  border: '1px solid #d1d5db',
  padding: '0.5em 0.8em',
  textAlign: 'left',
  background: '#f3f4f6',
  fontWeight: 600,
};
const tdStyle = {
  border: '1px solid #d1d5db',
  padding: '0.5em 0.8em',
  textAlign: 'left',
};

const AdminDashboard = () => {
  // Получаем всех студентов
  const students = users.filter(u => u.role === 'student');
  // Получаем все уникальные группы
  const groups = Array.from(new Set(students.map(s => s.group)));
  // Сигнал для выбранного студента (для подробного просмотра)
  const [selectedStudent, setSelectedStudent] = createSignal(null);
  const [editUser, setEditUser] = createSignal(null);
  const [editFields, setEditFields] = createSignal({ name: '', surname: '', email: '', password: '', group: '', subjects: [] });
  const [showAddModal, setShowAddModal] = createSignal(false);
  const [addRole, setAddRole] = createSignal('student');
  const [addFields, setAddFields] = createSignal({ name: '', surname: '', email: '', password: '', group: '', subjects: [] });

  function openEditUser(user) {
    setEditUser(user);
    setEditFields({
      name: user.name,
      surname: user.surname,
      email: user.email,
      password: user.password,
      group: user.group || '',
      subjects: user.subjects ? [...user.subjects] : [],
    });
  }

  function handleEditField(field, value) {
    setEditFields({ ...editFields(), [field]: value });
  }

  function handleEditSubjects(subjId) {
    const current = editFields().subjects || [];
    if (current.includes(subjId)) {
      setEditFields({ ...editFields(), subjects: current.filter(s => s !== subjId) });
    } else {
      setEditFields({ ...editFields(), subjects: [...current, subjId] });
    }
  }

  function saveUserEdit() {
    if (!editUser()) return;
    const updated = { ...editUser(), ...editFields() };
    // Обновляем пользователя в localStorage
    const all = getAllUsers().map(u => u.email === updated.email ? updated : u);
    localStorage.setItem('users', JSON.stringify(all));
    setEditUser(null);
    setSelectedStudent(null);
    window.location.reload(); // для простоты, чтобы обновить таблицу
  }

  function deleteUser() {
    if (!editUser()) return;
    const all = getAllUsers().filter(u => u.email !== editUser().email);
    localStorage.setItem('users', JSON.stringify(all));
    setEditUser(null);
    setSelectedStudent(null);
    window.location.reload();
  }

  function handleAddField(field, value) {
    setAddFields({ ...addFields(), [field]: value });
  }
  function handleAddSubjects(subjId) {
    const current = addFields().subjects || [];
    if (current.includes(subjId)) {
      setAddFields({ ...addFields(), subjects: current.filter(s => s !== subjId) });
    } else {
      setAddFields({ ...addFields(), subjects: [...current, subjId] });
    }
  }
  function saveNewUser() {
    const base = {
      id: Date.now().toString(),
      name: addFields().name,
      surname: addFields().surname,
      email: addFields().email,
      password: addFields().password,
      role: addRole(),
    };
    let newUser;
    if (addRole() === 'student') {
      newUser = { ...base, group: addFields().group };
    } else if (addRole() === 'teacher') {
      newUser = { ...base, subjects: addFields().subjects };
    }
    addUser(newUser);
    setShowAddModal(false);
    window.location.reload();
  }

  // Преподаватели
  const teachers = users.filter(u => u.role === 'teacher');
  // Для каждой группы: какие преподаватели ведут предметы в этой группе (по совпадению subject id)
  // Для простоты: покажем все группы, если у преподавателя есть хотя бы один предмет, совпадающий с предметом группы
  // (или просто список всех групп, если нет строгой связи)

  // Получить список групп, где есть студенты
  const studentGroups = Array.from(new Set(users.filter(u => u.role === 'student').map(s => s.group)));

  // Для каждого преподавателя — список групп, где есть студенты (можно доработать логику связи)

  return (
    <>
      <Header />
      <div class="main-shell">
        <h2 style={{ marginBottom: '2rem', color: '#2563eb', fontSize: '2rem', fontWeight: 700, letterSpacing: '0.01em' }}>Админ-панель: Ученики по классам</h2>
        <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '1.5em' }}>
          <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddModal(true)}>Добавить пользователя</button>
        </div>
        <For each={groups}>{(group, idx) => (
          <div style={{ marginBottom: '2.5em' }}>
            <h3 style={{ color: '#213547', marginBottom: '0.7em' }}>Класс: {group}</h3>
            <table style={tableStyle}>
              <thead>
                <tr>
                  <th style={thStyle}>ФИО</th>
                  <th style={thStyle}>Email</th>
                  <th style={thStyle}></th>
                </tr>
              </thead>
              <tbody>
                <For each={students.filter(s => s.group === group)}>{(student, sidx) => (
                  <tr style={sidx % 2 === 1 ? { background: '#f9fafb' } : {}}>
                    <td style={tdStyle}>{student.surname} {student.name}</td>
                    <td style={tdStyle}>{student.email}</td>
                    <td style={tdStyle}>
                      <button style={{ background: 'none', color: '#2563eb', border: '1px solid #2563eb', borderRadius: '3px', padding: '0.2em 0.7em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em' }}>
                        <A href={`/admin/profile/${student.id}`} style={{ color: '#2563eb', textDecoration: 'none' }}>К данным</A>
                      </button>
                    </td>
                  </tr>
                )}</For>
              </tbody>
            </table>
          </div>
        )}</For>
        {/* Модальное окно с подробной инфой о студенте */}
        <Show when={!!selectedStudent()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setSelectedStudent(null)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '10px', padding: '2.2rem 2rem 2rem 2rem', minWidth: '320px', maxWidth: '95vw', boxShadow: '0 8px 32px rgba(37,99,235,0.18)', position: 'relative' }} onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.2rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s' }} onMouseOver={e => (e.currentTarget.style.color = '#e76f51')} onMouseOut={e => (e.currentTarget.style.color = '#2563eb')} onClick={() => setSelectedStudent(null)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.2em' }}>Информация о студенте</h3>
              <div style={{ marginBottom: '0.7em' }}><b>ФИО:</b> {selectedStudent()?.surname} {selectedStudent()?.name}</div>
              <div style={{ marginBottom: '0.7em' }}><b>Email:</b> {selectedStudent()?.email}</div>
              <div style={{ marginBottom: '0.7em' }}><b>Класс:</b> {selectedStudent()?.group}</div>
            </div>
          </div>
        </Show>
        <Show when={!!editUser()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setEditUser(null)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '10px', padding: '2.2rem 2rem 2rem 2rem', minWidth: '340px', maxWidth: '95vw', boxShadow: '0 8px 32px rgba(37,99,235,0.18)', position: 'relative' }} onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.2rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s' }} onClick={() => setEditUser(null)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.2em' }}>Профиль пользователя</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Фамилия:</label><br />
                <input type="text" value={editFields().surname} onInput={e => handleEditField('surname', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Имя:</label><br />
                <input type="text" value={editFields().name} onInput={e => handleEditField('name', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Email:</label><br />
                <input type="email" value={editFields().email} onInput={e => handleEditField('email', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Пароль:</label><br />
                <input type="text" value={editFields().password} onInput={e => handleEditField('password', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <Show when={editUser()?.role === 'student'}>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Класс:</label><br />
                  <input type="text" value={editFields().group} onInput={e => handleEditField('group', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
                </div>
              </Show>
              <Show when={editUser()?.role === 'teacher'}>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Предметы:</label><br />
                  <For each={allSubjects}>{subj => (
                    <label style={{ marginRight: '1em' }}>
                      <input type="checkbox" checked={editFields().subjects.includes(subj.id)} onChange={() => handleEditSubjects(subj.id)} /> {subj.name}
                    </label>
                  )}</For>
                </div>
              </Show>
              <div style={{ display: 'flex', gap: '1em', marginTop: '1.5em' }}>
                <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={saveUserEdit}>Сохранить</button>
                <button style={{ background: '#fff', color: 'var(--accent-danger)', border: '1px solid var(--accent-danger)', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={deleteUser}>Удалить</button>
                <button style={{ background: '#f3f4f6', color: '#213547', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setEditUser(null)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
        <Show when={showAddModal()}>
          <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(37,99,235,0.10)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setShowAddModal(false)}>
            <div style={{ background: 'var(--bg-primary)', borderRadius: '10px', padding: '2.2rem 2rem 2rem 2rem', minWidth: '340px', maxWidth: '95vw', boxShadow: '0 8px 32px rgba(37,99,235,0.18)', position: 'relative' }} onClick={e => e.stopPropagation()}>
              <button style={{ position: 'absolute', top: '1.2rem', right: '1.5rem', fontSize: '1.7em', background: 'none', border: 'none', cursor: 'pointer', color: '#2563eb', transition: 'color 0.2s' }} onClick={() => setShowAddModal(false)} aria-label="Закрыть">&times;</button>
              <h3 style={{ color: '#2563eb', marginBottom: '1.2em' }}>Добавить пользователя</h3>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Роль:</label><br />
                <select value={addRole()} onInput={e => setAddRole(e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }}>
                  <option value="student">Ученик</option>
                  <option value="teacher">Учитель</option>
                </select>
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Фамилия:</label><br />
                <input type="text" value={addFields().surname} onInput={e => handleAddField('surname', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Имя:</label><br />
                <input type="text" value={addFields().name} onInput={e => handleAddField('name', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Email:</label><br />
                <input type="email" value={addFields().email} onInput={e => handleAddField('email', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <div style={{ marginBottom: '0.7em' }}>
                <label>Пароль:</label><br />
                <input type="text" value={addFields().password} onInput={e => handleAddField('password', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
              </div>
              <Show when={addRole() === 'student'}>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Класс:</label><br />
                  <input type="text" value={addFields().group} onInput={e => handleAddField('group', e.currentTarget.value)} style={{ width: '100%', padding: '0.4em', marginBottom: '0.5em' }} />
                </div>
              </Show>
              <Show when={addRole() === 'teacher'}>
                <div style={{ marginBottom: '0.7em' }}>
                  <label>Предметы:</label><br />
                  <For each={allSubjects}>{subj => (
                    <label style={{ marginRight: '1em' }}>
                      <input type="checkbox" checked={addFields().subjects.includes(subj.id)} onChange={() => handleAddSubjects(subj.id)} /> {subj.name}
                    </label>
                  )}</For>
                </div>
              </Show>
              <div style={{ display: 'flex', gap: '1em', marginTop: '1.5em' }}>
                <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={saveNewUser}>Сохранить</button>
                <button style={{ background: '#f3f4f6', color: '#213547', border: 'none', borderRadius: '5px', padding: '0.5em 1.5em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setShowAddModal(false)}>Закрыть</button>
              </div>
            </div>
          </div>
        </Show>
        <hr style={{ margin: '2.5em 0', border: 0, borderTop: '2px solid #e3eafc' }} />
        <h2 style={{ marginBottom: '1.5rem', color: '#2563eb', fontSize: '1.7rem', fontWeight: 700 }}>Преподаватели, их предметы и классы</h2>
        <table style={tableStyle}>
          <thead>
            <tr>
              <th style={thStyle}>ФИО</th>
              <th style={thStyle}>Email</th>
              <th style={thStyle}>Предметы</th>
              <th style={thStyle}>Классы</th>
              <th style={thStyle}></th>
            </tr>
          </thead>
          <tbody>
            <For each={teachers}>{(teacher, tidx) => (
              <tr style={tidx % 2 === 1 ? { background: '#f9fafb' } : {}}>
                <td style={tdStyle}>{teacher.surname} {teacher.name}</td>
                <td style={tdStyle}>{teacher.email}</td>
                <td style={tdStyle}>{(teacher.subjects || []).map(sid => (allSubjects.find(s => s.id === sid)?.name || sid)).join(', ')}</td>
                <td style={tdStyle}>{studentGroups.join(', ')}</td>
                <td style={tdStyle}>
                  <button style={{ background: 'none', color: '#2563eb', border: '1px solid #2563eb', borderRadius: '3px', padding: '0.2em 0.7em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em' }}>
                    <A href={`/admin/profile/${teacher.id}`} style={{ color: '#2563eb', textDecoration: 'none' }}>К данным</A>
                  </button>
                </td>
              </tr>
            )}</For>
          </tbody>
        </table>
      </div>
    </>
  );
};

export default AdminDashboard; 