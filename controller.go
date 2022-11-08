package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chsir-zy/anan/framework"
)

func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(2*time.Second))
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		time.Sleep(2 * time.Second)
		c.Json(200, "ok")

		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeOut()
	}

	return nil
}

func UserControllerHandler(c *framework.Context) error {
	fmt.Println("UserControllerHandler")
	time.Sleep(5 * time.Second)
	c.Json(200, "UserControllerHandler")
	return nil
}

func SubjectSubGetControllerHandler(c *framework.Context) error {
	fmt.Println("SubjectSubGetControllerHandler")
	time.Sleep(3 * time.Second)
	c.Json(200, "SubjectSubGetControllerHandler")
	return nil
}

func SubjectGetControllerHandler(c *framework.Context) error {
	c.Json(200, "SubjectGetControllerHandler")
	return nil
}
func SubjectPutControllerHandler(c *framework.Context) error {
	c.Json(200, "SubjectPutControllerHandler")
	return nil
}
func SubjectPostControllerHandler(c *framework.Context) error {
	c.Json(200, "SubjectPostControllerHandler")
	c.Next()
	return nil
}

func SubjectSubInfoGetControllerHandler(c *framework.Context) error {
	c.Json(200, "SubjectSubInfoGetControllerHandler")
	return nil
}
func SubjectSubInfoSunGetControllerHandler(c *framework.Context) error {
	c.Json(200, "SubjectSubInfoSunGetControllerHandler")
	return nil
}
