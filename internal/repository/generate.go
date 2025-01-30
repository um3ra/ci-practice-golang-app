package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/mockery --name UserRepository --output ./mocks --outpkg mocks --case underscore --with-expecter=true
