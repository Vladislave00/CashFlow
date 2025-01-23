document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    const usernameElement = document.getElementById('username');
    const logoutButton = document.getElementById('logout-btn');
    const accountingsList = document.getElementById('accountings-list');
    const createAccountingBtn = document.getElementById('create-accounting-btn');
    const createAccountingModal = new bootstrap.Modal(document.getElementById('createAccountingModal'));
    const editAccountingModal = new bootstrap.Modal(document.getElementById('editAccountingModal'));
    const createAccountingForm = document.getElementById('createAccountingForm');
    const accountingNameInput = document.getElementById('accountingName');
    const editAccountingForm = document.getElementById('editAccountingForm');
    const editAccountingNameInput = document.getElementById('editAccountingName');

    let currentEditingId = null; // Для хранения ID редактируемого учета

    if (!token) {
        window.location.href = '/login';
    } else {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            usernameElement.textContent = payload.name; 
            loadAccountings(); 
        } catch (error) {
            console.error("Failed to decode token:", error);
        }
    }

    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('token');
        window.location.href = '/login';
    });

    createAccountingBtn.addEventListener('click', function() {
        createAccountingModal.show(); // Показываем модальное окно
    });

    createAccountingForm.addEventListener('submit', function(event) {
        event.preventDefault();
        createAccounting();
    });

    editAccountingForm.addEventListener('submit', function(event) {
        event.preventDefault();
        updateAccounting();
    });

    function loadAccountings() {
        fetch('http://localhost:8000/api/accountings/', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => response.json())
        .then(data => {
            accountingsList.innerHTML = ''; // Очищаем список перед добавлением
            data.data.forEach(accounting => {
                const li = document.createElement('li');
                li.className = 'list-group-item d-flex justify-content-between align-items-center';
                
                // Создаем ссылку на страницу учета
                const accountingLink = document.createElement('a');
                accountingLink.href = `/accounting?id=${accounting.id}`; // Изменили путь на /accounting
                accountingLink.textContent = accounting.name; // Имя учета
                accountingLink.className = 'accounting-link';
                
                li.appendChild(accountingLink);
    
                const buttonGroup = document.createElement('div');
                buttonGroup.className = 'btn-group';
    
                const editButton = document.createElement('button');
                editButton.className = 'btn btn-warning btn-sm';
                editButton.textContent = 'Изменить';
                editButton.onclick = () => openEditModal(accounting.id, accounting.name);
                
                const deleteButton = document.createElement('button');
                deleteButton.className = 'btn btn-danger btn-sm';
                deleteButton.textContent = 'Удалить';
                deleteButton.onclick = () => deleteAccounting(accounting.id);
    
                buttonGroup.appendChild(editButton);
                buttonGroup.appendChild(deleteButton);
                li.appendChild(buttonGroup);
                accountingsList.appendChild(li);
            });
        })
        .catch(error => console.error("Error loading accountings:", error));
    }
    
    

    function createAccounting() {
        const name = accountingNameInput.value;

        fetch('http://localhost:8000/api/accountings/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name }) // Отправляем только имя
        })
        .then(response => {
            if (response.ok) {
                loadAccountings(); // Обновляем список учетов
                createAccountingModal.hide(); // Закрываем модальное окно
                accountingNameInput.value = ''; // Очищаем поле ввода
            } else {
                console.error("Failed to create accounting:", response.statusText);
            }
        })
        .catch(error => console.error("Error creating accounting:", error));
    }

    function openEditModal(id, name) {
        currentEditingId = id; // Сохраняем ID редактируемого учета
        editAccountingNameInput.value = name; // Заполняем поле ввода текущим именем
        editAccountingModal.show(); // Показываем модальное окно редактирования
    }

    function updateAccounting() {
        const newName = editAccountingNameInput.value;

        fetch(`http://localhost:8000/api/accountings/${currentEditingId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name: newName })
        })
        .then(response => {
            if (response.ok) {
                loadAccountings(); // Обновляем список учетов
                editAccountingModal.hide(); // Закрываем модальное окно
            } else {
                console.error("Failed to update accounting:", response.statusText);
            }
        })
        .catch(error => console.error("Error updating accounting:", error));
    }

    function deleteAccounting(id) {
        if (confirm("Вы уверены, что хотите удалить этот учет?")) {
            fetch(`http://localhost:8000/api/accountings/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
            .then(response => {
                if (response.ok) {
                    loadAccountings(); // Обновляем список учетов
                } else {
                    console.error("Failed to delete accounting:", response.statusText);
                }
            })
            .catch(error => console.error("Error deleting accounting:", error));
        }
    }
});
