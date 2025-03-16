#!/bin/bash

# Запрос порта у пользователя
read -p "Input server port: " port

# Проверка, что введено число от 1 до 65535
if ! [[ "$port" =~ ^[0-9]+$ ]] || [ "$port" -lt 1 ] || [ "$port" -gt 65535 ]; then
    echo "Ошибка: укажите корректный порт (число от 1 до 65535)"
    exit 1
fi

# Запуск бинарника main с указанным портом
echo "Starting server on pottу $port..."
sudo ./main -p $port 
