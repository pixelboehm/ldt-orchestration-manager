package monitoring_dependency_manager

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"longevity/src/communication"
	. "longevity/src/types"
	"net/http"
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
	rest.AddCustomHandler("/", m.mainpage)
	rest.Run(port)
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
	log.Printf("<Monitor>: New LDT %s with PID %d registered at %s\n", ldt.Name, ldt.Pid, ldt.Started)
}

func (m *Monitor) RemoveLDT(pid int) {
	var name string
	for i, ldt := range m.processes {
		if ldt.Pid == pid {
			name = ldt.Name
			m.processes = append(m.processes[:i], m.processes[i+1:]...)
		}
	}
	log.Printf("<Monitor>: Removing LDT %s with PID %d\n", name, pid)
}

func (m *Monitor) ListLDTs() string {
	if len(m.processes) > 0 {
		var buffer bytes.Buffer
		for _, process := range m.processes {
			line := fmt.Sprintf("%d\t%s\t%s\t%v\t%d\t%t\t%s\n", process.Pid, process.Ldt, process.Name, process.Started, process.Port, process.Pairable, process.DeviceMacAddress)
			buffer.WriteString(line)
		}
		return buffer.String()
	}
	return " "
}

func (m *Monitor) DoKeepAlive(seconds int64) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	for {
		log.Printf("<Monitor>: Currently Active LDTs %d\n", len(m.processes))
		for _, ldt := range m.processes {
			if !ldtIsRunning(ldt.Pid) {
				m.Stopped <- ldt.Pid
			}
		}
		<-ticker.C
	}
}

func (m *Monitor) GetLDTAddressForDevice(device Device) (string, error) {
	hostAddress, err := getIPAddress()
	if err != nil {
		return "", err
	}
	for i, ldt := range m.processes {
		if ldt.DeviceMacAddress == device.MacAddress && ldt.Pairable == false {
			var ldt_IPv4_full string
			if ldt.Port == 0 || ldt.Port == 80 {
				ldt_IPv4_full = hostAddress
			} else {
				ldt_IPv4_full = hostAddress + ":" + fmt.Sprint(ldt.Port)
			}
			return ldt_IPv4_full, nil
		}
		if ldt.Pairable == true && ldt.LdtType() == device.Name {
			res := hostAddress + ":" + fmt.Sprint(ldt.Port)
			m.processes[i].DeviceMacAddress = device.MacAddress
			m.processes[i].Pairable = false
			return res, nil
		}
	}
	return "nil", nil
}

func (m *Monitor) GetPidViaLdtName(name string) (int, error) {
	for _, ldt := range m.processes {
		if ldt.Name == name {
			return ldt.Pid, nil
		}
	}
	return -1, errors.New("<Monitor>: PID could not be resolved")
}

func (m *Monitor) mainpage(w http.ResponseWriter, r *http.Request) {
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
			"loadDescription": loadDescription,
			"convertTime":     convertTime,
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
