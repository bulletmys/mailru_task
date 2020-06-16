package analysis

//Инициализация и запуск обработчика
func RunWordAnalysis(paths []string, k int, word string) ([]PathValue, []error) {
	dataChan := make(chan PathReader)
	quitChan := make(chan bool)
	resChan := make(chan []PathValue)

	WCP := WordCounterParams{
		DataChan: dataChan,
		Quit:     quitChan,
		ResChan:  resChan,
		GoMax:    k,
		Word:     word,
	}

	//запуск планировщика для обслуживания дочерних горутин-обработчиков
	go WCP.WordCounter()

	dataErrors := GetData(paths, dataChan)

	quitChan <- true

	return <-resChan, dataErrors
}
