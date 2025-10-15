# 🚀 Деплой в Production

## Варианты деплоя

### 1. 🐳 Docker на VPS (Самый простой)

```bash
# На вашем сервере (DigitalOcean, AWS EC2, Hetzner, etc.)

# 1. Установите Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 2. Клонируйте репозиторий
git clone https://github.com/AlexRodving/advanced-user-api.git
cd advanced-user-api

# 3. Создайте .env файл
cat > .env << EOF
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=ВАША_БЕЗОПАСНАЯ_ПАРОЛЬ
DB_NAME=advanced_api
JWT_SECRET=ВАША_БЕЗОПАСНАЯ_SECRET_KEY
SERVER_PORT=8080
GIN_MODE=release
EOF

# 4. Запустите
sudo docker compose up -d

# 5. Проверьте
curl http://YOUR_SERVER_IP:8080/health
```

---

### 2. ☸️ Kubernetes (Для продакшна)

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: advanced-user-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: advanced-user-api
  template:
    metadata:
      labels:
        app: advanced-user-api
    spec:
      containers:
      - name: api
        image: your-registry/advanced-user-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-host
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
```

---

### 3. 🌐 Heroku (Бесплатный вариант)

```bash
# Установите Heroku CLI
curl https://cli-assets.heroku.com/install.sh | sh

# Логин
heroku login

# Создайте приложение
heroku create your-app-name

# Добавьте PostgreSQL
heroku addons:create heroku-postgresql:hobby-dev

# Установите переменные окружения
heroku config:set JWT_SECRET=your-secret-key
heroku config:set GIN_MODE=release

# Деплой
git push heroku main

# Откройте
heroku open
```

---

### 4. ☁️ AWS (Elastrolic Beanstalk)

```bash
# Установите EB CLI
pip install awsebcli

# Инициализируйте
eb init -p docker advanced-user-api

# Создайте окружение
eb create production

# Деплой
eb deploy

# Откройте
eb open
```

---

## 🔒 Важно для Production!

### 1. Переменные окружения

Обязательно измените:
- `JWT_SECRET` - длинная случайная строка (32+ символов)
- `DB_PASSWORD` - сильный пароль
- `GIN_MODE=release`

```bash
# Генерация безопасного секрета
openssl rand -base64 32
```

### 2. HTTPS

Используйте Nginx + Let's Encrypt:

```nginx
# /etc/nginx/sites-available/api
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Установите SSL:
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.yourdomain.com
```

### 3. Мониторинг

```yaml
# docker-compose.prod.yml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

### 4. Бэкапы БД

```bash
# Автоматический бэкап каждый день
0 2 * * * docker exec postgres pg_dump -U postgres advanced_api > /backups/db_$(date +\%Y\%m\%d).sql
```

---

## 📊 CI/CD

GitHub Actions уже настроен (`.github/workflows/ci.yml`):
- ✅ Тесты при каждом push
- ✅ Сборка Docker образа
- ✅ Security scan

---

## 💰 Стоимость хостинга

| Провайдер | Цена/месяц | Характеристики |
|-----------|------------|----------------|
| DigitalOcean | $5-10 | 1-2 GB RAM, 1 vCPU |
| Hetzner | €4-8 | 2-4 GB RAM, лучше по цене |
| AWS Free Tier | $0 (первый год) | t2.micro |
| Heroku | $0 (hobby) | Ограничения |

---

## 🎯 Чек-лист деплоя

- [ ] Создан production `.env`
- [ ] Сменён `JWT_SECRET`
- [ ] Сменён `DB_PASSWORD`
- [ ] `GIN_MODE=release`
- [ ] Настроен HTTPS
- [ ] Настроены бэкапы БД
- [ ] Настроен мониторинг
- [ ] Протестированы все endpoints
- [ ] Проверена безопасность

---

**Удачного деплоя!** 🚀

