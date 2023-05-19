#!/bin/sh

output_dir="docs/generated/coverage/go"
mkdir -p ${output_dir}
rm -f ${output_dir}/*.cov
PKG_LIST=$(go list ./internal/... | grep -v /vendor/ | grep -v /docs/ )
echo "mode: count" > ${output_dir}/coverage.cov
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${output_dir}/${package##*/}.cov" "$package" ;
	tail -q -n +2 "${output_dir}/${package##*/}.cov" >> ${output_dir}/coverage.cov
	rm "${output_dir}/${package##*/}.cov"
done
go tool cover -func=${output_dir}/coverage.cov
go tool cover -html=${output_dir}/coverage.cov -o ${output_dir}/coverage.html
