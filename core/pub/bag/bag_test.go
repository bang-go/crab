package bag

import (
	"errors"
	"testing"
)

func TestBaggerFinishWithError(t *testing.T) {
	b := NewBagger()

	executed := []int{}

	// 注册三个函数，第二个会报错
	b.Register(
		func() error {
			executed = append(executed, 1)
			return nil
		},
		func() error {
			executed = append(executed, 2)
			return errors.New("test error")
		},
		func() error {
			executed = append(executed, 3)
			return nil
		},
	)

	err := b.Finish()

	// 应该返回错误
	if err == nil {
		t.Error("expected error, got nil")
	}

	// 应该只执行了前两个函数，第三个不执行
	if len(executed) != 2 {
		t.Errorf("expected 2 functions executed, got %d", len(executed))
	}

	// 验证执行顺序
	if len(executed) >= 2 && (executed[0] != 1 || executed[1] != 2) {
		t.Errorf("unexpected execution order: %v", executed)
	}
}

func TestBaggerFinishOnlyOnce(t *testing.T) {
	b := NewBagger()

	counter := 0
	b.Register(func() error {
		counter++
		return nil
	})

	// 第一次调用
	err := b.Finish()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// 第二次调用不应该再执行
	err = b.Finish()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// counter 应该只增加一次
	if counter != 1 {
		t.Errorf("expected counter=1, got %d", counter)
	}
}

func TestBaggerConcurrency(t *testing.T) {
	b := NewBagger()

	// 并发注册
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			b.Register(func() error {
				return nil
			})
			done <- true
		}(i)
	}

	// 等待所有注册完成
	for i := 0; i < 10; i++ {
		<-done
	}

	err := b.Finish()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
