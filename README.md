# vk_watcher
Отслеживает новые сообщения в группах Vkontakte, подает звуковой сигнал
Программа задумана, как некая универсальная заготовка для отслеживания изменений на веб страницах.
Планировщиком служит крон.


Синтаксис запуска такой: vk_parser "https://vk.com/gorod34" "news_vlg"

команда url имя_файла


Получился весьма нетребовательный продукт, работает в районе пары секунд, занимает примерно 10 мб оперативы. 
cron настроил так: 
*/5 * * * * %username команда

Каждые 5 минут

Чтобы из crona работали уведомления и звуковой сигнал, в etc/crontab пришлось загнать почти все содержимое virtual environment (можно получить по printenv)
В ближайших планах прикрутить к ней черный и белый списки. Подробно закомментирована.
