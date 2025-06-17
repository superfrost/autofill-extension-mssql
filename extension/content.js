// Конфигурация
const CONFIG = {
  INPUT_CLASS: 'tracked-input', // Класс отслеживаемого поля
  DEBOUNCE_DELAY: 500,          // Задержка debounce (мс)
  MIN_CHARS: 3                  // Минимальное кол-во символов
};

// Глобальные переменные
let currentController = null;
let dropdown = null;
let apiBaseUrl = 'http://localhost';
let apiPort = '8070';

// Инициализация
async function initAutocomplete() {
  // Загрузка сохраненных настроек
  chrome.storage.sync.get(['apiBaseUrl', 'apiPort'], (data) => {
    if (data.apiBaseUrl) apiBaseUrl = data.apiBaseUrl;
    if (data.apiPort) apiPort = data.apiPort;
  });
  
  // Слушатель для обновления настроек
  chrome.runtime.onMessage.addListener((message) => {
    if (message.type === "API_SETTINGS_UPDATE") {
      apiBaseUrl = message.baseUrl;
      apiPort = message.port;
    }
  });

  const input = document.querySelector(`.${CONFIG.INPUT_CLASS}`);
  if (!input) return;
  
  createDropdown();
  setupInputListeners(input);
}

// Настройка обработчиков
function setupInputListeners(input) {
  let timeoutId = null;
  
  input.addEventListener('input', () => {
    const query = input.value.trim();
    
    // Скрыть список если символов < 3
    if (query.length < CONFIG.MIN_CHARS) {
      hideDropdown();
      return;
    }
    
    // Применить debounce
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => {
      fetchSuggestions(query);
    }, CONFIG.DEBOUNCE_DELAY);
  });
  
  // Скрыть список при клике вне
  document.addEventListener('click', (e) => {
    if (!dropdown.contains(e.target)) hideDropdown();
  });
}

// Обновленная функция запроса к API
async function fetchSuggestions(query) {
  if (currentController) currentController.abort();
  currentController = new AbortController();
  
  const apiUrl = `${apiBaseUrl}:${apiPort}/list`;
  
  try {
    const response = await fetch(apiUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
      signal: currentController.signal
    });
    
    const data = await response.json();
    showSuggestions(data);
  } catch (err) {
    if (err.name !== 'AbortError') {
      console.error('API error:', err);
    }
  }
}

// Отображение предложений
function showSuggestions(items) {
  if (!items || items.length === 0) {
    hideDropdown();
    return;
  }
  
  dropdown.innerHTML = '';
  items.forEach(item => {
    const option = document.createElement('div');
    option.className = 'dropdown-item';
    option.textContent = item.displayText; // Используйте нужное поле из API
    option.dataset.value = JSON.stringify(item);
    
    option.addEventListener('click', () => {
      fillFormFields(item);
      hideDropdown();
    });
    
    dropdown.appendChild(option);
  });
  
  // Позиционирование
  const input = document.querySelector(`.${CONFIG.INPUT_CLASS}`);
  const rect = input.getBoundingClientRect();
  dropdown.style.display = 'block';
  dropdown.style.top = `${rect.bottom + window.scrollY}px`;
  dropdown.style.left = `${rect.left + window.scrollX}px`;
  dropdown.style.width = `${rect.width}px`;
}

// Заполнение полей формы
function fillFormFields(data) {
  // Пример (настройте под вашу структуру):
  const fields = {
    '.field-email': data.email,
    '.field-name': data.name,
    '.field-phone': data.phone,
    '.field-address': data.address,
    '.field-comment': data.comment
  };
  
  Object.entries(fields).forEach(([selector, value]) => {
    const field = document.querySelector(selector);
    if (field) field.value = value || '';
  });
}

// Скрыть список
function hideDropdown() {
  if (dropdown) dropdown.style.display = 'none';
}

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', initAutocomplete);
