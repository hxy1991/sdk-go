package rand

import (
	"github.com/hxy1991/sdk-go/utils"
	"math/rand"
	"sync"
)

var RandPool = NewRandomPool(1)

// RandomPool 结构表示随机数生成器池
type RandomPool struct {
	pool *sync.Pool
}

// NewRandomPool 创建一个随机数生成器池，并初始时创建指定数量的随机数生成器
func NewRandomPool(size int) *RandomPool {
	pool := &sync.Pool{
		New: func() interface{} {
			// 创建新的随机数生成器
			source := rand.NewSource(utils.Now().UnixNano())
			return rand.New(source)
		},
	}

	// 预先创建指定数量的随机数生成器并放入池中
	for i := 0; i < size; i++ {
		pool.Put(pool.New())
	}

	return &RandomPool{
		pool: pool,
	}
}

// GenerateRandom 从随机数生成器池中生成随机数
func (rp *RandomPool) GenerateRandom() int {
	r := rp.pool.Get().(*rand.Rand)
	defer rp.pool.Put(r)

	return r.Intn(100)
}
