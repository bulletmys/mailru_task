package analysis

import (
	"io"
	"sync"
)

// Структура для передачи данных в обработчик по подсчету слов.
type PathReader struct {
	Path   string
	Reader *io.ReadCloser
}

// Структура для записи результата работы алгоритма подсчета слов.
type PathValue struct {
	Path  string
	Value int64
	Err   error
}

// Структура, отвечающая за инициализацию обработчика по подсчету
// вхождений подстроки (слова) в строке.
type WordCounterParams struct {

	// Канал, принимающий новые данные для анализа
	DataChan <-chan PathReader

	// Канал, по которому горутина завершит работу
	Quit <-chan bool

	// Канал для записи результата работы горутины
	ResChan chan<- []PathValue

	// Максимальное кол-во горутин-обработчиков
	GoMax int

	// Искомая подстрока
	Word string
}

// Структура для обеспечения внутренней работы горутин-обработчиков.
type wrapperParams struct {
	mu   *sync.RWMutex
	wg   *sync.WaitGroup
	word string

	//Кол-во запущенных горутин-обработчиков в данный момент
	goAmount int

	//Слайс для хранения результатов подсчетов
	counts []PathValue

	//Очередь источников данных на обработку
	queue []PathReader

	//Размер буффера, считывающего данные
	readBuffSize int
}
