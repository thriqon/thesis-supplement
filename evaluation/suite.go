package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type point struct {
	diffBytes int
	diffTime  time.Duration
}

type result struct {
	ttc   point
	ttcRe point
	size  int
}

type suite struct {
	group   string
	name    string
	f       func() result
	enabled bool
}

var _ flag.Value = (*suite)(nil)

func (s *suite) String() string {
	return "false"
}

func (s *suite) Set(_ string) error {
	(*s).enabled = true
	return nil
}

func (s *suite) IsBoolFlag() bool {
	return true
}

var suites []*suite

func registerSuite(group, name string, f func() result) {
	slug := slugify(group + "-" + name)
	sx := suite{group, name, f, false}
	suites = append(suites, &sx)
	flag.Var(&sx, slug, "Enable ["+group+"/"+name+"]")
}

func formatDurationInMs(d time.Duration) string {
	return fmt.Sprintf("%v", commify(int(d.Nanoseconds()/1e6)))
}

func formatBytes(v int) string {
	v /= 1024
	return commify(v)
}

// copied from https://github.com/dustin/go-humanize/blob/master/comma.go#L14
func commify(v int) string {
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(int64(v)%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return strings.Join(parts[j:], ",")
	//*/
}

const N = 1

func average(results [N]result) result {
	ttc := point{}
	ttcRe := point{}
	size := 0

	for _, x := range results {
		ttc.diffBytes += x.ttc.diffBytes
		ttc.diffTime += x.ttc.diffTime
		ttcRe.diffBytes += x.ttcRe.diffBytes
		ttcRe.diffTime += x.ttcRe.diffTime
		size += x.size
	}

	ttc.diffBytes /= N
	ttc.diffTime /= N
	ttcRe.diffBytes /= N
	ttcRe.diffTime /= N
	size /= N

	return result{ttc, ttcRe, size}
}

func (s *suite) run() result {
	var res [N]result
	for i := 0; i < N; i++ {
		rmAllContainers()
		rmDanglingImages()
		res[i] = s.f()
	}
	return average(res)
}

type namedres struct {
	name string
	r    result
}

func main() {
	flag.Parse()
	results := make(map[string][]namedres)

	for _, s := range suites {
		if s.enabled {
			results[s.group] = append(results[s.group], namedres{s.name, s.run()})
		}
	}

	fmt.Println(`
	\begin{table}[tbp]
	\begin{tabular}{@{\extracolsep{4pt}}l l r r r r r@{}}\hline
	\multicolumn{2}{l}{Name} & \multicolumn{2}{c}{Initial Compile} & \multicolumn{2}{c}{Recompile} & Size \\
	\cline{3-4}\cline{5-6}\cline{7-7}
	& & KiB & ms & KiB & ms & KiB\\\hline`)
	format := " & %-20s&%14s&%10v&%14v&%10v&%14v\\\\\n"

	for key := range results {
		val := results[key]
		fmt.Printf("\\hline\\multicolumn{7}{l}{%s}\\\\\n", key)
		for _, nres := range val {
			res := nres.r
			fmt.Printf(format, nres.name, formatBytes(res.ttc.diffBytes), formatDurationInMs(res.ttc.diffTime), formatBytes(res.ttcRe.diffBytes), formatDurationInMs(res.ttcRe.diffTime), formatBytes(res.size))
		}
	}

	fmt.Println(`\end{tabular}`)
	fmt.Println(`\caption{Measurements}`)
	fmt.Println(`\label{tab:measurements}`)
	fmt.Println(`\end{table}`)
	fmt.Println()
	fmt.Println()
	fmt.Println()

	for key := range results {
		val := results[key]
		coords := make([]string, len(val))
		for i, v := range val {
			coords[i] = fmt.Sprintf("%v\\\\(Size: %v KiB)", v.name, v.r.size/1024)
		}
		fmt.Printf("\\def\\diag%s{%%\n", slugify(key))
		fmt.Printf(`\pgfplotsset{enlargelimits=0.15,width=8cm,xticklabels={%v},xtick=data,ymin=0,tick label style={/pgf/number format/fixed,align=center}}`, strings.Join(coords, ","))
		fmt.Println(`\begin{figure}[htbp]`)
		fmt.Println(`\begin{tikzpicture}`)
		fmt.Println(`\begin{axis}[ybar,axis y line*=left,ylabel={Transmitted data in KiB},x tick label style={rotate=45,anchor=east},scaled y ticks = false]`)
		fmt.Print(`\addplot coordinates {`)
		for i, nres := range val {
			fmt.Printf("(%v, %v) ", i, nres.r.ttc.diffBytes/1024)
		}
		fmt.Printf("};\\label{plot:%s}\n", slugify(key+"-"+"ttc-kib"))
		fmt.Print("\\addplot coordinates {")
		for i, nres := range val {
			fmt.Printf("(%v, %v) ", i, nres.r.ttcRe.diffBytes/1024)
		}
		fmt.Printf("};\\label{plot:%s}\n", slugify(key+"-"+"ttrc-kib"))
		fmt.Println(`\end{axis}`)
		fmt.Println(`\begin{axis}[axis y line*=right,ylabel={Time in ms},xticklabels={},scaled y ticks = false]`)

		fmt.Print("\\addplot coordinates {")
		for i, nres := range val {
			fmt.Printf("(%v, %v) ", i, nres.r.ttc.diffTime.Nanoseconds()/1e6)
		}
		fmt.Printf("};\\label{plot:%s}\n", slugify(key+"-"+"ttc-ms"))
		fmt.Print("\\addplot coordinates {")
		for i, nres := range val {
			fmt.Printf("(%v, %v) ", i, nres.r.ttcRe.diffTime.Nanoseconds()/1e6)
		}
		fmt.Printf("};\\label{plot:%s}\n", slugify(key+"-"+"ttrc-ms"))
		fmt.Println(`\end{axis}`)
		fmt.Println(`\end{tikzpicture}`)
		fmt.Printf("\\caption{Measurement results for '%s' example}\n", key)
		fmt.Printf("\\label{fig:mes-res-%s}\n", slugify(key))
		fmt.Println(`\end{figure}`)
		fmt.Println(`}`)
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}
}

// utils

func writeFile(fn, contents string) {
	if err := ioutil.WriteFile(fn, []byte(contents), 0755); err != nil {
		panic(err)
	}
}

func readFileInt(f string) int {
	contents, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		panic(err)
	}
	return n
}

func rtxed(ifaces []string) int {
	val := 0
	for _, iface := range ifaces {
		rx := readFileInt("/sys/class/net/" + iface + "/statistics/rx_bytes")
		tx := readFileInt("/sys/class/net/" + iface + "/statistics/tx_bytes")
		val += rx + tx
	}
	return val
}

func rmDanglingImages() {
	_ = exec.Command("/bin/sh", "-c", "docker rmi $(docker images -q)").Run()
}

func rmAllContainers() {
	_ = exec.Command("/bin/sh", "-c", "docker rm $(docker ps -aq)").Run()
}

func measure(f interface{}, cleanup interface{}) point {
	ifaces := []string{"docker0", "wlan0"}
	beforeBytes := rtxed(ifaces)
	beforeTime := time.Now()

	switch ff := f.(type) {
	case func() error:
		if err := ff(); err != nil {
			panic(err)
		}
	case func():
		ff()
	default:
		panic("Unknown measured function type")
	}

	afterTime := time.Now()
	afterBytes := rtxed(ifaces)

	if cleanup != nil {
		if cf, ok := cleanup.(func()); ok {
			cf()
		}
	}

	return point{afterBytes - beforeBytes, afterTime.Sub(beforeTime)}
}

func ensureNotPresent(images ...string) {
	for _, img := range images {
		_ = exec.Command("docker", "rmi", "-f", img).Run()
	}
}

func imageSize(image string) int {
	output, err := exec.Command("docker", "inspect", "-f", "{{.VirtualSize}}", image).Output()
	if err != nil {
		panic(err)
	}

	num, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		panic(err)
	}
	return num
}

func mustOutputString(args ...string) string {
	cmd := exec.Command(args[0], args[1:]...)
	//cmd.Stderr = os.Stderr
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(o))
}

func commandRunner(args ...string) func() error {
	return func() error {
		return exec.Command(args[0], args[1:]...).Run()
	}
}

func mustRun(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("mustRun", args[0], "failed:", err)
		panic(err)
	}
}

func slugify(s string) string {
	mapper := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'a' + (r - 'A')
		case r >= 'a' && r <= 'z':
			return r
		default:
			return '-'
		}
	}

	return strings.Map(mapper, s)
}
