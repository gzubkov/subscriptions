#!/bin/bash

echo "=== Testing creating and calculating ==="
echo

BASE_URL="http://localhost:8080/api/v1"
USER_ID="2057e892-f836-43b7-81d2-ae6ffecd7989"

# 1. Создаем несколько подписок для тестирования
echo "1. Creating test subscriptions..."
echo

# Подписка 1: Yandex на 3 месяца (январь-март 2026)
echo "Creating Yandex subscription (Jan-Mar 2026)..."
curl -X POST "$BASE_URL/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex",
    "price": 299,
    "user_id": "'$USER_ID'",
    "start_date": "01-2026",
    "end_date": "03-2026"
  }'
echo -e "\n---\n"

# Подписка 2: Kion на 6 месяцев (январь-июнь 2026)
echo "Creating Kion subscription (Jan-Jun 2026)..."
curl -X POST "$BASE_URL/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Kion",
    "price": 199,
    "user_id": "'$USER_ID'",
    "start_date": "01-2026",
    "end_date": "06-2026"
  }'
echo -e "\n---\n"

# Подписка 3: Okko (с января 2026)
echo "Creating Okko subscription (from Jan 2026)..."
curl -X POST "$BASE_URL/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Okko",
    "price": 99,
    "user_id": "'$USER_ID'",
    "start_date": "01-2026"
  }'
echo -e "\n---\n"

# 2. Тесты расчета стоимости
echo "2. Testing total cost calculations..."
echo

# Тест 1: Первый квартал 2026 (3 месяца)
echo "Test 1: January-March 2026 - 3 months"
echo "Expected: Yandex (299 * 3) + Kion (199 * 3) + Okko (99 * 3) = 1791"
curl -X POST "$BASE_URL/subscriptions/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "01-2026",
    "end_date": "03-2026"
  }'
echo -e "\n---\n"

# Тест 2: Апрель 2026 (1 месяц)
echo "Test 2: April 2026 (1 month)"
echo "Expected: Kion (199) + Okko (99) = 298 (Yandex закончился)"
curl -X POST "$BASE_URL/subscriptions/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "04-2026",
    "end_date": "04-2026"
  }'
echo -e "\n---\n"

# Тест 3: Весь 2026 год (12 месяцев)
echo "Test 3: Full Year 2026 (12 months)"
echo "Expected: Yandex (299 * 3) + Kion (199 * 6) + Okko (99 * 12) = 3279"
curl -X POST "$BASE_URL/subscriptions/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "01-2026",
    "end_date": "12-2026"
  }'
echo -e "\n---\n"

# Тест 4: Фильтрация по сервису
echo "Test 4: Yandex only in January-March 2026"
echo "Expected: Yandex (299 * 3) = 897"
curl -X POST "$BASE_URL/subscriptions/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "01-2026",
    "end_date": "03-2026",
    "service_name": "Yandex"
  }'
echo -e "\n---\n"

# Тест 5: Фильтрация по пользователю
echo "Test 5: Filter by user"
curl -X POST "$BASE_URL/subscriptions/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "01-2026",
    "end_date": "12-2026",
    "user_id": "'$USER_ID'"
  }'
echo -e "\n---\n"

echo "=== Testing Complete ==="