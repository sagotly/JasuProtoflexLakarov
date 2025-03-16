package init

import (
	"log"
	"os"
)

func generateInstallScript(scriptPath string) error {
	// Создаём или перезаписываем файл скрипта.
	f, err := os.Create(scriptPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(initScriptContent)
	if err != nil {
		return err
	}

	// Делаем скрипт исполняемым.
	if err = os.Chmod(scriptPath, 0755); err != nil {
		return err
	}

	return nil
}

func main() {
	// Генерируем скрипт с именем install_script.sh
	err := generateInstallScript("init.sh")
	if err != nil {
		log.Fatalf("Ошибка при генерации скрипта: %v", err)
	}
	log.Println("Скрипт успешно сгенерирован: install_script.sh")
}
