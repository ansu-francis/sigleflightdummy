package main

import(
	"fmt"
	"sync"
	"time"
)

type Call struct {
	val interface{}
	err error
	wg sync.WaitGroup
}

type Group struct {
	m map[string] *Call
	mu sync.Mutex
}

func(g *Group) Do(key string, fn func()(interface{}, error)) (interface{}, error) {
	
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string] *Call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := Call{}
	c.wg.Add(1)
	g.m[key] = &c
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()
	g.mu.Lock()
	if g.m[key] == &c {
		delete(g.m, key)
		g.mu.Unlock()
	}
		
	return c.val, c.err	
}

var g Group

func test(f func()(interface{}, error), done chan bool) {
	res, err := g.Do("test", f)
	fmt.Println("Inside test")
	fmt.Println(res.(string))
	fmt.Println(err)
	
}
	
func main() {
	done := make(chan bool)
	f := func()(interface{}, error) {
		fmt.Println("inside fn")
		time.Sleep(10*time.Second)
		fmt.Println("completed fn")
		return "done", nil
	}
	
	go func() {
		test(f, done)
		close(done)
	}()
	
	res, err := g.Do("test", f)
	fmt.Println("inside main")
	fmt.Println(res.(string))
	fmt.Println(err)
	<-done
	
	test(f, done)
	
	
}
