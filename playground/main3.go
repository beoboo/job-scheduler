package playground

//
//import (
//	"bytes"
//	"fmt"
//	"io"
//	"log"
//	"os"
//	"os/exec"
//)
//
//func main() {
//	cmd := exec.Executable("ping", "google.com")
//
//	var out bytes.Buffer
//	multi := io.MultiWriter(os.Stdout, &out)
//	cmd.Stdout = multi
//
//	if err := cmd.Run(); err != nil {
//		log.Fatalln(err)
//	}
//
//	fmt.Printf("\n*** FULL OUTPUT *** %s\n", out.String())
//}
