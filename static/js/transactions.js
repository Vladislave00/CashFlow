document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    const usernameElement = document.getElementById('username');
    const accountTitle = document.getElementById('account-title');
    const transactionsContent = document.getElementById('transactions-content');
    const backButton = document.getElementById('back-btn');
    const logoutButton = document.getElementById('logout-btn');

    if (!token) {
        window.location.href = '/login';
    } else {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            usernameElement.textContent = payload.name;
            
            const urlParams = new URLSearchParams(window.location.search);
            const accountId = urlParams.get('account_id');
            const accountingId = urlParams.get('accounting_id');
            
            if (accountId && accountingId) {
                loadAccountName(accountId);
                loadTransactions(accountingId, accountId);
            }
        } catch (error) {
            console.error("Failed to decode token:", error);
        }
    }

    backButton.addEventListener('click', function() {
        window.history.back();
    });

    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('token');
        window.location.href = '/login';
    });

    async function loadAccountName(accountId) {
        try {
            const response = await fetch(`http://localhost:8000/api/accounts/${accountId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                const account = await response.json();
                accountTitle.textContent = `Транзакции счета: ${account.name}`;
            }
        } catch (error) {
            console.error("Error loading account name:", error);
        }
    }

    async function loadTransactions(accountingId, accountId) {
        try {
            const response = await fetch(`http://localhost:8000/api/accountings/${accountingId}/accounts/transactions/${accountId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (!response.ok) {
                throw new Error('Failed to load transactions');
            }
            
            const transactions = await response.json();
            displayTransactions(transactions);
        } catch (error) {
            console.error("Error loading transactions:", error);
        }
    }

    async function displayTransactions(transactions) {
        transactionsContent.innerHTML = '';
        
        if (transactions.length === 0) {
            transactionsContent.innerHTML = '<p>Нет транзакций для отображения.</p>';
            return;
        }

        for (const transaction of transactions) {
            const senderName = await getAccountName(transaction.account_id);
            const receiverName = await getAccountName(transaction.external_account_id);
            
            const transactionElement = document.createElement('div');
            transactionElement.className = 'transaction-item border p-3 mb-3';
            transactionElement.innerHTML = `
                <div class="d-flex justify-content-between">
                    <div>
                        <strong>Отправитель:</strong> ${senderName}<br>
                        <strong>Получатель:</strong> ${receiverName}
                    </div>
                    <div class="text-right">
                        <strong>Сумма:</strong> ${transaction.money_amount}<br>
                        <small>${new Date(transaction.created_at).toLocaleString()}</small>
                    </div>
                </div>
            `;
            transactionsContent.appendChild(transactionElement);
        }
    }

    async function getAccountName(accountId) {
        try {
            const response = await fetch(`http://localhost:8000/api/accounts/${accountId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                const account = await response.json();
                return account.name;
            }
            return 'Неизвестный счет';
        } catch (error) {
            console.error("Error fetching account name:", error);
            return 'Неизвестный счет';
        }
    }
});