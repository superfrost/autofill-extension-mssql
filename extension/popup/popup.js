document.addEventListener('DOMContentLoaded', () => {
  const saveBtn = document.getElementById('saveBtn');
  const baseUrlInput = document.getElementById('baseUrl');
  const portInput = document.getElementById('port');
  const statusEl = document.getElementById('status');
  
  // Загрузка сохраненных настроек
  chrome.storage.sync.get(['apiBaseUrl', 'apiPort'], (data) => {
    if (data.apiBaseUrl) baseUrlInput.value = data.apiBaseUrl;
    if (data.apiPort) portInput.value = data.apiPort;
  });
  
  // Сохранение настроек
  saveBtn.addEventListener('click', () => {
    const baseUrl = baseUrlInput.value.trim();
    const port = portInput.value.trim();
    
    if (!baseUrl || !port) {
      showStatus('Заполните оба поля!', true);
      return;
    }
    
    chrome.storage.sync.set({
      apiBaseUrl: baseUrl,
      apiPort: port
    }, () => {
      showStatus('Настройки сохранены!');
      
      // Сообщаем контент-скриптам об обновлении настроек
      chrome.tabs.query({}, (tabs) => {
        for (const tab of tabs) {
          chrome.tabs.sendMessage(tab.id, {
            type: "API_SETTINGS_UPDATE",
            baseUrl,
            port
          });
        }
      });
    });
  });
  
  function showStatus(message, isError = false) {
    statusEl.textContent = message;
    statusEl.style.color = isError ? '#f44336' : '#4CAF50';
    
    setTimeout(() => {
      statusEl.textContent = '';
    }, 3000);
  }
});
