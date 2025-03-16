#!/bin/bash

set -e  # Остановить скрипт при ошибке

# Установка WireGuard
read -p "Установить WireGuard? (y/n): " install_wg
if [ "$install_wg" = "y" ] || [ "$install_wg" = "Y" ]; then
    echo "Устанавливаем WireGuard..."
    sudo apt update
    sudo apt install -y wireguard

    # Параметры по умолчанию
    INTERFACE_NAME=${1:-wg0}
    ADDRESS=${2:-10.0.0.1/24}
    LISTEN_PORT=${3:-5432}

    # Генерация ключей
    PRIVATE_KEY=$(sudo wg genkey)
    PUBLIC_KEY=$(echo "$PRIVATE_KEY" | sudo wg pubkey)
    echo "Сгенерированы ключи WireGuard:"
    echo "Private Key: $PRIVATE_KEY"
    echo "Public Key: $PUBLIC_KEY"

    # Конфигурация интерфейса WireGuard
    echo "Настраиваем интерфейс $INTERFACE_NAME..."
    sudo bash -c "cat <<EOF > /etc/wireguard/$INTERFACE_NAME.conf
[Interface]
PrivateKey = $PRIVATE_KEY
Address = $ADDRESS
ListenPort = $LISTEN_PORT
SaveConfig = true

PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
EOF"

    # Включаем IP forwarding
    echo "Включаем IP forwarding..."
    sudo bash -c "echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.conf"
    sudo sysctl -p

    # Запускаем WireGuard
    echo "Запускаем WireGuard..."
    sudo systemctl enable wg-quick@$INTERFACE_NAME
    sudo systemctl start wg-quick@$INTERFACE_NAME

    echo "WireGuard успешно установлен."
fi

# Установка Shadowsocks
read -p "Установить Shadowsocks? (y/n): " install_ss
if [ "$install_ss" = "y" ] || [ "$install_ss" = "Y" ]; then
    echo "Устанавливаем Shadowsocks..."
    sudo apt update
    sudo apt install -y python3-pip
    sudo pip3 install shadowsocks

    read -p "Введите порт для Shadowsocks (по умолчанию 8388): " SHADOWS_PORT
    SHADOWS_PORT=${SHADOWS_PORT:-8388}
    read -p "Введите пароль для Shadowsocks (по умолчанию 'password'): " SHADOWS_PASS
    SHADOWS_PASS=${SHADOWS_PASS:-password}

    # Создаём конфигурационный файл для Shadowsocks
    sudo bash -c "cat <<EOF > /etc/shadowsocks.json
{
    \"server\": \"0.0.0.0\",
    \"server_port\": $SHADOWS_PORT,
    \"password\": \"$SHADOWS_PASS\",
    \"timeout\": 300,
    \"method\": \"aes-256-cfb\"
}
EOF"

    # Запускаем Shadowsocks сервер в фоне
    sudo ssserver -c /etc/shadowsocks.json -d start
    echo "Shadowsocks успешно установлен."
fi

# Установка кастомного сервера
  read -p "Введите порт для кастомного сервера (по умолчанию 8080): " CUSTOM_PORT
  CUSTOM_PORT=${CUSTOM_PORT:-8080}
	sudo ./main -p $CUSTOM_PORT -i $INTERFACE_NAME

echo "Все установки завершены."