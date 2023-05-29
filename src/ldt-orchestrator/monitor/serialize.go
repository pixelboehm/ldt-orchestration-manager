package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	. "longevity/src/types"
	"os"
)

func (m *Monitor) SerializeLDTs() error {
	file, err := os.OpenFile(m.ldt_list_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Printf("Could not create file: %s\n", m.ldt_list_path)
		return err
	}
	defer file.Close()

	template := "%s\t%d\t%s\t%d\t%s\t%t\t%s\n"
	writer := bufio.NewWriter(file)
	for _, ldt := range m.processes {
		res := fmt.Sprintf(template, ldt.Ldt, ldt.Pid, ldt.Name, ldt.Port, ldt.Started, ldt.Pairable, ldt.DeviceMacAddress)
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
			var deviceMacAddress string
			var desc json.RawMessage
			_, err := fmt.Sscanf(scanner.Text(), "%s\t%d\t%s\t%d\t%s%s\t%t\t%s", &ldt, &pid, &name, &port, &day, &hour, &pairable, &deviceMacAddress)
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
				Pid:              pid,
				Ldt:              ldt,
				Name:             name,
				Port:             port,
				Started:          started,
				Desc:             desc,
				Pairable:         pairable,
				DeviceMacAddress: deviceMacAddress,
			})
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
			return err
		}
		os.Remove(m.ldt_list_path)
	}
	return nil
}
