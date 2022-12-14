### Как запустить программу
0. Предполагается, что у вас установлен Go, если нет, тогда читаем [https://go.dev/doc/install](https://go.dev/doc/install)
1. Клонируем репозиторий: ``git clone <repo>``
2. Далее есть 2 варианта запуска программы:  
2.a ``go run main.go https://example.com http://some.url``  
2.b используя файл urls.txt, лежащий в корне репы: ``cat urls.txt | xargs go run main.go``

### Вы можете переопределить значения по умолчанию, используя опции:
- максимальное число одновременно запущенных обработчиков ``--max-handlers <value>``. По умолчанию 5
- искомая строка, для которой необходимо подсчитать суммарное число вхождений на всех URL ``--needle <value>``. Default "Go"
- также можно включить подробный вывод с помощью ``--debug``

### Некоторые заметки
Метод http.Get внутри себя использует горутины, поэтому runtime.NumGoroutine() возвращает число намного большее, чем заданное в ``--max-handlers``. Это можно заметить, если передать фиктивные URL-ы программе и включить опцию ``--debug``, например ``go run main.go --debug a b c d``, runtime.NumGoroutine() вернет число 5 (функция main, плюс 4 обработчика на каждый URL).
Или запустить команду ``cat fake_urls.txt | xargs go run main.go --debug --max-handlers 5``