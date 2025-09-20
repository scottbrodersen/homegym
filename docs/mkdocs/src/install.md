# Build and install

Your first task is to build and install the HomeGym server. Once complete, you'll be able to access the server over a secure connection and authenticate with a user name and password.

## Build

To build HomeGym, you need to get a copy of the GitHub repo onto your computer by cloning it. Follow the [GitHub docs](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) to clone the repo at [https://github.com/scottbrodersen/homegym](https://github.com/scottbrodersen/homegym).

Install the latest version of the Go programming language. The [Go download page](https://go.dev/doc/install) provides instructions.

Now you're ready to build:

1. Locate the build.sh file in the root directory of the homegym repository.
2. Double-click the file to run it.

Two executable files are generated and saved in the bin folder. Each file is the HomeGym program that you can run on your computer:

- Mac: homegym_mac_amd64
- Windows: homegym_win_amd64.exe

## Create a self-signed certificate

Use the following command to create a private key. Do not share the key with anyone!

`openssl genrsa -out homegym.key 2048`

Now, create a certificate signing request using the private key. After entering the following command, you'll be prompted for more information about your location and identity:

`openssl req -new -key homegym.key -out homegym.csr`

Now, create the self-signed certificate using the private key and certificate signing request:

`openssl x509 -req -days 365 -in homegym.csr -signkey homegym.key -out homegym.crt`

## Install and run

1. Create a folder in your home directory named HomeGym.
2. In the HomeGym folder, create a folder named database.
3. Create an environment variable named `HOMEGYM_DB_PATH` and the value is the path to the database folder that you just created.
4. Copy the HomeGym executable file to the HomeGym folder.
5. Copy the homegym.key and homegym.crt files to the HomeGym folder.
6. Double-click the executable file to run it. (When you want to stop it, press Ctrl+c.)

## Create a user account and sign in

1. In your web browser, go to [https://127.0.0.1:443/homegym/signup/](https://127.0.0.1:443/homegym/signup/).
2. Enter a user name, your email address, and a password and click Sign Up.
3. Go to [http://127.0.0.1:3000/homegym/login/](http://127.0.0.1:3000/homegym/login/) and sign in.
