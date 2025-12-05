import type { Component } from 'solid-js';
import Header from '../components/Header';

const Home: Component = () => {
  return (
    <>
      <Header />
      <div class="main-shell" style={{
        'max-width': '900px',
        margin: '4.5rem auto 0 auto',
        background: 'var(--bg-secondary, #f8f9fa)',
        'border-radius': '24px',
        'box-shadow': '0 4px 24px 0 rgba(0,0,0,0.07)',
        padding: '2.5rem 2rem',
        'text-align': 'center',
        'min-height': '60vh',
        'display': 'flex',
        'flex-direction': 'column',
        'align-items': 'center',
        'gap': '2.5rem'
      }}>
        <h1 style={{
          'font-family': 'TT Hoves Pro Trial, a_futuricanord, Segoe UI, Arial, sans-serif',
          'font-size': '2.7rem',
          'font-weight': 700,
          color: '#2563eb',
          margin: 0,
          'letter-spacing': '0.01em',
          'line-height': 1.15
        }}>
          Learn Hub — образовательная платформа нового поколения
        </h1>
        <h2 style={{
          'font-size': '1.35rem',
          color: '#6c757d',
          'font-weight': 400,
          margin: 0,
          'letter-spacing': '0.01em',
          'line-height': 1.3
        }}>
          Для учеников, учителей, родителей и администраторов
        </h2>
        <div style={{
          'font-size': '1.1rem',
          color: '#333',
          'max-width': '700px',
          'margin': '0 auto',
          'line-height': 1.7
        }}>
          <p>
            <b>Learn Hub</b> — это современная образовательная платформа, объединяющая всех участников учебного процесса. Здесь вы можете:
          </p>
          <ul style={{'text-align': 'left', 'margin': '1.5rem auto', 'max-width': '500px', 'font-size': '1.05rem', color: '#444', 'line-height': 1.6}}>
            <li>Загружать, просматривать и скачивать учебные материалы любого формата (PDF, видео, презентации, аудио, документы и др.)</li>
            <li>Выполнять интерактивные задания и тесты</li>
            <li>Следить за своей успеваемостью и прогрессом</li>
            <li>Пользоваться удобной системой фильтров и поиска</li>
            <li>Работать с любого устройства: ПК, планшета или смартфона</li>
            <li>Учителям и администраторам — управлять материалами, заданиями и статистикой</li>
            <li>Родителям — отслеживать успехи своих детей</li>
          </ul>
          <p>
            Платформа создана для того, чтобы сделать обучение удобным, доступным и интересным для всех!
          </p>
        </div>
        <div style={{
          display: 'grid',
          'grid-template-columns': 'repeat(auto-fit, minmax(180px, 1fr))',
          gap: '1.5rem',
          width: '100%',
          'max-width': '700px',
          'margin': '0 auto'
        }}>
          <div style={{'background': '#e3eafc', 'border-radius': '16px', padding: '1.2rem'}}>
            <div style={{'font-size': '2.2rem'}}>📚</div>
            <div style={{'font-weight': 600, 'margin-top': '0.5rem'}}>Все материалы в одном месте</div>
          </div>
          <div style={{'background': '#e3eafc', 'border-radius': '16px', padding: '1.2rem'}}>
            <div style={{'font-size': '2.2rem'}}>🧑‍🏫</div>
            <div style={{'font-weight': 600, 'margin-top': '0.5rem'}}>Интерактивные задания</div>
          </div>
          <div style={{'background': '#e3eafc', 'border-radius': '16px', padding: '1.2rem'}}>
            <div style={{'font-size': '2.2rem'}}>📈</div>
            <div style={{'font-weight': 600, 'margin-top': '0.5rem'}}>Статистика и прогресс</div>
          </div>
          <div style={{'background': '#e3eafc', 'border-radius': '16px', padding: '1.2rem'}}>
            <div style={{'font-size': '2.2rem'}}>🔍</div>
            <div style={{'font-weight': 600, 'margin-top': '0.5rem'}}>Умный поиск и фильтры</div>
          </div>
          <div style={{'background': '#e3eafc', 'border-radius': '16px', padding: '1.2rem'}}>
            <div style={{'font-size': '2.2rem'}}>💻</div>
            <div style={{'font-weight': 600, 'margin-top': '0.5rem'}}>Доступ с любого устройства</div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Home; 