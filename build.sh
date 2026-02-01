# build the frontend
cd frontend
npm install --force
npm run build
cd ..
go mod tidy
# build backend for mac
GOOS=darwin GOARCH=amd64 go build -o bin/homegym_mac_amd64 -v -buildmode=exe ./homegym
# build backend for windows
GOOS=windows GOARCH=amd64 go build -o bin/homegym_win_amd64.exe -v -buildmode=exe ./homegym
