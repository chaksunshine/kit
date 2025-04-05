package thread

import (
	"context"
	"sync"
)

// 线程池
// @author fuzeyu
// @date 2025/4/4
type Pool struct {

	// 上下文
	ctx        context.Context
	cancelFunc context.CancelFunc

	// 协程工作数
	size int
	wait *sync.WaitGroup

	// 错误信息
	err error
	rw  *sync.Mutex
}

// 记录错误信息
// @param err 错误信息
func (obj *Pool) recordErr(err error) {

	// 执行取消
	obj.cancelFunc()

	// 保存错误
	obj.rw.Lock()
	if err != nil && obj.err == nil {
		obj.err = err
	}
	obj.rw.Unlock()
}

// 创建并执行工作
// @param call 回调方法
func (obj *Pool) Working(call func(ctx context.Context, index int) error) {

	obj.wait.Add(obj.size)

	// 执行
	for index := 0; index < obj.size; index++ {
		go func(index int) {
			defer obj.wait.Done()

			// 执行,如果存在错误记录错误
			if err := call(obj.ctx, index); err != nil {
				obj.recordErr(err)
				return
			}

			// 检查上下文是否取消
			if IsCancel(obj.ctx) {
				obj.recordErr(obj.ctx.Err())
				return
			}
		}(index)
	}
}

// 获取结果
func (obj *Pool) Result() error {
	obj.wait.Wait()

	obj.rw.Lock()
	var err = obj.err
	obj.rw.Unlock()
	return err
}

// @param ctx 上下文
// @param size 执行的协程池大小
func NewPool(ctx context.Context, size int) *Pool {

	cancel, cancelFunc := context.WithCancel(ctx)

	c := &Pool{
		ctx:        cancel,
		cancelFunc: cancelFunc,
		size:       size,
		wait:       new(sync.WaitGroup),
		rw:         new(sync.Mutex),
	}
	return c
}
