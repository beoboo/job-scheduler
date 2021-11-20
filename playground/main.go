package playground

//import (
//	"bufio"
//	"fmt"
//	"os/exec"
//	"strings"
//)
//
//func main() {
//	args := "5 1"
//	cmd := exec.Command("./test.sh", strings.Split(args, " ")...)
//
//	stdout, _ := cmd.StdoutPipe()
//	cmd.Start()
//
//	scanner := bufio.NewScanner(stdout)
//	scanner.Split(bufio.ScanWords)
//	for scanner.Scan() {
//		m := scanner.Text()
//		fmt.Println(m)
//	}
//	cmd.Wait()
//}
