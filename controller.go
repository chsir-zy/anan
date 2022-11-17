package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		c.SetOkStatus().Json("ok")

		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.SetStatus(http.StatusInternalServerError).Json("panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.SetStatus(http.StatusInternalServerError).Json("time out")
		c.SetHasTimeOut()
	}

	return nil
}

func UserControllerHandler(c *framework.Context) error {
	fmt.Println("UserControllerHandler")
	ff, _ := c.FormFile("file")
	fmt.Println(ff.Filename)
	c.Json("UserControllerHandler")
	return nil
}

func SubjectSubGetControllerHandler(c *framework.Context) error {
	fmt.Println("SubjectSubGetControllerHandler")
	time.Sleep(3 * time.Second)
	c.Json("SubjectSubGetControllerHandler")
	return nil
}

func SubjectGetControllerHandler(c *framework.Context) error {
	c.Json("SubjectGetControllerHandler")
	return nil
}
func SubjectPutControllerHandler(c *framework.Context) error {
	c.Json("SubjectPutControllerHandler")
	return nil
}
func SubjectPostControllerHandler(c *framework.Context) error {
	c.Json("SubjectPostControllerHandler")
	c.Next()
	return nil
}

func SubjectSubInfoGetControllerHandler(c *framework.Context) error {
	c.Json("SubjectSubInfoGetControllerHandler")
	return nil
}
func SubjectSubInfoSunGetControllerHandler(c *framework.Context) error {
	c.Json("SubjectSubInfoSunGetControllerHandler")
	return nil
}
