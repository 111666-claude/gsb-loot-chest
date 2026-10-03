// Command loot-chest 跑宝箱样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/chest"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("loot-chest", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "keys", "keys / cap / refresh / idem / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "keys":
		world := chest.New(2, 100000)
		fmt.Fprintf(stdout, "ok=%v keys=%d\n", world.Open("p1", "c1", 0, "r1"), world.Keys("p1"))
	case "cap":
		world := chest.New(2, 100000)
		world.AddKeys("p1", 5)
		world.Open("p1", "c1", 0, "r1")
		world.Open("p1", "c1", 0, "r2")
		fmt.Fprintf(stdout, "third=%v opens=%d\n", world.Open("p1", "c1", 0, "r3"), world.Opens("c1", 0))
	case "refresh":
		world := chest.New(1, 1000)
		world.AddKeys("p1", 5)
		world.Open("p1", "c1", 0, "r1")
		before := world.Open("p1", "c1", 900, "r2")
		after := world.Open("p1", "c1", 1500, "r3")
		fmt.Fprintf(stdout, "before=%v after=%v\n", before, after)
	case "idem":
		world := chest.New(5, 100000)
		world.AddKeys("p1", 5)
		world.Open("p1", "c1", 0, "same")
		fmt.Fprintf(stdout, "second=%v keys=%d\n", world.Open("p1", "c1", 0, "same"), world.Keys("p1"))
	case "scan":
		world := chest.New(5, 100000)
		for index := 0; index < 3000; index++ {
			world.AddKeys("p1", 1)
			world.Open("p1", fmt.Sprintf("c-%d", index), 0, fmt.Sprintf("r-%d", index))
		}
		for index := 0; index < 3000; index++ {
			world.Opens("c-2999", 0)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", world.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample keys|cap|refresh|idem|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
