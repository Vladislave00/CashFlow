document.getElementById('registerForm').addEventListener('submit', async function(event) {
    event.preventDefault(); 

    const formData = {
        name: document.getElementById('nickname').value,
        email: document.getElementById('email').value,
        password: document.getElementById('password').value
    };

    try {
        const response = await fetch('/auth/sign-up', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            const result = await response.json();
            alert(`User created with ID: ${result.user_id}`); 
            window.location.href = '/login'; 
        } else {
            const errorText = await response.text();
            alert(`Error: ${errorText}`); 
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred. Please try again.');
    }
});
