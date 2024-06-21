document.addEventListener('DOMContentLoaded', function () {
    const loginContainer = document.getElementById('loginContainer');
    const registerContainer = document.getElementById('registerContainer');
    const showRegisterLink = document.getElementById('showRegister');
    const showLoginLink = document.getElementById('showLogin');

    // Event listener to show the register form and hide the login form
    showRegisterLink.addEventListener('click', function (event) {
        event.preventDefault();
        registerContainer.classList.remove('hidden');
        loginContainer.classList.add('hidden');
    });

    // Event listener to show the login form and hide the register form
    showLoginLink.addEventListener('click', function (event) {
        event.preventDefault();
        loginContainer.classList.remove('hidden');
        registerContainer.classList.add('hidden');
    });
});
