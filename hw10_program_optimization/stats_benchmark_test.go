package hw10programoptimization

import (
	"archive/zip"
	"github.com/stretchr/testify/require"
	"testing"
)

func Benchmark_GetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(b, err)

		require.Equal(b, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(b, err)
		_, err = GetDomainStat(data, "biz")
		require.NoError(b, err)
		r.Close()
	}
}
