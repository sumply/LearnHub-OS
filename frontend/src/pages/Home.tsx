import type { Component } from 'solid-js';
import Header from '../components/Header';

const Home: Component = () => {
  return (
    <>
      <Header />
      <div class="main-shell" style={{
        background: 'var(--bg-secondary, #f8f9fa)',
        'text-align': 'center',
        'min-height': '60vh',
        'display': 'flex',
        'flex-direction': 'column',
        'align-items': 'center',
        'gap': 'var(--spacing-xl)'
      }}>
        <h1 style={{
          'font-family': 'TT Hoves Pro Trial, a_futuricanord, Segoe UI, Arial, sans-serif',
          'font-size': 'var(--font-size-3xl)',
          'font-weight': 700,
          color: '#2563eb',
          margin: 0,
          'letter-spacing': '0.01em',
          'line-height': 1.15
        }}>
          Learn Hub — образовательная платформа нового поколения
        </h1>
        <h2 style={{
          'font-size': 'var(--font-size-lg)',
          color: 'var(--text-secondary)',
          'font-weight': 400,
          margin: 0,
          'letter-spacing': '0.01em',
          'line-height': 1.3
        }}>
          Для учеников, учителей, родителей и администраторов
        </h2>
        <div style={{
          'font-size': 'var(--font-size-base)',
          color: 'var(--text-primary)',
          'max-width': 'var(--container-md)',
          'margin': '0 auto',
          'line-height': 1.7
        }}>
          <p>
            <b>Learn Hub</b> — это современная образовательная платформа, объединяющая всех участников учебного процесса. Здесь вы можете:
          </p>
          <ul style={{'text-align': 'left', 'margin': 'var(--spacing-lg) auto', 'max-width': '500px', 'font-size': 'var(--font-size-base)', color: 'var(--text-primary)', 'line-height': 1.6}}>
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
        <div class="grid-responsive" style={{
          width: '100%',
          'max-width': 'var(--container-md)',
          'margin': '0 auto'
        }}>
          <div class="card" style={{color: 'var(--text-primary)'}}>
            <div style={{'font-size': 'var(--font-size-2xl)'}}>📚</div>
            <div style={{'font-weight': 600, 'margin-top': 'var(--spacing-xs)'}}>Все материалы в одном месте</div>
          </div>
          <div class="card" style={{color: 'var(--text-primary)'}}>
            <div style={{'font-size': 'var(--font-size-2xl)'}}>🧑‍🏫</div>
            <div style={{'font-weight': 600, 'margin-top': 'var(--spacing-xs)'}}>Интерактивные задания</div>
          </div>
          <div class="card" style={{color: 'var(--text-primary)'}}>
            <div style={{'font-size': 'var(--font-size-2xl)'}}>📈</div>
            <div style={{'font-weight': 600, 'margin-top': 'var(--spacing-xs)'}}>Статистика и прогресс</div>
          </div>
          <div class="card" style={{color: 'var(--text-primary)'}}>
            <div style={{'font-size': 'var(--font-size-2xl)'}}>🔍</div>
            <div style={{'font-weight': 600, 'margin-top': 'var(--spacing-xs)'}}>Умный поиск и фильтры</div>
          </div>
          <div class="card" style={{color: 'var(--text-primary)'}}>
            <div style={{'font-size': 'var(--font-size-2xl)'}}>💻</div>
            <div style={{'font-weight': 600, 'margin-top': 'var(--spacing-xs)'}}>Доступ с любого устройства</div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Home; 