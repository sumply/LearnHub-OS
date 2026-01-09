import { type Component, createSignal, onMount, For, Show, createMemo } from 'solid-js';
import { materials, formatFileSize } from '../config/materials';
import type { Material, MaterialType } from '../config/materials';
import { subjects } from '../config/subjectsGroups';
import { getCurrentUser } from '../utils/api';
import Header from '../components/Header';
import AddMaterialForm from '../components/AddMaterialForm';

const Library: Component = () => {

  const [allMaterials, setAllMaterials] = createSignal<Material[]>([]);
  const [error, setError] = createSignal<string | null>(null);
  // Удалено: const [loading, setLoading] = createSignal(true);
  
  // Фильтры
  const [searchQuery, setSearchQuery] = createSignal('');
  const [selectedSubject, setSelectedSubject] = createSignal<string>('');
  const [selectedGrade, setSelectedGrade] = createSignal<string>('');
  const [selectedType, setSelectedType] = createSignal<MaterialType | ''>('');
  
  // Состояние для управления материалами
  const [showAddForm, setShowAddForm] = createSignal(false);
  const [editingMaterial, setEditingMaterial] = createSignal<Material | null>(null);

  // Получение текущего пользователя
  const currentUser = createMemo(() => getCurrentUser() as any);
  const isAdmin = createMemo(() => currentUser()?.role === 'admin');
  const isTeacher = createMemo(() => currentUser()?.role === 'teacher');
  const canEdit = createMemo(() => isAdmin() || isTeacher());

  // Фильтрация материалов
  const filteredMaterials = createMemo(() => {
    let filtered = allMaterials();
    // Только фильтры поиска, предмета, класса, типа
    if (searchQuery()) {
      const query = searchQuery().toLowerCase();
      filtered = filtered.filter(m => 
        m.title.toLowerCase().includes(query) ||
        m.description.toLowerCase().includes(query) ||
        (m.tags && m.tags.some(tag => tag.toLowerCase().includes(query)))
      );
    }
    if (selectedSubject()) {
      filtered = filtered.filter(m => m.category === selectedSubject());
    }
    if (selectedGrade()) {
      filtered = filtered.filter(m => m.grade === selectedGrade());
    }
    if (selectedType()) {
      filtered = filtered.filter(m => m.type === selectedType());
    }
    return filtered;
  });

  // Получение уникальных классов
  const availableGrades = createMemo(() => {
    const grades = new Set(allMaterials().map(m => m.grade).filter(Boolean));
    return Array.from(grades).sort();
  });

  // Получение уникальных типов
  const availableTypes = createMemo(() => {
    const types = new Set(allMaterials().map(m => m.type));
    return Array.from(types).sort();
  });

  onMount(() => {
    try {
      // Преобразуем uploadDate в объект Date
      const fixedMaterials = materials.map(m => ({
        ...m,
        uploadDate: typeof m.uploadDate === 'string' ? new Date(m.uploadDate) : m.uploadDate
      }));
      setAllMaterials(fixedMaterials);
    } catch (e: any) {
      setError(e.message || 'Ошибка загрузки материалов');
      // eslint-disable-next-line no-console
      console.error('Ошибка при загрузке материалов:', e);
    }
  });

  // Функции управления материалами
  const handleDeleteMaterial = (materialId: string) => {
    if (confirm('Вы уверены, что хотите удалить этот материал?')) {
      setAllMaterials(prev => prev.filter(m => m.id !== materialId));
    }
  };

  const handleAddMaterial = (materialData: Omit<Material, 'id' | 'uploadDate' | 'uploadedBy'>) => {
    const newMaterial: Material = {
      ...materialData,
      id: `mat${Date.now()}`,
      uploadDate: new Date(),
      uploadedBy: currentUser()?.id || 'unknown'
    };
    setAllMaterials(prev => [newMaterial, ...prev]);
    setShowAddForm(false);
  };

  const handleDownload = (material: Material) => {
    const link = document.createElement('a');
    link.href = material.url;
    link.download = material.fileName || material.title;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handleView = (material: Material) => {
    window.open(material.url, '_blank');
  };

  const getTypeIcon = (type: MaterialType) => {
    switch (type) {
      case 'pdf': return '📄';
      case 'video': return '🎥';
      case 'audio': return '🎵';
      case 'presentation': return '📊';
      case 'document': return '📝';
      case 'image': return '🖼️';
      case 'archive': return '📦';
      case 'link': return '🔗';
      default: return '📄';
    }
  };

  const getTypeName = (type: MaterialType) => {
    switch (type) {
      case 'pdf': return 'PDF';
      case 'video': return 'Видео';
      case 'audio': return 'Аудио';
      case 'presentation': return 'Презентация';
      case 'document': return 'Документ';
      case 'image': return 'Изображение';
      case 'archive': return 'Архив';
      case 'link': return 'Ссылка';
      default: return 'Файл';
    }
  };


  if (error()) {
    return <div style={{'text-align': 'center', 'margin-top': '2rem', color: '#e76f51'}}>{error()}</div>;
  }

  return (
    <>
      <Header />
      <div class="main-shell">
        <div style={{
          'display': 'flex',
          'justify-content': 'space-between',
          'align-items': 'center',
          'margin-bottom': '2rem',
          'flex-wrap': 'wrap',
          'gap': '1rem'
        }}>
          <h2 style={{
            'font-family': 'TT Hoves Pro Trial, sans-serif',
            'margin': '0'
          }}>Библиотека материалов</h2>
          
          <Show when={canEdit()}>
            <button 
              onClick={() => setShowAddForm(true)}
              style={{
                'background': '#264653',
                'color': 'white',
                'border': 'none',
                'padding': '0.75rem 1.5rem',
                'border-radius': '8px',
                'cursor': 'pointer',
                'font-weight': '500'
              }}
            >
              + Добавить материал
            </button>
          </Show>
        </div>

        {/* Фильтры */}
        <div style={{
          'background': 'var(--bg-secondary)',
          'padding': '1.5rem',
          'border-radius': '12px',
          'margin-bottom': '2rem',
          'box-shadow': '0 2px 8px rgba(0,0,0,0.04)',
          color: 'var(--text-primary)'
        }}>
          <h3 style={{
            'font-family': 'TT Hoves Pro Trial, sans-serif',
            'margin-top': '0',
            'margin-bottom': '1rem',
            'font-size': '1.1rem',
            color: 'var(--text-primary)'
          }}>Фильтры</h3>
          
          <div style={{
            'display': 'grid',
            'grid-template-columns': 'repeat(auto-fit, minmax(200px, 1fr))',
            'gap': '1rem'
          }}>
            {/* Поиск */}
            <div>
              <label style={{'display': 'block', 'margin-bottom': '0.5rem', 'font-weight': '500'}}>
                Поиск по названию
              </label>
              <input
                type="text"
                placeholder="Введите название материала..."
                value={searchQuery()}
                onInput={(e) => setSearchQuery(e.currentTarget.value)}
                style={{
                  'width': '100%',
                  'padding': '0.75rem',
                  'border': '1px solid #ddd',
                  'border-radius': '6px',
                  'font-size': '0.9rem'
                }}
              />
            </div>

            {/* Фильтр по предмету */}
            <div>
              <label style={{'display': 'block', 'margin-bottom': '0.5rem', 'font-weight': '500'}}>
                Предмет
              </label>
              <select
                value={selectedSubject()}
                onChange={(e) => setSelectedSubject(e.currentTarget.value)}
                style={{
                  'width': '100%',
                  'padding': '0.75rem',
                  'border': '1px solid #ddd',
                  'border-radius': '6px',
                  'font-size': '0.9rem'
                }}
              >
                <option value="">Все предметы</option>
                <For each={subjects}>
                  {subject => <option value={subject.id}>{subject.name}</option>}
                </For>
              </select>
            </div>

            {/* Фильтр по классу */}
            <div>
              <label style={{'display': 'block', 'margin-bottom': '0.5rem', 'font-weight': '500'}}>
                Класс
              </label>
              <select
                value={selectedGrade()}
                onChange={(e) => setSelectedGrade(e.currentTarget.value)}
                style={{
                  'width': '100%',
                  'padding': '0.75rem',
                  'border': '1px solid #ddd',
                  'border-radius': '6px',
                  'font-size': '0.9rem'
                }}
              >
                <option value="">Все классы</option>
                <For each={availableGrades()}>
                  {grade => <option value={grade}>{grade}</option>}
                </For>
              </select>
            </div>

            {/* Фильтр по типу */}
            <div>
              <label style={{'display': 'block', 'margin-bottom': '0.5rem', 'font-weight': '500'}}>
                Тип материала
              </label>
              <select
                value={selectedType()}
                onChange={(e) => setSelectedType(e.currentTarget.value as MaterialType | '')}
                style={{
                  'width': '100%',
                  'padding': '0.75rem',
                  'border': '1px solid #ddd',
                  'border-radius': '6px',
                  'font-size': '0.9rem'
                }}
              >
                <option value="">Все типы</option>
                <For each={availableTypes()}>
                  {type => <option value={type}>{getTypeName(type)}</option>}
                </For>
              </select>
            </div>
          </div>

          {/* Сброс фильтров */}
          <button
            onClick={() => {
              setSearchQuery('');
              setSelectedSubject('');
              setSelectedGrade('');
              setSelectedType('');
            }}
            style={{
              'background': 'none',
              'border': '1px solid #ddd',
              'padding': '0.5rem 1rem',
              'border-radius': '6px',
              'cursor': 'pointer',
              'margin-top': '1rem',
              'font-size': '0.9rem'
            }}
          >
            Сбросить фильтры
          </button>
        </div>

        {/* Список материалов */}
        <Show when={filteredMaterials().length > 0} fallback={
          <div style={{
            'text-align': 'center',
            'padding': '3rem',
            'color': '#666',
            'background': '#f8f9fa',
            'border-radius': '12px'
          }}>
            <div style={{'font-size': '3rem', 'margin-bottom': '1rem'}}>📚</div>
            <h3>Материалы не найдены</h3>
            <p>Попробуйте изменить параметры поиска или фильтры</p>
          </div>
        }>
          <div style={{
            'display': 'grid',
            'grid-template-columns': 'repeat(auto-fill, minmax(350px, 1fr))',
            'gap': '3rem'
          }}>
            <For each={filteredMaterials()}>
              {material => (
                <div style={{
                  'background': 'white',
                  'border-radius': '12px',
                  'padding': '1.5rem',
                  'box-shadow': '0 4px 12px rgba(0,0,0,0.08)',
                  'border': '1px solid #e9ecef',
                  'transition': 'transform 0.2s, box-shadow 0.2s'
                }}>
                  {/* Заголовок с иконкой */}
                  <div style={{
                    'display': 'flex',
                    'align-items': 'center',
                    'margin-bottom': '1rem'
                  }}>
                    <span style={{'font-size': '1.5rem', 'margin-right': '0.75rem'}}>
                      {getTypeIcon(material.type)}
                    </span>
                    <div>
                      <h3 style={{
                        'margin': '0',
                        'font-size': '1.1rem',
                        'font-weight': '600',
                        'color': '#264653'
                      }}>
                        {material.title}
                      </h3>
                      <div style={{
                        'font-size': '0.8rem',
                        'color': '#6c757d',
                        'margin-top': '0.25rem'
                      }}>
                        {getTypeName(material.type)} • {material.fileSize ? formatFileSize(material.fileSize) : 'Ссылка'}
                      </div>
                    </div>
                  </div>

                  {/* Описание */}
                  <p style={{
                    'color': '#495057',
                    'margin': '0 0 1rem 0',
                    'line-height': '1.5',
                    'font-size': '0.9rem'
                  }}>
                    {material.description}
                  </p>

                  {/* Метаданные */}
                  <div style={{
                    'display': 'grid',
                    'grid-template-columns': '1fr 1fr',
                    'gap': '0.5rem',
                    'margin-bottom': '1rem',
                    'font-size': '0.8rem',
                    'color': '#6c757d'
                  }}>
                    <div><strong>Предмет:</strong> {subjects.find(s => s.id === material.category)?.name || material.category}</div>
                    <div><strong>Класс:</strong> {material.grade || 'Не указан'}</div>
                    <div><strong>Загружен:</strong> {material.uploadDate.toLocaleDateString('ru-RU')}</div>
                    <div><strong>Тип:</strong> {getTypeName(material.type)}</div>
                  </div>

                  {/* Теги */}
                  <Show when={material.tags && material.tags.length > 0}>
                    <div style={{
                      'display': 'flex',
                      'flex-wrap': 'wrap',
                      'gap': '0.5rem',
                      'margin-bottom': '1rem'
                    }}>
                      <For each={material.tags}>
                        {tag => (
                          <span style={{
                            'background': '#e9ecef',
                            'color': '#495057',
                            'padding': '0.25rem 0.5rem',
                            'border-radius': '12px',
                            'font-size': '0.75rem'
                          }}>
                            {tag}
                          </span>
                        )}
                      </For>
                    </div>
                  </Show>

                  {/* Действия */}
                  <div style={{
                    'display': 'flex',
                    'gap': '0.5rem',
                    'flex-wrap': 'wrap'
                  }}>
                    <button
                      onClick={() => handleView(material)}
                      style={{
                        'background': '#264653',
                        'color': 'white',
                        'border': 'none',
                        'padding': '0.5rem 1rem',
                        'border-radius': '6px',
                        'cursor': 'pointer',
                        'font-size': '0.85rem',
                        'flex': '1'
                      }}
                    >
                      👁️ Просмотр
                    </button>
                    
                    <Show when={material.fileName}>
                      <button
                        onClick={() => handleDownload(material)}
                        style={{
                          'background': '#2a9d8f',
                          'color': 'white',
                          'border': 'none',
                          'padding': '0.5rem 1rem',
                          'border-radius': '6px',
                          'cursor': 'pointer',
                          'font-size': '0.85rem',
                          'flex': '1'
                        }}
                      >
                        ⬇️ Скачать
                      </button>
                    </Show>

                    <Show when={canEdit()}>
                      <button
                        onClick={() => handleDeleteMaterial(material.id)}
                        style={{
                          'background': '#e76f51',
                          'color': 'white',
                          'border': 'none',
                          'padding': '0.5rem 1rem',
                          'border-radius': '6px',
                          'cursor': 'pointer',
                          'font-size': '0.85rem'
                        }}
                      >
                        🗑️ Удалить
                      </button>
                    </Show>
                  </div>
                </div>
              )}
            </For>
          </div>
        </Show>

        {/* Статистика */}
        <div style={{
          'margin-top': '2rem',
          'padding': '1rem',
          'background': '#f8f9fa',
          'border-radius': '8px',
          'text-align': 'center',
          'color': '#6c757d',
          'font-size': '0.9rem'
        }}>
          Показано {filteredMaterials().length} из {allMaterials().length} материалов
        </div>
      </div>

      {/* Модальное окно добавления материала */}
      <Show when={showAddForm()}>
        <div style={{ 'background': 'var(--bg-primary)', 'padding': '2rem', 'border-radius': '12px', 'box-shadow': '0 2px 8px var(--shadow-heavy)', 'margin-bottom': '2rem', 'color': 'var(--text-primary)' }}>
          <AddMaterialForm onMaterialAdded={handleAddMaterial} />
          <button onClick={() => setShowAddForm(false)} style={{ 'margin-top': '1rem', 'background': 'var(--bg-secondary)', 'color': 'var(--text-primary)', 'border': '1px solid var(--border-color)' }}>Отмена</button>
        </div>
      </Show>
    </>
  );
};

export default Library; 