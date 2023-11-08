cd frontend
npm install
npm run build
cd ..
go mod tidy
GOOS=darwin GOARCH=amd64 go build -o bin/homegym_mac_amd64 -v -buildmode=exe ./homegym
GOOS=windows GOARCH=amd64 go build -o bin/homegym_win_amd64.exe -v -buildmode=exe ./homegym
