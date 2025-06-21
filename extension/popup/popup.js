// options.js
document.addEventListener('DOMContentLoaded', () => {
  applyLocalization();
  restoreOptions();
});
document.getElementById('saveButton').addEventListener('click', saveOptions);

function applyLocalization() {
  document.title = chrome.i18n.getMessage('optionsPageTitle');
  document.getElementById('optionsPageTitle').textContent = chrome.i18n.getMessage('optionsPageTitle');
  document.getElementById('pageHeading').textContent = chrome.i18n.getMessage('optionsPageTitle');
  document.getElementById('serverAddressLabel').textContent = chrome.i18n.getMessage('serverAddressLabel');
  document.getElementById('serverAddress').placeholder = chrome.i18n.getMessage('serverAddressPlaceholder');
  document.getElementById('saveButton').textContent = chrome.i18n.getMessage('saveButtonText');
}

function saveOptions() {
  const serverAddress = document.getElementById('serverAddress').value;
  const statusElement = document.getElementById('status');

  try {
    new URL(serverAddress);
  } catch (e) {
    statusElement.textContent = chrome.i18n.getMessage('invalidUrlError');
    statusElement.classList.add('error');
    return;
  }

  chrome.storage.sync.set({
    serverAddress: serverAddress
  }, () => {
    statusElement.textContent = chrome.i18n.getMessage('settingsSaved');
    statusElement.classList.remove('error');
    setTimeout(() => {
      statusElement.textContent = '';
    }, 2000);
  });
}

function restoreOptions() {
  chrome.storage.sync.get({
    serverAddress: 'http://localhost:3000'
  }, (items) => {
    document.getElementById('serverAddress').value = items.serverAddress;
  });
}
