# Микросервис отправки писем

Микросервис получает из Kafka информацию о письмах, которые необходимо отправить, выполняет работу по подготовке письма по шаблону. 
- реализован паттерн Workers pool с целью конкурентного выполнения подготовки писем;
- гексагональная архитектура.
