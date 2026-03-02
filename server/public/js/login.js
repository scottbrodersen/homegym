const userNameInput = document.querySelector('input#username');
const passwordInput = document.querySelector('input#password');
const nameInput = document.querySelector('input#name');
const emailInput = document.querySelector('input#email');
const togglePasswordButton = document.querySelector('button#toggle-password');
const submitButton = document.getElementById('signin');

submitButton.disabled = true;

togglePasswordButton.addEventListener('click', togglePassword);

function togglePassword() {
  if (passwordInput.type === 'password') {
    passwordInput.type = 'text';
    togglePasswordButton.textContent = 'Hide password';
    togglePasswordButton.setAttribute('aria-label', 'Hide password.');
  } else {
    passwordInput.type = 'password';
    togglePasswordButton.textContent = 'Show password';
    togglePasswordButton.setAttribute(
      'aria-label',
      'Show password as plain text. ' +
        'Warning: this will display your password on the screen.',
    );
  }
}

passwordInput.addEventListener('input', validatePassword);
passwordInput.addEventListener('keydown', function (event) {
  if (event.key === 'Enter') {
    event.preventDefault();
    submitButton.click();
  }
});

submitButton.addEventListener('click', login);

function validatePassword() {
  let message = '';
  if (!/.{10,}/.test(passwordInput.value)) {
    message = 'At least ten characters. ';
    submitButton.disabled = true;
  } else {
    submitButton.disabled = false;
  }
  passwordInput.setCustomValidity(message);
}

async function login() {
  const url = '/homegym/login';
  const body = `{"username": "${userNameInput.value}", "password": "${passwordInput.value}", "form": "true"}`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: body,
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    const unsuccessfulMessage = document.getElementById('badlogin');
    unsuccessfulMessage.style.visibility = 'visible';
  } else if (resp.status < 200 || resp.status >= 300) {
    throw new Error('failed to login');
  }

  if (resp.redirected === true) {
    window.location.href = resp.url;
  }
}
