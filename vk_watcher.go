package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/antchfx/htmlquery"
)

func main() {
	var file_content_list []string
	var filename, url string
	
	// здесь получаем путь к каталогу исполняемого файла
	path, _ := os.Executable()
	path = filepath.Dir(path)
	

	url = os.Args[1]      // 1-й аргумент URL
	filename = os.Args[2] // 2-й аргумент файл, для промежуточной записи

	// открываем файл с предыдущими записями, или создаем новый
	file, _ := os.OpenFile(path+"/"+filename, os.O_RDWR|os.O_CREATE, 0755)

	// этот блок читает файл в список
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		file_content_list = append(file_content_list, fileScanner.Text())
	}

	{ //блок парсинга
		var news_list []string             // список для записи новостей
		doc, err := htmlquery.LoadURL(url) // получаем страничку
		if err != nil {
			panic(err) // выходим при неуспехе отпечатав ошибку
		}

		// разбираем полученную страничку с помощью xpath
		s, err := htmlquery.QueryAll(doc, "//div[contains(@class,'pi_text')]")
		if len(s) == 0 || err != nil {
			log.Fatal("could not find viewpoint") // выходим, если ничего не нашли
		}

		// в этом блоке выдираем из каждого xpath его текст
		for i := 0; i < len(s); i++ {
			news_text := htmlquery.InnerText(s[i])
			news_list = append(news_list, news_text)
		}

		if !check_list(file_content_list, news_list) { // если списки не одинаковые...
			notification("VK Parser: "+filename, news_list[0], path+"/"+"Ахтунг.wav") //уведомляем о новости

			// в этом блоке переписываем файл новыми значениями. При повторном запуске будет сравниваться уже с ним
			file.Truncate(0)                 // очищаем файл и перемещаем указатель...
			file.Seek(0, 0)                  // в нулевую позицию. Иначе запишет много нулей
			writer := bufio.NewWriter(file)  // создаем построчного записывателя
			for _, news := range news_list { // для каждой новости в списке новостей (кэп)
				writer.WriteString(news + "\n") // записываем в файл строку + символ переноса строки (кэп)
				writer.Flush()                  // сохраняем изменения
			}
		}
		file.Close() // (кэп)
	}
}
func notification(title string, text string, sound_file string) {
	// просто выполняем системные команды
	exec.Command("/usr/bin/notify-send", "-t", "5000", "-u", "low", "-i", "-a", title, text).Start()
	exec.Command("/usr/bin/paplay", sound_file).Start()
}
func check_list(listOne, listTwo []string) bool {

	if len(listOne) != len(listTwo) {
		return false // если листы не равны по длине, можно сразу отдавать false
	}
	sort.Strings(listOne)
	sort.Strings(listTwo) // сортируем оба списка
	// теперь их элементы расположены рядом и в случае равных списков, элементы каждый раз должны быть равны
	for i := 0; i < len(listOne); i++ {
		if listOne[i] != listTwo[i] { // если хоть один элемент выбивается из этого правила...
			return false // значит списки точно не равны.
		} // именно для этого делалась сортировка, организовать сравнение списков за один цикл
	}
	return true // раз ни один элемент не равен, значит, списки равны
}
