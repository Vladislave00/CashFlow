document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    const usernameElement = document.getElementById('username');
    const accountingTitle = document.getElementById('accounting-title');
    const accountingContent = document.getElementById('accounting-content');
    const backButton = document.getElementById('back-btn');
    const logoutButton = document.getElementById('logout-btn');
    const createAccountModal = new bootstrap.Modal(document.getElementById('createAccountModal'));
    const editAccountModal = new bootstrap.Modal(document.getElementById('editAccountModal'));
    const createAccountForm = document.getElementById('create-account-form');
    const updateAccountForm = document.getElementById('update-account-form');

    let currentEditingId = null; // Для хранения ID редактируемого счета

    if (!token) {
        window.location.href = '/login';
    } else {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            usernameElement.textContent = payload.name;
            const accountingId = new URLSearchParams(window.location.search).get('id');
            if (accountingId) {
                loadAccountingTitle(accountingId);
                loadAccounts(accountingId);
            }
        } catch (error) {
            console.error("Failed to decode token:", error);
        }
    }

    // Обработчик события для кнопки "Назад"
    backButton.addEventListener('click', function() {
        window.location.href = '/';
    });

    // Обработчик события для кнопки "Logout"
    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('token');
        window.location.href = '/login';
    });

    // Обработчик события для создания счета
    createAccountForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const accountingId = new URLSearchParams(window.location.search).get('id');
        createAccount(accountingId);
    });

    // Обработчик события для редактирования счета
    updateAccountForm.addEventListener('submit', function(event) {
        event.preventDefault();
        if (currentEditingId) { // Проверяем, установлен ли ID
            updateAccount(currentEditingId); // Вызываем функцию обновления счета
        } else {
            console.error("No account ID set for updating.");
        }
    });

    // Обработчик события для кнопки "Изменить"
    accountingContent.addEventListener('click', function(event) {
        if (event.target.classList.contains('btn-warning')) { // Проверяем, что кликнули по кнопке "Изменить"
            const accountRow = event.target.closest('.account-row'); // Находим родительский элемент
            const accountId = accountRow.dataset.accountId; // Получаем ID счета из data-атрибута
            const accountName = accountRow.dataset.accountName; // Получаем название счета из data-атрибута
            const accountAmount = accountRow.dataset.accountAmount; // Получаем сумму счета из data-атрибута

            openEditModal(accountId, accountName, accountAmount); // Открываем модальное окно
        } else if (event.target.classList.contains('btn-danger')) { // Проверяем, что кликнули по кнопке "Удалить"
            const accountRow = event.target.closest('.account-row'); // Находим родительский элемент
            const accountId = accountRow.dataset.accountId; // Получаем ID счета из data-атрибута

            if (confirm("Вы уверены, что хотите удалить этот счет?")) {
                deleteAccount(accountId); // Удаляем счет
            }
        }
    });

    function loadAccountingTitle(accountingId) {
        fetch(`http://localhost:8000/api/accountings/${accountingId}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load accounting title');
            }
            return response.json();
        })
        .then(data => {
            accountingTitle.textContent = data.name;
        })
        .catch(error => console.error("Error loading accounting title:", error));
    }

    function loadAccounts(accountingId) {
        // Очищаем контент только один раз в начале
        accountingContent.innerHTML = '<div class="text-center">Загрузка данных...</div>';
    
        // Загрузка общего счета
        fetch(`http://localhost:8000/api/accountings/${accountingId}/accounts/general`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => {
            if (!response.ok) throw new Error('Ошибка загрузки общего счета');
            return response.json();
        })
        .then(generalData => {
            displayGeneralAccount(generalData);
            
            // Затем загружаем остальные счета
            fetch(`http://localhost:8000/api/accountings/${accountingId}/accounts`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
            .then(response => {
                if (!response.ok) throw new Error('Ошибка загрузки счетов');
                return response.json();
            })
            .then(accounts => {
                displayAccounts(accounts);
            })
            .catch(error => {
                accountingContent.innerHTML += `<div class="alert alert-danger mt-3">${error.message}</div>`;
            });
        })
        .catch(error => {
            accountingContent.innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
        });
    }
    
    function displayGeneralAccount(generalData) {
        const generalHtml = `
            <div class="card mb-4 border-primary">
                <div class="card-body">
                    <h5 class="card-title text-primary mb-3">
                        <i class="fas fa-wallet"></i> Общий баланс
                    </h5>
                    <div class="d-flex flex-wrap gap-3">
                        ${Object.entries(generalData).map(([currency, amount]) => `
                            <div class="bg-light p-2 rounded">
                                <span class="font-weight-bold">${currency}</span>: ${amount}
                            </div>
                        `).join('')}
                    </div>
                </div>
            </div>
        `;
        accountingContent.innerHTML = generalHtml; // Перезаписываем контент
    }
    
    function displayAccounts(accounts) {
        const filteredAccounts = accounts.filter(account => 
            account.name.trim() !== "Внешний счёт"
        );
    
        if (filteredAccounts.length === 0) {
            accountingContent.innerHTML += '<p class="mt-3">Нет счетов для отображения.</p>';
            return;
        }
    
        filteredAccounts.forEach(account => {
            fetch(`http://localhost:8000/api/values/${account.value_id}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
            .then(response => {
                if (!response.ok) throw new Error('Ошибка загрузки валюты');
                return response.json();
            })
            .then(currency => {
                const accountRow = document.createElement('div');
                accountRow.className = 'account-row card mb-3';
                accountRow.dataset.accountId = account.account_id;
                accountRow.innerHTML = `
                    <div class="card-body d-flex justify-content-between align-items-center">
                        <div class="d-flex justify-content-between w-100 align-items-center">
                            <!-- Название и сумма -->
                            <div class="d-flex align-items-center gap-4">
                                <h5 class="card-title mb-0" style="margin-right: 10px;"> <!-- Инлайновый стиль -->
                                    <a href="/transactions?accounting_id=${account.accounting_id}&account_id=${account.account_id}" 
                                    class="text-decoration-none text-dark">
                                    ${account.name}
                                    </a>
                                </h5>
                                <div class="text-success fs-5">
                                    ${account.money_amount} ${currency.name}
                                </div>
                            </div>

                            <!-- Кнопки -->
                            <div>
                                <button class="btn btn-warning btn-sm me-2">Изменить</button>
                                <button class="btn btn-danger btn-sm">Удалить</button>
                            </div>
                        </div>
                    </div>
                `;
                accountingContent.appendChild(accountRow);
            })
            .catch(error => {
                console.error("Ошибка:", error);
            });
        });
    }

    function openEditModal(id, name, amount) {
        currentEditingId = id; // Сохраняем ID редактируемого счета
        document.getElementById('update-account-name').value = name; // Заполняем поле "Название счета"
        document.getElementById('update-account-balance').value = amount; // Заполняем поле "Состояние счета"
        editAccountModal.show(); // Открываем модальное окно
    }

    function updateAccount(id) {
        const name = document.getElementById('update-account-name').value;
        const amount = parseFloat(document.getElementById('update-account-balance').value);

        fetch(`http://localhost:8000/api/accounts/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name, money_amount: amount })
        })
        .then(response => {
            if (response.ok) {
                const accountingId = new URLSearchParams(window.location.search).get('id');
                loadAccounts(accountingId); // Обновляем список счетов
                editAccountModal.hide();
            } else {
                console.error("Failed to update account:", response.statusText);
            }
        })
        .catch(error => console.error("Error updating account:", error));
    }

    function deleteAccount(id) {
        fetch(`http://localhost:8000/api/accounts/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => {
            if (response.ok) {
                const accountingId = new URLSearchParams(window.location.search).get('id');
                loadAccounts(accountingId); // Обновляем список счетов
            } else {
                console.error("Failed to delete account:", response.statusText);
            }
        })
        .catch(error => console.error("Error deleting account:", error));
    }

    function createAccount(accountingId) {
        const name = document.getElementById('account-name').value;
        const amount = document.getElementById('account-balance').value;
        const currency = document.getElementById('account-currency').value;

        fetch(`http://localhost:8000/api/accountings/${accountingId}/accounts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                name: name,
                value: currency,
                money_amount: amount
            })
        })
        .then(response => {
            if (response.ok) {
                loadAccounts(accountingId); // Обновляем список счетов
                createAccountModal.hide();
                createAccountForm.reset(); // Сбрасываем форму
            } else {
                console.error("Failed to create account:", response.statusText);
            }
        })
        .catch(error => console.error("Error creating account:", error));
    }

    function loadAccountsForTransaction() {
        const accountingId = new URLSearchParams(window.location.search).get('id');
        fetch(`http://localhost:8000/api/accountings/${accountingId}/accounts`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load accounts');
            }
            return response.json();
        })
        .then(accounts => {
            const senderSelect = document.getElementById('sender-account');
            const receiverSelect = document.getElementById('receiver-account');

            accounts.forEach(account => {
                // Создаем опции для выпадающего списка
                const option = document.createElement('option');
                option.value = account.account_id;
                option.textContent = account.name;
                senderSelect.appendChild(option.cloneNode(true)); // Клон для отправителя
                receiverSelect.appendChild(option); // Для получателя
            });
        })
        .catch(error => console.error("Error loading accounts:", error));
    }

    loadAccountsForTransaction(); // Вызываем функцию после загрузки страницы

    // Обработчик события для создания транзакции
    document.getElementById('create-transaction-form').addEventListener('submit', function(event) {
        event.preventDefault();
        createTransaction();
    });

    function createTransaction() {
        const senderId = parseInt(document.getElementById('sender-account').value); // ID счета отправителя
        const receiverId = parseInt(document.getElementById('receiver-account').value); // ID счета получателя
        const amount = parseFloat(document.getElementById('transaction-amount').value); // Сумма
    
        const accountingId = new URLSearchParams(window.location.search).get('id'); // ID учета
    
        const transactionData = {
            account_id: senderId,          // Счет отправителя
            external_account_id: receiverId, // Счет получателя
            money_amount: amount           // Сумма
        };
    
        fetch(`http://localhost:8000/api/accountings/${accountingId}/accounts/transactions`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(transactionData)
        })
        .then(response => {
            if (response.ok) {
                alert("Транзакция успешно создана!");
                location.reload(); // Обновляем страницу для обновления данных
            } else {
                console.error("Failed to create transaction:", response.statusText);
                alert("Ошибка при создании транзакции. Пожалуйста, проверьте данные.");
            }
        })
        .catch(error => console.error("Error creating transaction:", error));
    }
});