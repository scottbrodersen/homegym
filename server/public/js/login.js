const userNameInput = document.querySelector("input#username");
const passwordInput = document.querySelector("input#password");
const nameInput = document.querySelector("input#name");
const emailInput = document.querySelector("input#email");

const togglePasswordButton = document.querySelector("button#toggle-password");
const badCreds = document.querySelector("div#bad");

const submitButton = document.querySelector("button");
togglePasswordButton.addEventListener("click", togglePassword);

function togglePassword() {
  if (passwordInput.type === "password") {
    passwordInput.type = "text";
    togglePasswordButton.textContent = "Hide password";
    togglePasswordButton.setAttribute("aria-label", "Hide password.");
  } else {
    passwordInput.type = "password";
    togglePasswordButton.textContent = "Show password";
    togglePasswordButton.setAttribute(
      "aria-label",
      "Show password as plain text. " +
        "Warning: this will display your password on the screen."
    );
  }
}

passwordInput.addEventListener("input", validatePassword);

function validatePassword() {
  let message = "";
  if (!/.{10,}/.test(passwordInput.value)) {
    message = "At least ten characters. ";
    submitButton.disabled = true;
  } else {
    submitButton.disabled = false;
  }
  passwordInput.setCustomValidity(message);
}
