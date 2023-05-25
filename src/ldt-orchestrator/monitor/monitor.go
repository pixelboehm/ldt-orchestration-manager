package monitor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"longevity/src/communication"
	. "longevity/src/types"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Monitor struct {
	Started       chan *Process
	Stopped       chan int
	processes     []Process
	ldt_list_path string
}

func NewMonitor(ldt_list_path string) *Monitor {
	return &Monitor{
		Started:       make(chan *Process),
		Stopped:       make(chan int),
		ldt_list_path: ldt_list_path,
	}
}

func (m *Monitor) Run(port int) {
	fs := http.FileServer(http.Dir("static"))
	rest := communication.NewRestInterface(nil)
	rest.Router().Handle("/static/", http.StripPrefix("/static/", fs))
	rest.AddCustomHandler("/", m.handler)
	rest.Run(port)
}

func (m *Monitor) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	tmp := template.Must(
		template.New("index.html").Funcs(template.FuncMap{
			"formatJSON":  formatJSON,
			"convertTime": convertTime,
		}).ParseFiles("static/index.html"),
	)

	data := map[string]interface{}{
		"Processes": m.processes,
	}

	err := tmp.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func formatJSON(data json.RawMessage) string {
	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return string(data)
	}
	return string(formatted)
}

func convertTime(started string) string {
	currentTime := time.Now().Format("2006-1-2 15:4:5")
	newCurrentTime, err := time.Parse("2006-1-2 15:4:5", currentTime)
	if err != nil {
		log.Println("Monitor: Failed to parse time")
		return "Unknown"
	}
	startTime, err := time.Parse("2006-1-2 15:4:5", started)
	if err != nil {
		log.Println("Monitor: Failed to parse time")
		return "Unknown"
	}
	uptime := newCurrentTime.Sub(startTime)
	return fmt.Sprint(uptime)
}

func (m *Monitor) DoKeepAlive() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		log.Printf("Monitor: Currently Active LDTs %d\n", len(m.processes))
		for _, ldt := range m.processes {
			if !ldtIsRunning(ldt.Pid) {
				m.Stopped <- ldt.Pid
			}
		}
		<-ticker.C
	}
}

func (m *Monitor) RefreshLDTs() {
	for {
		select {
		case started := <-m.Started:
			m.RegisterLDT(started)
		case stopped := <-m.Stopped:
			m.RemoveLDT(stopped)
		default:
		}
	}
}

func (m *Monitor) RegisterLDT(ldt *Process) {
	m.processes = append(m.processes, *ldt)
	log.Printf("Monitor: New LDT %s with PID %d registered at %s\n", ldt.Name, ldt.Pid, ldt.Started)
}

func (m *Monitor) RemoveLDT(pid int) {
	for i, ldt := range m.processes {
		if ldt.Pid == pid {
			m.processes = append(m.processes[:i], m.processes[i+1:]...)
		}
	}
	log.Printf("Monitor: Removing LDT with PID %d\n", pid)
}

func (m *Monitor) ListLDTs() string {
	if len(m.processes) > 0 {
		var buffer bytes.Buffer
		for _, process := range m.processes {
			line := fmt.Sprintf("%d \t %s \t %s \t %v \t %d \t %t\n", process.Pid, process.Ldt, process.Name, process.Started, process.Port, process.Pairable)
			buffer.WriteString(line)
		}
		return buffer.String()
	}
	return " "
}

func (m *Monitor) SerializeLDTs() error {
	file, err := os.OpenFile(m.ldt_list_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Printf("Could not create file: %s\n", m.ldt_list_path)
		return err
	}
	defer file.Close()

	template := "%s\t%d\t%s\t%d\t%s\t%t\n"
	writer := bufio.NewWriter(file)
	for _, ldt := range m.processes {
		res := fmt.Sprintf(template, ldt.Ldt, ldt.Pid, ldt.Name, ldt.Port, ldt.Started, ldt.Pairable)
		writer.WriteString(res)
		writer.WriteString(string(ldt.Desc) + "\n")
	}

	writer.Flush()
	return nil
}

func (m *Monitor) DeserializeLDTs() error {
	if checkFileExists(m.ldt_list_path) {
		file, err := os.Open(m.ldt_list_path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var ldt string
			var pid int
			var name string
			var port int
			var day string
			var hour string
			var pairable bool
			var desc json.RawMessage
			_, err := fmt.Sscanf(scanner.Text(), "%s\t%d\t%s\t%d\t%s%s\t%t", &ldt, &pid, &name, &port, &day, &hour, &pairable)
			if err != nil {
				log.Println("Monitor: failed to deserialize the LDT", err)
			}

			started := day + " " + hour

			scanner.Scan()
			err = json.Unmarshal([]byte(scanner.Text()), &desc)

			if err != nil {
				log.Println("Monitor: Failed to deserialize the LDT description", err)
			}

			m.processes = append(m.processes, Process{
				Pid:      pid,
				Ldt:      ldt,
				Name:     name,
				Port:     port,
				Started:  started,
				Desc:     desc,
				Pairable: pairable})
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
			return err
		}
		os.Remove(m.ldt_list_path)
	}
	return nil
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func ldtIsRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Println(err)
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}

func (m *Monitor) GetPairaibleLDTAddress(name string) (string, error) {
	hostAddress, err := getIPAddress()
	if err != nil {
		return "", err
	}
	for i, ldt := range m.processes {
		if ldt.Pairable == true && ldt.LdtType() == name {
			res := hostAddress + ":" + fmt.Sprint(ldt.Port)
			m.processes[i].Pairable = false
			return res, nil
		}
	}
	return "No pairable LDT available", nil
}

func getIPAddress() (string, error) {
	hostname, _ := os.Hostname()

	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return "", errors.New(fmt.Sprint("Monitor: Failed wo obtain Host-IP Address"))
	}
	return ipAddr.IP.String(), nil
}
