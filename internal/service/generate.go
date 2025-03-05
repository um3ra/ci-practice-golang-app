package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/mockery --name AuthService --output ./mocks --outpkg mocks --case underscore --with-expecter=true
