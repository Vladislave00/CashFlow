document.getElementById('loginForm').addEventListener('submit', async function(event) {
    event.preventDefault(); // Отменяем стандартное поведение формы

    const formData = {
        email: document.getElementById('email').value,
        password: document.getElementById('password').value
    };

    try {
        const response = await fetch('/auth/sign-in', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            const result = await response.json();
            localStorage.setItem('token', result.token); // Сохраняем токен в localStorage
            window.location.href = '/'; // Перенаправляем на главную страницу
        } else {
            const errorText = await response.text();
            alert(`Error: ${errorText}`); // Выводим ошибку
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred. Please try again.');
    }
});
