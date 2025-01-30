package jwt

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/mockery --name JwtService --output ./mocks --outpkg mocks --case underscore --with-expecter=true
