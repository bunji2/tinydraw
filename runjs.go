package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

var d *Data

// RunJS は JavaScript のコードを実行する関数
func RunJS(jsFilePath string) (err error) {
	/*
	   	start := time.Now()
	       defer func() {
	           duration := time.Since(start)
	           if caught := recover(); caught != nil {
	               if caught == halt {
	                   fmt.Fprintf(os.Stderr, "Some code took to long! Stopping after: %v\n", duration)
	                   return
	               }
	               panic(caught) // Something else happened, repanic!
	           }
	           fmt.Fprintf(os.Stderr, "Ran code successfully: %v\n", duration)
	   	}()
	*/
	vm := otto.New()

	/*
	   	var script *otto.Scriptvm.Interrupt = make(chan func(), 1) // The buffer prevents blocking
	       go func() {
	           time.Sleep(10 * time.Second) // Stop after ten seconds
	           vm.Interrupt <- func() {
	               panic(halt)
	           }
	   	}()
	*/
	var script *otto.Script
	script, err = vm.Compile(jsFilePath, nil)
	if err != nil {
		return
	}
	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})

	vm.Set("setParams", func(call otto.FunctionCall) otto.Value {
		var gw, gh, margin int64
		if len(call.ArgumentList) > 1 {
			gw, _ = call.Argument(0).ToInteger()
			gh, _ = call.Argument(1).ToInteger()
		} else if len(call.ArgumentList) > 2 {
			margin, _ = call.Argument(2).ToInteger()
		}
		SetParams(int(gw), int(gh), int(margin))
		return otto.Value{}
	})

	vm.Set("setGrid", func(call otto.FunctionCall) otto.Value {
		w, _ := call.Argument(0).ToInteger()
		h, _ := call.Argument(1).ToInteger()
		d = NewDraw(int(w), int(h))
		return otto.Value{}
	})

	/*
		vm.Set("setGridVisible", func(call otto.FunctionCall) otto.Value {
			if d != nil {
				b, _ := call.Argument(0).ToBoolean()
				d.setGridVisible(b)
			}
			return otto.Value{}
		})
	*/

	vm.Set("fillGridSquare", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			x1, _ := call.Argument(0).ToInteger()
			y1, _ := call.Argument(1).ToInteger()
			x2, _ := call.Argument(2).ToInteger()
			y2, _ := call.Argument(3).ToInteger()
			d.fillGridSquare(int(x1), int(y1), int(x2), int(y2))
		}
		return otto.Value{}
	})

	vm.Set("fillAll", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			cs, _ := call.Argument(0).ToString()
			c := getColorOfStr(cs)
			if c != nil {
				d.fillAll(c)
			}
		}
		return otto.Value{}
	})

	vm.Set("fillSquare", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			x1, _ := call.Argument(0).ToInteger()
			y1, _ := call.Argument(1).ToInteger()
			x2, _ := call.Argument(2).ToInteger()
			y2, _ := call.Argument(3).ToInteger()
			d.fillSquare(int(x1), int(y1), int(x2), int(y2))
		}
		return otto.Value{}
	})

	vm.Set("setFgColor", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			if len(call.ArgumentList) == 1 {
				cs, _ := call.Argument(0).ToString()
				c := getColorOfStr(cs)
				if c != nil {
					d.setFgColor(c)
				}
			} else {
				r, _ := call.Argument(0).ToInteger()
				g, _ := call.Argument(1).ToInteger()
				b, _ := call.Argument(2).ToInteger()
				d.setFgColor(NewColor(uint8(r), uint8(g), uint8(b)))
			}
		}
		return otto.Value{}
	})

	vm.Set("setBgColor", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			if len(call.ArgumentList) == 1 {
				cs, _ := call.Argument(0).ToString()
				c := getColorOfStr(cs)
				if c != nil {
					d.setBgColor(c)
				}
			} else {
				r, _ := call.Argument(0).ToInteger()
				g, _ := call.Argument(1).ToInteger()
				b, _ := call.Argument(2).ToInteger()
				d.setBgColor(NewColor(uint8(r), uint8(g), uint8(b)))
			}
		}
		return otto.Value{}
	})

	vm.Set("drawGridLine", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			x1, _ := call.Argument(0).ToInteger()
			y1, _ := call.Argument(1).ToInteger()
			x2, _ := call.Argument(2).ToInteger()
			y2, _ := call.Argument(3).ToInteger()
			d.drawGridLine(int(x1), int(y1), int(x2), int(y2))
		}
		return otto.Value{}
	})

	vm.Set("drawLine", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			x1, _ := call.Argument(0).ToInteger()
			y1, _ := call.Argument(1).ToInteger()
			x2, _ := call.Argument(2).ToInteger()
			y2, _ := call.Argument(3).ToInteger()
			d.drawLine(int(x1), int(y1), int(x2), int(y2))
		}
		return otto.Value{}
	})

	vm.Set("saveFile", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			b, _ := call.Argument(0).ToString()
			d.saveFile(b)
		}
		return otto.Value{}
	})

	vm.Set("drawText", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			gx, _ := call.Argument(0).ToInteger()
			gy, _ := call.Argument(1).ToInteger()
			text, _ := call.Argument(2).ToString()
			d.drawText(int(gx), int(gy), text)
		}
		return otto.Value{}
	})

	vm.Set("drawGridSquares", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			d.drawGridSquares()
		}
		return otto.Value{}
	})

	vm.Set("drawGridText", func(call otto.FunctionCall) otto.Value {
		if d != nil {
			gx, _ := call.Argument(0).ToInteger()
			gy, _ := call.Argument(1).ToInteger()
			text, _ := call.Argument(2).ToString()
			d.drawGridText(int(gx), int(gy), text)
		}
		return otto.Value{}
	})

	_, err = vm.Run(script)

	return
}
