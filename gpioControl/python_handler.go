/*
    直接调用py脚本操作gpio方式
*/

package gpioControl

import (
	"fmt"
	"os/exec"
	"regexp"
	// "encoding/json"
)

func ExecPython(params []string) (string, error) {
	var args = append([]string{}, params...)
	out, err := exec.Command("python3", args...).Output()
	if err != nil {
		fmt.Println("exec python err =>", err.Error())
		return "", err
	}

	// out_emp_ind := bytes.IndexByte(out, 0)

	// if out_emp_ind < 0 {
	// 	out_emp_ind = out_emp_ind
	// }
	// out = out[:out_emp_ind]

	outData := string(out)

	reg := regexp.MustCompile("\\s+")
	outData = reg.ReplaceAllString(outData, "")

	return outData, nil
}

func PyGpioExec(params ...string) (string, error) {
    tmpParams := []string{"gpio.py"}
    tmpParams = append(tmpParams, params...)
    return ExecPython(tmpParams)
}
